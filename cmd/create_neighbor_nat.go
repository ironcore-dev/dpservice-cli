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

	"github.com/onmetal/dpservice-cli/dpdk/api"
	"github.com/onmetal/dpservice-cli/flag"
	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func CreateNeighborNat(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts CreateNeighborNatOptions
	)

	cmd := &cobra.Command{
		Use:     "neighbornat <NatIP> [flags]",
		Short:   "Create a Neighbor NAT",
		Example: "dpservice-cli create neighbornat 10.20.30.40 --vni=100 --minport=30000 --maxport=30100 --underlayroute=ff80::1",
		Aliases: NeighborNatAliases,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			natVIPIP, err := netip.ParseAddr(args[0])
			if err != nil {
				return fmt.Errorf("error parsing nat vip ip: %w", err)
			}

			return RunCreateNeighborNat(
				cmd.Context(),
				dpdkClientFactory,
				rendererFactory,
				natVIPIP,
				opts,
			)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type CreateNeighborNatOptions struct {
	Vni           uint32
	MinPort       uint32
	MaxPort       uint32
	UnderlayRoute netip.Addr
}

func (o *CreateNeighborNatOptions) AddFlags(fs *pflag.FlagSet) {
	fs.Uint32Var(&o.Vni, "vni", o.Vni, "VNI of neighbor NAT.")
	fs.Uint32Var(&o.MinPort, "minport", o.MinPort, "MinPort of neighbor NAT.")
	fs.Uint32Var(&o.MaxPort, "maxport", o.MaxPort, "MaxPort of neighbor NAT.")
	flag.AddrVar(fs, &o.UnderlayRoute, "underlayroute", o.UnderlayRoute, "Underlay route of neighbor NAT.")
}

func (o *CreateNeighborNatOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"vni", "minport", "maxport", "underlayroute"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunCreateNeighborNat(ctx context.Context, dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory, natVIPIP netip.Addr, opts CreateNeighborNatOptions) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	err = client.CreateNeighborNat(ctx, &api.NeighborNat{
		NeighborNatMeta: api.NeighborNatMeta{
			NatVIPIP: natVIPIP,
		},
		Spec: api.NeighborNatSpec{
			Vni:           opts.Vni,
			MinPort:       opts.MinPort,
			MaxPort:       opts.MaxPort,
			UnderlayRoute: opts.UnderlayRoute,
		},
	})

	if err != nil {
		return fmt.Errorf("error creating neighbor nat: %w", err)
	}
	fmt.Printf("Neighbor NAT with IP: %s created\n", natVIPIP.String())

	return nil
}
