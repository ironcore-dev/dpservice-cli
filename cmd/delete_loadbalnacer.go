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

	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func DeleteLoadBalancer(factory DPDKClientFactory) *cobra.Command {
	var (
		opts DeleteLoadBalancerOptions
	)

	cmd := &cobra.Command{
		Use:     "loadbalancer <id> [<ids> ...]",
		Short:   "Delete loadbalancer(s)",
		Aliases: LoadBalancerAliases,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			loadbalancerIDs := args
			return RunDeleteLoadBalancer(cmd.Context(), factory, loadbalancerIDs, opts)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type DeleteLoadBalancerOptions struct {
}

func (o *DeleteLoadBalancerOptions) AddFlags(fs *pflag.FlagSet) {
}

func (o *DeleteLoadBalancerOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	return nil
}

func RunDeleteLoadBalancer(ctx context.Context, factory DPDKClientFactory, loadbalancerIDs []string, opts DeleteLoadBalancerOptions) error {
	client, cleanup, err := factory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	for _, loadbalancerID := range loadbalancerIDs {
		if err := client.DeleteLoadBalancer(ctx, loadbalancerID); err != nil {
			return fmt.Errorf("Error deleting loadbalancer %s: %v\n", loadbalancerID, err)
		}

		fmt.Println("Deleted loadbalancer", loadbalancerID)
	}
	return nil
}
