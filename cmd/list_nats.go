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
	"strings"

	"github.com/onmetal/net-dpservice-go/api"
	"github.com/onmetal/net-dpservice-go/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func ListNats(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nats",
		Short:   "List all nats",
		Example: "dpservice-cli list nats",
		Aliases: InterfaceAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunListNats(
				cmd.Context(),
				dpdkClientFactory,
				rendererFactory,
			)
		},
	}

	return cmd
}

type ListNatsOptions struct {
}

func (o *ListNatsOptions) AddFlags(fs *pflag.FlagSet) {
}

func (o *ListNatsOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	return nil
}

func RunListNats(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error getting dpdk client: %w", err)
	}
	defer DpdkClose(cleanup)

	interfaceList, err := client.ListInterfaces(ctx)
	if err != nil {
		return fmt.Errorf("error listing interfaces: %w", err)
	}
	interfaces := interfaceList.Items
	var nats []api.Nat
	for _, iface := range interfaces {
		nat, err := client.GetNat(ctx, iface.InterfaceMeta.ID)
		if strings.Contains(err.Error(), errors.StatusErrorString) {
			continue
		}
		if err != nil && !strings.Contains(err.Error(), errors.StatusErrorString) {
			return fmt.Errorf("error getting nat: %w", err)
		}
		nats = append(nats, *nat)
	}

	natList := api.NatList{
		TypeMeta: api.TypeMeta{Kind: api.NatListKind},
		Items:    nats,
	}

	return rendererFactory.RenderList("", os.Stdout, &natList)
}
