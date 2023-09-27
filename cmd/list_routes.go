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
	"sort"
	"strings"

	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func ListRoutes(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts ListRoutesOptions
	)

	cmd := &cobra.Command{
		Use:     "routes <--vni>",
		Short:   "List routes of specified VNI",
		Example: "dpservice-cli list routes --vni=100",
		Aliases: RouteAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunGetRoute(
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

type ListRoutesOptions struct {
	VNI    uint32
	SortBy string
}

func (o *ListRoutesOptions) AddFlags(fs *pflag.FlagSet) {
	fs.Uint32Var(&o.VNI, "vni", o.VNI, "VNI to get the routes from.")
	fs.StringVar(&o.SortBy, "sort-by", "", "Column to sort by.")
}

func (o *ListRoutesOptions) MarkRequiredFlags(cmd *cobra.Command) error {
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
	opts ListRoutesOptions,
) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer DpdkClose(cleanup)

	routeList, err := client.ListRoutes(ctx, opts.VNI)
	if err != nil {
		return fmt.Errorf("error listing routes: %w", err)
	}

	// sort items in list
	routes := routeList.Items
	sort.SliceStable(routes, func(i, j int) bool {
		mi, mj := routes[i], routes[j]
		switch strings.ToLower(opts.SortBy) {
		case "nexthopvni":
			return mi.Spec.NextHop.VNI < mj.Spec.NextHop.VNI
		case "nexthopip":
			return mi.Spec.NextHop.IP.String() < mj.Spec.NextHop.IP.String()
		default:
			return mi.Spec.Prefix.String() < mj.Spec.Prefix.String()
		}
	})
	routeList.Items = routes

	return rendererFactory.RenderList("", os.Stdout, routeList)
}
