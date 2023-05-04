// Copyright 2022 OnMetal authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"net/netip"
	"os"

	"github.com/onmetal/dpservice-cli/dpdk/api"
	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func CreateLoadBalancerTarget(
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
) *cobra.Command {
	var (
		opts CreateLoadBalancerTargetOptions
	)

	cmd := &cobra.Command{
		Use:     "lbtarget <targetIP> [flags]",
		Short:   "Create a loadbalancer target",
		Example: "dpservice-cli create lbtarget ff80::5 --lb-id 2",
		Args:    cobra.ExactArgs(1),
		Aliases: LoadBalancerTargetAliases,
		RunE: func(cmd *cobra.Command, args []string) error {
			ip, err := netip.ParseAddr(args[0])
			if err != nil {
				return fmt.Errorf("error parsing ip: %w", err)
			}

			return RunCreateLoadBalancerTarget(
				cmd.Context(),
				dpdkClientFactory,
				rendererFactory,
				ip,
				opts,
			)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type CreateLoadBalancerTargetOptions struct {
	LoadBalancerID string
}

func (o *CreateLoadBalancerTargetOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.LoadBalancerID, "lb-id", o.LoadBalancerID, "ID of the loadbalancer to create the target for.")
}

func (o *CreateLoadBalancerTargetOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"lb-id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunCreateLoadBalancerTarget(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	ip netip.Addr,
	opts CreateLoadBalancerTargetOptions,
) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	renderer, err := rendererFactory.NewRenderer("created", os.Stdout)
	if err != nil {
		return fmt.Errorf("error creating renderer: %w", err)
	}

	targetIP := api.ProtoLbipToLbip(*api.LbipToProtoLbip(ip))
	res, err := client.CreateLoadBalancerTarget(ctx, &api.LoadBalancerTarget{
		TypeMeta:               api.TypeMeta{Kind: api.LoadBalancerTargetKind},
		LoadBalancerTargetMeta: api.LoadBalancerTargetMeta{ID: opts.LoadBalancerID},
		Spec:                   api.LoadBalancerTargetSpec{TargetIP: *targetIP},
	})
	if err != nil {
		return fmt.Errorf("error creating loadbalancer target: %w", err)
	}

	if err := renderer.Render(res); err != nil {
		return fmt.Errorf("error rendering loadbalancer target: %w", err)
	}
	return nil
}
