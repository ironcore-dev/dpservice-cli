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

func DeleteNeighborNat(factory DPDKClientFactory) *cobra.Command {
	var (
		opts DeleteNeighborNatOptions
	)

	cmd := &cobra.Command{
		Use:     "neighbornat <--natip> <--vni> <--minport> <--maxport>",
		Short:   "Delete neighbor nat",
		Example: "dpservice-cli delete neighbornat --natip=10.20.30.40 --vni=100 --minport=30000 --maxport=30100",
		Aliases: NeighborNatAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunDeleteNeighborNat(cmd.Context(), factory, opts)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type DeleteNeighborNatOptions struct {
	NatIP   netip.Addr
	Vni     uint32
	MinPort uint32
	MaxPort uint32
}

func (o *DeleteNeighborNatOptions) AddFlags(fs *pflag.FlagSet) {
	flag.AddrVar(fs, &o.NatIP, "natip", o.NatIP, "Neighbor NAT IP.")
	fs.Uint32Var(&o.Vni, "vni", o.Vni, "VNI of neighbor NAT.")
	fs.Uint32Var(&o.MinPort, "minport", o.MinPort, "MinPort of neighbor NAT.")
	fs.Uint32Var(&o.MaxPort, "maxport", o.MaxPort, "MaxPort of neighbor NAT.")
}

func (o *DeleteNeighborNatOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"natip", "vni", "minport", "maxport"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunDeleteNeighborNat(ctx context.Context, factory DPDKClientFactory, opts DeleteNeighborNatOptions) error {
	client, cleanup, err := factory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	neigbhorNat := api.NeighborNat{
		TypeMeta:        api.TypeMeta{Kind: api.NatKind},
		NeighborNatMeta: api.NeighborNatMeta{NatVIPIP: opts.NatIP},
		Spec: api.NeighborNatSpec{
			Vni:     opts.Vni,
			MinPort: opts.MinPort,
			MaxPort: opts.MaxPort,
		},
	}
	if err := client.DeleteNeighborNat(ctx, neigbhorNat); err != nil {
		return fmt.Errorf("error deleting neighbor nat with ip %s: %v", opts.NatIP, err)
	}

	fmt.Printf("Deleted neighbor NAT with IP %s\n", opts.NatIP)

	return nil
}
