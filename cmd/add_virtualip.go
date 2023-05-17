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
	"net/netip"
	"os"

	"github.com/onmetal/dpservice-cli/dpdk/api"
	"github.com/onmetal/dpservice-cli/flag"
	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func AddVirtualIP(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts AddVirtualIPOptions
	)

	cmd := &cobra.Command{
		Use:     "virtualip <--vip> <--interface-id>",
		Short:   "Add a virtual IP to interface.",
		Example: "dpservice-cli add virtualip --vip=20.20.20.20 --interface-id=vm1",
		Aliases: VirtualIPAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunAddVirtualIP(
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

type AddVirtualIPOptions struct {
	Vip         netip.Addr
	InterfaceID string
}

func (o *AddVirtualIPOptions) AddFlags(fs *pflag.FlagSet) {
	flag.AddrVar(fs, &o.Vip, "vip", o.Vip, "Virtual IP to add on interface.")
	fs.StringVar(&o.InterfaceID, "interface-id", o.InterfaceID, "Interface ID to add the virtual ip for.")
}

func (o *AddVirtualIPOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"vip", "interface-id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunAddVirtualIP(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	opts AddVirtualIPOptions,
) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	renderer, err := rendererFactory.NewRenderer("added", os.Stdout)
	if err != nil {
		return fmt.Errorf("error creating renderer: %w", err)
	}

	virtualIP, err := client.AddVirtualIP(ctx, &api.VirtualIP{
		VirtualIPMeta: api.VirtualIPMeta{
			InterfaceID: opts.InterfaceID,
			IP:          opts.Vip,
		},
	})
	if err != nil {
		return fmt.Errorf("error adding virtual ip: %w", err)
	}

	if err := renderer.Render(virtualIP); err != nil {
		return fmt.Errorf("error rendering virtual ip: %w", err)
	}
	fmt.Println("Underlay route is:", virtualIP.Spec.UnderlayRoute)
	return nil
}
