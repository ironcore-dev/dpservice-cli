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

	"github.com/onmetal/dpservice-cli/dpdk/api"
	"github.com/onmetal/dpservice-cli/dpdk/api/errors"
	"github.com/onmetal/dpservice-cli/flag"
	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func AddRoute(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts AddRouteOptions
	)

	cmd := &cobra.Command{
		Use:     "route <--prefix> <--next-hop-vni> <--next-hop-ip> <--vni>",
		Short:   "Add a route",
		Example: "dpservice-cli add route --prefix=10.100.3.0/24 --next-hop-vni=0 --next-hop-ip=fc00:2::64:0:1 --vni=100",
		Aliases: RouteAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunAddRoute(
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

type AddRouteOptions struct {
	Prefix     netip.Prefix
	NextHopVNI uint32
	NextHopIP  netip.Addr
	VNI        uint32
}

func (o *AddRouteOptions) AddFlags(fs *pflag.FlagSet) {
	flag.PrefixVar(fs, &o.Prefix, "prefix", o.Prefix, "Prefix for the route.")
	fs.Uint32Var(&o.NextHopVNI, "next-hop-vni", o.NextHopVNI, "Next hop VNI for the route.")
	flag.AddrVar(fs, &o.NextHopIP, "next-hop-ip", o.NextHopIP, "Next hop IP for the route.")
	fs.Uint32Var(&o.VNI, "vni", o.VNI, "Source VNI for the route.")
}

func (o *AddRouteOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"prefix", "next-hop-vni", "next-hop-ip", "vni"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunAddRoute(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	opts AddRouteOptions,
) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer DpdkClose(cleanup)

	route, err := client.AddRoute(ctx, &api.Route{
		RouteMeta: api.RouteMeta{
			VNI: opts.VNI,
		},
		Spec: api.RouteSpec{Prefix: &opts.Prefix,
			NextHop: &api.RouteNextHop{
				VNI: opts.NextHopVNI,
				IP:  &opts.NextHopIP,
			}},
	})
	if err != nil && err != errors.ErrServerError {
		return fmt.Errorf("error adding route: %w", err)
	}

	route.TypeMeta.Kind = api.RouteKind
	route.RouteMeta.VNI = opts.VNI
	route.Spec.Prefix = &opts.Prefix
	return rendererFactory.RenderObject(fmt.Sprintf("added, Next Hop IP: %s", opts.NextHopIP), os.Stdout, route)
}
