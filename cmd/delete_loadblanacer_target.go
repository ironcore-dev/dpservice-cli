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

	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func DeleteLoadBalancerTarget(factory DPDKClientFactory) *cobra.Command {
	var (
		opts DeleteLoadBalancerTargetOptions
	)

	cmd := &cobra.Command{
		Use:     "lbtarget [<targetIPs>] --lb-id <loadbalancerID>",
		Short:   "Delete a loadbalancer target",
		Example: "dpservice-cli delete lbtarget ff80::4 ff80::5 --lb-id=2",
		Aliases: LoadBalancerTargetAliases,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targets := args

			return RunDeleteLoadBalancerTarget(cmd.Context(), factory, targets, opts)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type DeleteLoadBalancerTargetOptions struct {
	LoadBalancerID string
}

func (o *DeleteLoadBalancerTargetOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.LoadBalancerID, "lb-id", o.LoadBalancerID, "LoadBalancerID where to delete target.")
}

func (o *DeleteLoadBalancerTargetOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"lb-id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunDeleteLoadBalancerTarget(ctx context.Context, factory DPDKClientFactory, targets []string, opts DeleteLoadBalancerTargetOptions) error {
	client, cleanup, err := factory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error deleting dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	for _, target := range targets {
		targetIP, err := netip.ParseAddr(target)
		if err != nil {
			return fmt.Errorf("not valid IP Address: %w", err)
		}
		if err := client.DeleteLoadBalancerTarget(ctx, opts.LoadBalancerID, targetIP); err != nil {
			return fmt.Errorf("error deleting loadbalancer target %s/%v: %v", opts.LoadBalancerID, target, err)
		}
		fmt.Printf("Deleted loadbalancer target %s/%v\n", opts.LoadBalancerID, target)
	}
	return nil
}
