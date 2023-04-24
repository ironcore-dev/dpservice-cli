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

func GetRoute(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts GetRouteOptions
	)

	cmd := &cobra.Command{
		Use:     "route [<prefix> <next-hop-vni> <next-hop-ip>...]",
		Short:   "Get or list route(s)",
		Aliases: RouteAliases,
		Args:    MultipleOfArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			keys, err := ParseRouteKeyArgs(args)
			if err != nil {
				return err
			}

			return RunGetRoute(
				cmd.Context(),
				dpdkClientFactory,
				rendererFactory,
				keys,
				opts,
			)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type GetRouteOptions struct {
	VNI uint32
}

func (o *GetRouteOptions) AddFlags(fs *pflag.FlagSet) {
	fs.Uint32Var(&o.VNI, "vni", o.VNI, "VNI to get the routes from.")
}

func (o *GetRouteOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"vni"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunGetRoute(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	keys []RouteKey,
	opts GetRouteOptions,
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

	renderer, err := rendererFactory.NewRenderer("", os.Stdout)
	if err != nil {
		return fmt.Errorf("error creating renderer: %w", err)
	}

	if len(keys) == 0 {
		routeList, err := client.ListRoutes(ctx, opts.VNI)
		if err != nil {
			return fmt.Errorf("error listing routes: %w", err)
		}

		if err := renderer.Render(routeList); err != nil {
			return fmt.Errorf("error rendering list: %w", err)
		}
		return nil
	}
	return fmt.Errorf("getting individual routes is not implemented")
}
