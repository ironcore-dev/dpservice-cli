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

	"github.com/onmetal/dpservice-cli/dpdk/api"
	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func DeleteVirtualIP(factory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts DeleteVirtualIPOptions
	)

	cmd := &cobra.Command{
		Use:     "virtualip <--interface-id>",
		Short:   "Delete virtual IP from interface",
		Example: "dpservice-cli delete virtualip --interface-id=vm1",
		Aliases: VirtualIPAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunDeleteVirtualIP(
				cmd.Context(),
				factory,
				rendererFactory,
				opts,
			)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type DeleteVirtualIPOptions struct {
	InterfaceID string
}

func (o *DeleteVirtualIPOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.InterfaceID, "interface-id", o.InterfaceID, "Interface ID of the Virtual IP.")
}

func (o *DeleteVirtualIPOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"interface-id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunDeleteVirtualIP(ctx context.Context, factory DPDKClientFactory, rendererFactory RendererFactory, opts DeleteVirtualIPOptions) error {
	client, cleanup, err := factory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	if err := client.DeleteVirtualIP(ctx, opts.InterfaceID); err != nil {
		return fmt.Errorf("error deleting virtual ip of interface %s: %v", opts.InterfaceID, err)
	}

	renderer, err := rendererFactory.NewRenderer("deleted", os.Stdout)
	if err != nil {
		return fmt.Errorf("error creating renderer: %w", err)
	}
	virtualIP := api.VirtualIP{
		TypeMeta: api.TypeMeta{Kind: api.VirtualIPKind},
		VirtualIPMeta: api.VirtualIPMeta{
			InterfaceID: opts.InterfaceID,
		},
		Status: api.Status{
			Message: "Deleted",
		},
	}
	if err := renderer.Render(&virtualIP); err != nil {
		return fmt.Errorf("error rendering prefix: %w", err)
	}

	return nil
}
