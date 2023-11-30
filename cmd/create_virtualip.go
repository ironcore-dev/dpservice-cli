// Copyright 2022 IronCore authors
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

	"github.com/ironcore-dev/dpservice-cli/flag"
	"github.com/ironcore-dev/dpservice-cli/util"
	"github.com/ironcore-dev/dpservice-go/api"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func CreateVirtualIP(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts CreateVirtualIPOptions
	)

	cmd := &cobra.Command{
		Use:     "virtualip <--vip> <--interface-id>",
		Short:   "Create a virtual IP on interface.",
		Example: "dpservice-cli create virtualip --vip=20.20.20.20 --interface-id=vm1",
		Aliases: VirtualIPAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunCreateVirtualIP(
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

type CreateVirtualIPOptions struct {
	Vip         netip.Addr
	InterfaceID string
}

func (o *CreateVirtualIPOptions) AddFlags(fs *pflag.FlagSet) {
	flag.AddrVar(fs, &o.Vip, "vip", o.Vip, "Virtual IP to create on interface.")
	fs.StringVar(&o.InterfaceID, "interface-id", o.InterfaceID, "Interface ID where to create the virtual IP.")
}

func (o *CreateVirtualIPOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"vip", "interface-id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunCreateVirtualIP(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	opts CreateVirtualIPOptions,
) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer DpdkClose(cleanup)

	virtualIP, err := client.CreateVirtualIP(ctx, &api.VirtualIP{
		VirtualIPMeta: api.VirtualIPMeta{
			InterfaceID: opts.InterfaceID,
		},
		Spec: api.VirtualIPSpec{
			IP: &opts.Vip,
		},
	})
	if err != nil && virtualIP.Status.Code == 0 {
		return fmt.Errorf("error creating virtual ip: %w", err)
	}

	return rendererFactory.RenderObject(fmt.Sprintf("created, underlay route: %s", virtualIP.Spec.UnderlayRoute), os.Stdout, virtualIP)
}
