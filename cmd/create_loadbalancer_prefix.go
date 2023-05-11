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

func CreateLoadBalancerPrefix(
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
) *cobra.Command {
	var (
		opts CreateLoadBalancerPrefixOptions
	)

	cmd := &cobra.Command{
		Use:     "lbprefix <prefix> <--interface-id>",
		Short:   "Create a loadbalancer prefix",
		Example: "dpservice-cli create lbprefix ff80::1/64 --interface-id=vm1",
		Args:    cobra.ExactArgs(1),
		Aliases: PrefixAliases,
		RunE: func(cmd *cobra.Command, args []string) error {
			prefix, err := netip.ParsePrefix(args[0])
			if err != nil {
				return fmt.Errorf("error parsing prefix: %w", err)
			}

			return RunCreateLoadBalancerPrefix(
				cmd.Context(),
				dpdkClientFactory,
				rendererFactory,
				prefix,
				opts,
			)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type CreateLoadBalancerPrefixOptions struct {
	InterfaceID string
}

func (o *CreateLoadBalancerPrefixOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.InterfaceID, "interface-id", o.InterfaceID, "ID of the interface to create the prefix for.")
}

func (o *CreateLoadBalancerPrefixOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"interface-id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunCreateLoadBalancerPrefix(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	prefix netip.Prefix,
	opts CreateLoadBalancerPrefixOptions,
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

	res, err := client.CreateLoadBalancerPrefix(ctx, &api.Prefix{
		PrefixMeta: api.PrefixMeta{
			InterfaceID: opts.InterfaceID,
			Prefix:      prefix,
		},
		Spec: api.PrefixSpec{},
	})
	if err != nil {
		return fmt.Errorf("error creating prefix: %w", err)
	}

	if err := renderer.Render(res); err != nil {
		return fmt.Errorf("error rendering prefix: %w", err)
	}
	return nil
}
