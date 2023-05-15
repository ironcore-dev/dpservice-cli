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

	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func DeleteInterface(factory DPDKClientFactory) *cobra.Command {
	var (
		opts DeleteInterfaceOptions
	)

	cmd := &cobra.Command{
		Use:     "interface <interfaceID> [<interfaceIDs> ...]",
		Short:   "Delete interface(s)",
		Example: "dpservice-cli delete interface vm1",
		Aliases: InterfaceAliases,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			interfaceIDs := args
			return RunDeleteInterface(cmd.Context(), factory, interfaceIDs, opts)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type DeleteInterfaceOptions struct {
}

func (o *DeleteInterfaceOptions) AddFlags(fs *pflag.FlagSet) {
}

func (o *DeleteInterfaceOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	return nil
}

func RunDeleteInterface(ctx context.Context, factory DPDKClientFactory, interfaceIDs []string, opts DeleteInterfaceOptions) error {
	client, cleanup, err := factory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	for _, interfaceID := range interfaceIDs {
		if err := client.DeleteInterface(ctx, interfaceID); err != nil {
			return fmt.Errorf("error deleting interface %s: %v", interfaceID, err)
		}

		fmt.Println("Deleted interface", interfaceID)
	}
	return nil
}
