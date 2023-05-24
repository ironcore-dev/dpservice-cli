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

func GetVirtualIP(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts GetVirtualIPOptions
	)

	cmd := &cobra.Command{
		Use:     "virtualip <--interface-id>",
		Short:   "Get Virtual IP on interface",
		Example: "dpservice-cli get virtualip --interface-id=vm1",
		Aliases: VirtualIPAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunGetVirtualIP(
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

type GetVirtualIPOptions struct {
	InterfaceID string
}

func (o *GetVirtualIPOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.InterfaceID, "interface-id", o.InterfaceID, "Interface ID of the Virtual IP.")
}

func (o *GetVirtualIPOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"interface-id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunGetVirtualIP(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	opts GetVirtualIPOptions,
) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer DpdkClose(cleanup)

	virtualIP, err := client.GetVirtualIP(ctx, opts.InterfaceID)
	if err != nil && err != errors.ErrServerError {
		return fmt.Errorf("error getting virtual ip: %w", err)
	}

	virtualIP.TypeMeta.Kind = api.VirtualIPKind
	virtualIP.VirtualIPMeta.InterfaceID = opts.InterfaceID
	return rendererFactory.RenderObject("", os.Stdout, virtualIP)
}
