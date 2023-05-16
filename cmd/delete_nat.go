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

func DeleteNat(factory DPDKClientFactory) *cobra.Command {
	var (
		opts DeleteNatOptions
	)

	cmd := &cobra.Command{
		Use:     "nat <--interface-id>",
		Short:   "Delete nat from interface",
		Example: "dpservice-cli delete nat --interface-id=vm1",
		Aliases: NatAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunDeleteNat(cmd.Context(), factory, opts)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type DeleteNatOptions struct {
	InterfaceID string
}

func (o *DeleteNatOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.InterfaceID, "interface-id", o.InterfaceID, "Interface ID of the Virtual IP.")
}

func (o *DeleteNatOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"interface-id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunDeleteNat(ctx context.Context, factory DPDKClientFactory, opts DeleteNatOptions) error {
	client, cleanup, err := factory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	if err := client.DeleteNat(ctx, opts.InterfaceID); err != nil {
		return fmt.Errorf("error deleting nat of interface %s: %v", opts.InterfaceID, err)
	}

	fmt.Printf("Deleted NAT of interface %s\n", opts.InterfaceID)

	return nil
}
