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

func GetLoadBalancer(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts GetLoadBalancerOptions
	)

	cmd := &cobra.Command{
		Use:     "loadbalancer <ID>",
		Short:   "Get or list loadbalancer(s)",
		Aliases: LoadBalancerAliases,
		RunE: func(cmd *cobra.Command, args []string) error {
			loadbalancerIDs := args
			return RunGetLoadBalancer(
				cmd.Context(),
				dpdkClientFactory,
				rendererFactory,
				loadbalancerIDs,
				opts,
			)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type GetLoadBalancerOptions struct {
}

func (o *GetLoadBalancerOptions) AddFlags(fs *pflag.FlagSet) {
}

func (o *GetLoadBalancerOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	return nil
}

func RunGetLoadBalancer(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	loadbalancerIDs []string,
	opts GetLoadBalancerOptions,
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

	if len(loadbalancerIDs) == 0 {
		return fmt.Errorf("list loadbalancers not implemented")
	}

	for _, loadbalancerID := range loadbalancerIDs {
		lb, err := client.GetLoadBalancer(ctx, loadbalancerID)
		if err != nil {
			return fmt.Errorf("error getting loadbalancer: %w", err)
		}

		if err := renderer.Render(lb); err != nil {
			return fmt.Errorf("error rendering loadbalancer %s: %w", loadbalancerID, err)
		}
	}
	return nil
}
