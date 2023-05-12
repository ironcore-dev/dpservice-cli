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

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func ListInterfaces(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "interfaces",
		Short:   "List all interfaces",
		Example: "dpservice-cli list interfaces",
		Aliases: InterfaceAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunListInterfaces(
				cmd.Context(),
				dpdkClientFactory,
				rendererFactory,
			)
		},
	}

	return cmd
}

type ListInterfacesOptions struct {
}

func (o *ListInterfacesOptions) AddFlags(fs *pflag.FlagSet) {
}

func (o *ListInterfacesOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	return nil
}

func RunListInterfaces(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
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

	interfaceList, err := client.ListInterfaces(ctx)
	if err != nil {
		return fmt.Errorf("error listing firewall rules: %w", err)
	}

	if err := renderer.Render(interfaceList); err != nil {
		return fmt.Errorf("error rendering interfaces: %w", err)
	}
	return nil
}
