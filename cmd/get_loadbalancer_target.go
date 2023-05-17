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
	"os"

	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func GetLoadBalancerTargets(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts GetLoadBalancerTargetOptions
	)

	cmd := &cobra.Command{
		Use:     "lbtarget <--lb-id>",
		Short:   "Get LoadBalancer Targets",
		Example: "dpservice-cli get lbtarget --lb-id=1",
		Aliases: LoadBalancerTargetAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunGetLoadBalancerTargets(
				cmd.Context(),
				dpdkClientFactory,
				rendererFactory,
				opts,
			)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type GetLoadBalancerTargetOptions struct {
	LoadBalancerID string
}

func (o *GetLoadBalancerTargetOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.LoadBalancerID, "lb-id", o.LoadBalancerID, "ID of the loadbalancer to get the targets for.")
}

func (o *GetLoadBalancerTargetOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"lb-id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunGetLoadBalancerTargets(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	opts GetLoadBalancerTargetOptions,
) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error getting dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	renderer, err := rendererFactory.NewRenderer("", os.Stdout)
	if err != nil {
		return fmt.Errorf("error creating renderer: %w", err)
	}

	lbtarget, err := client.GetLoadBalancerTargets(ctx, opts.LoadBalancerID)
	if err != nil {
		return fmt.Errorf("error getting loadbalancer target for interface %s: %v", opts.LoadBalancerID, err)
	}

	if err := renderer.Render(lbtarget); err != nil {
		return fmt.Errorf("error rendering loadbalancer target: %w", err)
	}

	return nil
}
