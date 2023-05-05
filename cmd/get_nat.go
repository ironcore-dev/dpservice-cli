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

	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func GetNat(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts GetNatOptions
	)

	cmd := &cobra.Command{
		Use:     "nat [<interfaceIDs>...]",
		Short:   "Get or list nat(s)",
		Aliases: NatAliases,
		RunE: func(cmd *cobra.Command, args []string) error {
			interfaceIDs := args
			return RunGetNat(
				cmd.Context(),
				dpdkClientFactory,
				rendererFactory,
				interfaceIDs,
				opts,
			)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type GetNatOptions struct {
}

func (o *GetNatOptions) AddFlags(fs *pflag.FlagSet) {
}

func (o *GetNatOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	return nil
}

func RunGetNat(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	interfaceIDs []string,
	opts GetNatOptions,
) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error getting dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	renderer, err := rendererFactory.NewRenderer("", os.Stdout)
	if err != nil {
		return fmt.Errorf("error creating renderer: %w", err)
	}

	if len(interfaceIDs) == 0 {
		return fmt.Errorf("listing all nats is not implemented")
	}

	for _, interfaceID := range interfaceIDs {
		nat, err := client.GetNat(ctx, interfaceID)
		if err != nil {
			return fmt.Errorf("error getting nat for interface %s: %v", interfaceID, err)
		}

		if err := renderer.Render(nat); err != nil {
			return fmt.Errorf("error rendering  nat: %w", err)
		}
	}
	return nil
}
