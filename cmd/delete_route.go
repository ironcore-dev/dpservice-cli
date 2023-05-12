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

func DeleteRoute(factory DPDKClientFactory) *cobra.Command {
	var (
		opts DeleteRouteOptions
	)

	cmd := &cobra.Command{
		Use:     "route <prefix> <next-hop-vni> <next-hop-ip> [<prefix-n> <next-hop-ip-n> <next-hop-ip-n>...]",
		Short:   "Delete a route",
		Example: "dpservice-cli delete route 10.100.2.0/24 0 fc00:2::64:0:1 --vni=100",
		Aliases: RouteAliases,
		Args:    MultipleOfArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			keys, err := ParseRouteKeyArgs(args)
			if err != nil {
				return err
			}

			return RunDeleteRoute(cmd.Context(), factory, keys, opts)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type DeleteRouteOptions struct {
	VNI uint32
}

func (o *DeleteRouteOptions) AddFlags(fs *pflag.FlagSet) {
	fs.Uint32Var(&o.VNI, "vni", o.VNI, "VNI of the route.")
}

func (o *DeleteRouteOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"vni"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunDeleteRoute(ctx context.Context, factory DPDKClientFactory, keys []RouteKey, opts DeleteRouteOptions) error {
	client, cleanup, err := factory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error deleting dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	for _, key := range keys {
		if err := client.DeleteRoute(ctx, opts.VNI, key.Prefix, key.NextHopVNI, key.NextHopIP); err != nil {
			return fmt.Errorf("error deleting route %d-%v:%d-%v: %v", opts.VNI, key.Prefix, key.NextHopVNI, key.NextHopIP, err)
		}
		fmt.Printf("Deleted route %d-%v:%d-%v\n", opts.VNI, key.Prefix, key.NextHopVNI, key.NextHopIP)
	}

	return nil
}
