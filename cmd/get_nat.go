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
	"os"

	"github.com/ironcore-dev/dpservice-cli/util"
	"github.com/ironcore-dev/dpservice-go/api"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func GetNat(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts GetNatOptions
	)

	cmd := &cobra.Command{
		Use:     "nat <--interface-id>",
		Short:   "Get NAT on interface",
		Example: "dpservice-cli get nat --interface-id=vm1",
		Aliases: NatAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunGetNat(
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

type GetNatOptions struct {
	InterfaceID string
}

func (o *GetNatOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.InterfaceID, "interface-id", o.InterfaceID, "Interface ID of the NAT.")
}

func (o *GetNatOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	return nil
}

func RunGetNat(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	opts GetNatOptions,
) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer DpdkClose(cleanup)

	if opts.InterfaceID == "" {
		ifaces, err := client.ListInterfaces(ctx)
		if err != nil && ifaces.Status.Code == 0 {
			return fmt.Errorf("error listing interfaces: %w", err)
		}
		natList := api.NatList{
			TypeMeta: api.TypeMeta{Kind: api.NatListKind},
		}
		for _, iface := range ifaces.Items {
			nat, err := client.GetNat(ctx, iface.ID)
			if err != nil && nat.Status.Code == 0 {
				return fmt.Errorf("error getting nat: %w", err)
			}
			natList.Items = append(natList.Items, *nat)
		}

		return rendererFactory.RenderList("", os.Stdout, &natList)
	}

	nat, err := client.GetNat(ctx, opts.InterfaceID)
	if err != nil && nat.Status.Code == 0 {
		return fmt.Errorf("error getting nat: %w", err)
	}

	return rendererFactory.RenderObject("", os.Stdout, nat)
}
