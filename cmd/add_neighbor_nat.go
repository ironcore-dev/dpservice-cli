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

	"github.com/onmetal/dpservice-cli/flag"
	"github.com/onmetal/dpservice-cli/util"
	"github.com/onmetal/net-dpservice-go/api"
	"github.com/onmetal/net-dpservice-go/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func AddNeighborNat(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts AddNeighborNatOptions
	)

	cmd := &cobra.Command{
		Use:     "neighbornat <--nat-ip> <--vni> <--minport> <--maxport> <--underlayroute>",
		Short:   "Add a Neighbor NAT",
		Example: "dpservice-cli add neighbornat --nat-ip=10.20.30.40 --vni=100 --minport=30000 --maxport=30100 --underlayroute=ff80::1",
		Aliases: NeighborNatAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunAddNeighborNat(
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

type AddNeighborNatOptions struct {
	NatIP         netip.Addr
	Vni           uint32
	MinPort       uint32
	MaxPort       uint32
	UnderlayRoute netip.Addr
}

func (o *AddNeighborNatOptions) AddFlags(fs *pflag.FlagSet) {
	flag.AddrVar(fs, &o.NatIP, "nat-ip", o.NatIP, "Neighbor NAT IP.")
	fs.Uint32Var(&o.Vni, "vni", o.Vni, "VNI of neighbor NAT.")
	fs.Uint32Var(&o.MinPort, "minport", o.MinPort, "MinPort of neighbor NAT.")
	fs.Uint32Var(&o.MaxPort, "maxport", o.MaxPort, "MaxPort of neighbor NAT.")
	flag.AddrVar(fs, &o.UnderlayRoute, "underlayroute", o.UnderlayRoute, "Underlay route of neighbor NAT.")
}

func (o *AddNeighborNatOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"nat-ip", "vni", "minport", "maxport", "underlayroute"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunAddNeighborNat(ctx context.Context, dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory, opts AddNeighborNatOptions) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer DpdkClose(cleanup)

	neigbhorNat := &api.NeighborNat{
		TypeMeta:        api.TypeMeta{Kind: api.NeighborNatKind},
		NeighborNatMeta: api.NeighborNatMeta{NatVIPIP: &opts.NatIP},
		Spec: api.NeighborNatSpec{
			Vni:           opts.Vni,
			MinPort:       opts.MinPort,
			MaxPort:       opts.MaxPort,
			UnderlayRoute: &opts.UnderlayRoute,
		},
	}

	nnat, err := client.AddNeighborNat(ctx, neigbhorNat)
	if err != nil && err != errors.ErrServerError {
		return fmt.Errorf("error adding neighbor nat: %w", err)
	}

	return rendererFactory.RenderObject(fmt.Sprintf("added, underlay route: %s", nnat.Spec.UnderlayRoute), os.Stdout, nnat)
}
