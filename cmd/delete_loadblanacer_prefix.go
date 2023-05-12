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

func DeleteLoadBalancerPrefix(factory DPDKClientFactory) *cobra.Command {
	var (
		opts DeleteLoadBalancerPrefixOptions
	)

	cmd := &cobra.Command{
		Use:     "lbprefix <prefix> [<prefixes>...]",
		Short:   "Delete a loadbalancer prefix",
		Example: "dpservice-cli delete lbprefix ff80::1/64 --interface-id=vm1",
		Aliases: LoadBalancerPrefixAliases,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			prefixes, err := ParsePrefixArgs(args)
			if err != nil {
				return err
			}

			return RunDeleteLoadBalancerPrefix(cmd.Context(), factory, prefixes, opts)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type DeleteLoadBalancerPrefixOptions struct {
	InterfaceID string
}

func (o *DeleteLoadBalancerPrefixOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.InterfaceID, "interface-id", o.InterfaceID, "Interface ID of the loadbalancer prefix.")
}

func (o *DeleteLoadBalancerPrefixOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"interface-id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunDeleteLoadBalancerPrefix(ctx context.Context, factory DPDKClientFactory, prefixes []netip.Prefix, opts DeleteLoadBalancerPrefixOptions) error {
	client, cleanup, err := factory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error deleting dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	for _, prefix := range prefixes {
		if err := client.DeleteLoadBalancerPrefix(ctx, opts.InterfaceID, prefix); err != nil {
			fmt.Printf("Error deleting loadbalancer prefix %s/%v: %v\n", opts.InterfaceID, prefix, err)
		}
		fmt.Printf("Deleted loadbalancer prefix %s/%v\n", opts.InterfaceID, prefix)
	}
	return nil
}
