// Copyright 2022 IronCore authors
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

	"github.com/ironcore-dev/dpservice-cli/flag"
	"github.com/ironcore-dev/dpservice-cli/util"
	"github.com/ironcore-dev/dpservice-go/api"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func CreateNat(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts CreateNatOptions
	)

	cmd := &cobra.Command{
		Use:     "nat <--interface-id> <--nat-ip> <--minport> <--maxport>",
		Short:   "Create a NAT on interface",
		Example: "dpservice-cli create nat --interface-id=vm1 --nat-ip=10.20.30.40 --minport=30000 --maxport=30100",
		Aliases: NatAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunCreateNat(
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

type CreateNatOptions struct {
	InterfaceID string
	NatIP       netip.Addr
	MinPort     uint32
	MaxPort     uint32
}

func (o *CreateNatOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.InterfaceID, "interface-id", o.InterfaceID, "Interface ID where to create NAT.")
	fs.Uint32Var(&o.MinPort, "minport", o.MinPort, "MinPort of NAT.")
	fs.Uint32Var(&o.MaxPort, "maxport", o.MaxPort, "MaxPort of NAT.")
	flag.AddrVar(fs, &o.NatIP, "nat-ip", o.NatIP, "NAT IP to assign to the interface.")
}

func (o *CreateNatOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"interface-id", "minport", "maxport", "nat-ip"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunCreateNat(ctx context.Context, dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory, opts CreateNatOptions) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer DpdkClose(cleanup)

	nat, err := client.CreateNat(ctx, &api.Nat{
		NatMeta: api.NatMeta{
			InterfaceID: opts.InterfaceID,
		},
		Spec: api.NatSpec{
			NatIP:   &opts.NatIP,
			MinPort: opts.MinPort,
			MaxPort: opts.MaxPort,
		},
	})
	if err != nil && nat.Status.Code == 0 {
		return fmt.Errorf("error creating nat: %w", err)
	}

	return rendererFactory.RenderObject(fmt.Sprintf("created, underlay route: %s", nat.Spec.UnderlayRoute), os.Stdout, nat)
}
