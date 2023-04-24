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

func DeleteVirtualIP(factory DPDKClientFactory) *cobra.Command {
	var (
		opts DeleteVirtualIPOptions
	)

	cmd := &cobra.Command{
		Use:     "virtualip <interface-id> [<interface-ids>...]",
		Short:   "Delete virtual ip(s)",
		Aliases: VirtualIPAliases,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			interfaceIDs := args
			return RunDeleteVirtualIP(cmd.Context(), factory, interfaceIDs, opts)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type DeleteVirtualIPOptions struct {
}

func (o *DeleteVirtualIPOptions) AddFlags(fs *pflag.FlagSet) {
}

func (o *DeleteVirtualIPOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	return nil
}

func RunDeleteVirtualIP(ctx context.Context, factory DPDKClientFactory, interfaceIDs []string, opts DeleteVirtualIPOptions) error {
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
		if err := client.DeleteVirtualIP(ctx, interfaceID); err != nil {
			fmt.Printf("Error deleting virtual ip of interface %s: %v\n", interfaceID, err)
		}

		fmt.Printf("Deleted virtual ip of interface %s\n", interfaceID)
	}
	return nil
}
