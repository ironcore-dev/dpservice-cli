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
	"strconv"

	"github.com/onmetal/dpservice-go-library/dpdk/api"
	"github.com/onmetal/dpservice-go-library/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func CreateRoute(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts CreateRouteOptions
	)

	cmd := &cobra.Command{
		Use:     "route <prefix> <next-hop-vni> <next-hop-ip>",
		Short:   "Create a route",
		Aliases: []string{"rt"},
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			prefix, err := netip.ParsePrefix(args[0])
			if err != nil {
				return fmt.Errorf("error parsing prefix: %w", err)
			}

			nextHopVNI, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return fmt.Errorf("error parsing next hop vni: %w", err)
			}

			nextHopIP, err := netip.ParseAddr(args[2])
			if err != nil {
				return fmt.Errorf("error parsing next hop ip: %w", err)
			}

			return RunCreateRoute(
				cmd.Context(),
				dpdkClientFactory,
				rendererFactory,
				prefix,
				uint32(nextHopVNI),
				nextHopIP,
				opts,
			)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type CreateRouteOptions struct {
	VNI uint32
}

func (o *CreateRouteOptions) AddFlags(fs *pflag.FlagSet) {
	fs.Uint32Var(&o.VNI, "vni", o.VNI, "Source VNI for the route.")
}

func (o *CreateRouteOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"vni"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunCreateRoute(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	prefix netip.Prefix,
	nextHopVNI uint32,
	nextHopIP netip.Addr,
	opts CreateRouteOptions,
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

	renderer, err := rendererFactory.NewRenderer(os.Stdout)
	if err != nil {
		return fmt.Errorf("error creating renderer: %w", err)
	}

	route, err := client.CreateRoute(ctx, &api.Route{
		RouteMeta: api.RouteMeta{
			VNI:    opts.VNI,
			Prefix: prefix,
			NextHop: api.RouteNextHop{
				VNI: nextHopVNI,
				IP:  nextHopIP,
			},
		},
		Spec: api.RouteSpec{},
	})
	if err != nil {
		return fmt.Errorf("error creating route: %w", err)
	}

	if err := renderer.Render(route); err != nil {
		return fmt.Errorf("error rendering route: %w", err)
	}
	return nil
}
