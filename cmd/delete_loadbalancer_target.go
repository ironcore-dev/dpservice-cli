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

	"github.com/onmetal/dpservice-cli/flag"
	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func DeleteLoadBalancerTarget(factory DPDKClientFactory) *cobra.Command {
	var (
		opts DeleteLoadBalancerTargetOptions
	)

	cmd := &cobra.Command{
		Use:     "lbtarget <--target-ip> <--lb-id>",
		Short:   "Delete a loadbalancer target",
		Example: "dpservice-cli delete lbtarget --target-ip=ff80::1 --lb-id=1",
		Aliases: LoadBalancerTargetAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunDeleteLoadBalancerTarget(cmd.Context(), factory, opts)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type DeleteLoadBalancerTargetOptions struct {
	TargetIP       netip.Addr
	LoadBalancerID string
}

func (o *DeleteLoadBalancerTargetOptions) AddFlags(fs *pflag.FlagSet) {
	flag.AddrVar(fs, &o.TargetIP, "target-ip", o.TargetIP, "LoadBalancer target IP to delete.")
	fs.StringVar(&o.LoadBalancerID, "lb-id", o.LoadBalancerID, "LoadBalancerID where to delete target.")
}

func (o *DeleteLoadBalancerTargetOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"target-ip", "lb-id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunDeleteLoadBalancerTarget(ctx context.Context, factory DPDKClientFactory, opts DeleteLoadBalancerTargetOptions) error {
	client, cleanup, err := factory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error deleting dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	if err := client.DeleteLoadBalancerTarget(ctx, opts.LoadBalancerID, opts.TargetIP); err != nil {
		return fmt.Errorf("error deleting loadbalancer target %s/%v: %v", opts.LoadBalancerID, opts.TargetIP, err)
	}
	fmt.Printf("Deleted loadbalancer target %s/%v\n", opts.LoadBalancerID, opts.TargetIP)

	return nil
}
