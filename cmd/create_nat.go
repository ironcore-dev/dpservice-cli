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
	"github.com/onmetal/dpservice-cli/flag"
	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func CreateNat(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts CreateNatOptions
	)

	cmd := &cobra.Command{
		Use:     "nat <interfaceID> [flags]",
		Short:   "Create a NAT",
		Example: "dpservice-cli create nat vm1 --natip=10.20.30.40 --minport=30000 --maxport=30100",
		Aliases: NatAliases,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			interfaceID := args[0]
			return RunCreateNat(
				cmd.Context(),
				dpdkClientFactory,
				rendererFactory,
				interfaceID,
				opts,
			)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type CreateNatOptions struct {
	NATVipIP netip.Addr
	MinPort  uint32
	MaxPort  uint32
}

func (o *CreateNatOptions) AddFlags(fs *pflag.FlagSet) {
	fs.Uint32Var(&o.MinPort, "minport", o.MinPort, "MinPort of NAT.")
	fs.Uint32Var(&o.MaxPort, "maxport", o.MaxPort, "MaxPort of NAT.")
	flag.AddrVar(fs, &o.NATVipIP, "natip", o.NATVipIP, "NAT IP to assign to the interface.")
}

func (o *CreateNatOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"minport", "maxport", "natip"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunCreateNat(ctx context.Context, dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory, interfaceID string, opts CreateNatOptions) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	renderer, err := rendererFactory.NewRenderer("created", os.Stdout)
	if err != nil {
		return fmt.Errorf("error creating renderer: %w", err)
	}

	nat, err := client.CreateNat(ctx, &api.Nat{
		NatMeta: api.NatMeta{
			InterfaceID: interfaceID,
		},
		Spec: api.NatSpec{
			NatVIPIP: opts.NATVipIP,
			MinPort:  opts.MinPort,
			MaxPort:  opts.MaxPort,
		},
	})
	if err != nil {
		return fmt.Errorf("error creating nat: %w", err)
	}

	if err := renderer.Render(nat); err != nil {
		return fmt.Errorf("error rendering nat: %w", err)
	}
	return nil
}
