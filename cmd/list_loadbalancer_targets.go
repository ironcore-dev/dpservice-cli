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

func ListLoadBalancerTargets(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts ListLoadBalancerTargetOptions
	)

	cmd := &cobra.Command{
		Use:     "lbtargets <--lb-id>",
		Short:   "List LoadBalancer Targets",
		Example: "dpservice-cli list lbtargets --lb-id=1",
		Aliases: LoadBalancerTargetAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunListLoadBalancerTargets(
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

type ListLoadBalancerTargetOptions struct {
	LoadBalancerID string
}

func (o *ListLoadBalancerTargetOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.LoadBalancerID, "lb-id", o.LoadBalancerID, "ID of the loadbalancer to get the targets for.")
}

func (o *ListLoadBalancerTargetOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"lb-id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunListLoadBalancerTargets(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	opts ListLoadBalancerTargetOptions,
) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer DpdkClose(cleanup)

	lbtargets, err := client.ListLoadBalancerTargets(ctx, opts.LoadBalancerID)
	if err != nil && lbtargets.Status.Code == 0 {
		return fmt.Errorf("error listing loadbalancer targets: %w", err)
	}

	return rendererFactory.RenderList("", os.Stdout, lbtargets)
}
