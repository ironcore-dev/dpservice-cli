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

	"github.com/onmetal/dpservice-go-library/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func GetInterface(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts GetInterfaceOptions
	)

	cmd := &cobra.Command{
		Use:     "interface",
		Short:   "Get or list interface(s)",
		Aliases: InterfaceAliases,
		RunE: func(cmd *cobra.Command, args []string) error {
			interfaceIDs := args
			return RunGetInterface(
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

type GetInterfaceOptions struct {
}

func (o *GetInterfaceOptions) AddFlags(fs *pflag.FlagSet) {
}

func (o *GetInterfaceOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	return nil
}

func RunGetInterface(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	interfaceIDs []string,
	opts GetInterfaceOptions,
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
		ifaceList, err := client.ListInterfaces(ctx)
		if err != nil {
			return fmt.Errorf("error listing interfaces: %w", err)
		}

		if err := renderer.Render(ifaceList); err != nil {
			return fmt.Errorf("error rendering list: %w", err)
		}
		return nil
	}

	for _, interfaceID := range interfaceIDs {
		iface, err := client.GetInterface(ctx, interfaceID)
		if err != nil {
			return fmt.Errorf("error getting interface: %w", err)
		}

		if err := renderer.Render(iface); err != nil {
			return fmt.Errorf("error rendering interface %s: %w", interfaceID, err)
		}
	}
	return nil
}
