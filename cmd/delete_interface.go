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
	"github.com/onmetal/dpservice-cli/dpdk/api/errors"
	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func DeleteInterface(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts DeleteInterfaceOptions
	)

	cmd := &cobra.Command{
		Use:     "interface <--id>",
		Short:   "Delete interface",
		Example: "dpservice-cli delete interface --id=vm1",
		Aliases: InterfaceAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunDeleteInterface(
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

type DeleteInterfaceOptions struct {
	ID string
}

func (o *DeleteInterfaceOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ID, "id", o.ID, "Interface ID to delete.")
}

func (o *DeleteInterfaceOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunDeleteInterface(ctx context.Context, dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory, opts DeleteInterfaceOptions) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer DpdkClose(cleanup)

	iface, err := client.DeleteInterface(ctx, opts.ID)
	if err != nil && err != errors.ErrServerError {
		return fmt.Errorf("error deleting interface: %w", err)
	}

	iface.TypeMeta.Kind = api.InterfaceKind
	iface.InterfaceMeta.ID = opts.ID
	return rendererFactory.RenderObject("deleted", os.Stdout, iface)
}
