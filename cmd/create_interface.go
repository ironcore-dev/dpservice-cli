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

	"github.com/onmetal/dpservice-cli/flag"
	"github.com/onmetal/dpservice-cli/util"
	"github.com/onmetal/net-dpservice-go/api"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func CreateInterface(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts CreateInterfaceOptions
	)

	cmd := &cobra.Command{
		Use:     "interface <--id> [<--ip>] <--vni> <--device>",
		Short:   "Create an interface",
		Example: "dpservice-cli create interface --id=vm4 --ipv4=10.200.1.4 --ipv6=2000:200:1::4 --vni=200 --device=net_tap5",
		Aliases: InterfaceAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunCreateInterface(
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

type CreateInterfaceOptions struct {
	ID          string
	VNI         uint32
	IPv4        netip.Addr
	IPv6        netip.Addr
	Device      string
	PxeServer   string
	PxeFileName string
}

func (o *CreateInterfaceOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ID, "id", o.ID, "ID of the interface.")
	fs.Uint32Var(&o.VNI, "vni", o.VNI, "VNI to add the interface to.")
	flag.AddrVar(fs, &o.IPv4, "ipv4", o.IPv4, "IPv4 address to assign to the interface.")
	flag.AddrVar(fs, &o.IPv6, "ipv6", netip.IPv6Unspecified(), "IPv6 address to assign to the interface.")
	fs.StringVar(&o.Device, "device", o.Device, "Device to allocate.")
	fs.StringVar(&o.PxeServer, "pxe-server", o.PxeServer, "PXE next server.")
	fs.StringVar(&o.PxeFileName, "pxe-file-name", o.PxeFileName, "PXE boot file name.")
}

func (o *CreateInterfaceOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"id", "vni", "ipv4", "device"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunCreateInterface(ctx context.Context, dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory, opts CreateInterfaceOptions) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer DpdkClose(cleanup)

	iface, err := client.CreateInterface(ctx, &api.Interface{
		InterfaceMeta: api.InterfaceMeta{
			ID: opts.ID,
		},
		Spec: api.InterfaceSpec{
			VNI:    opts.VNI,
			Device: opts.Device,
			IPv4:   &opts.IPv4,
			IPv6:   &opts.IPv6,
			PXE:    &api.PXE{Server: opts.PxeServer, FileName: opts.PxeFileName},
		},
	})
	if err != nil && iface.Status.Code == 0 {
		return fmt.Errorf("error creating interface: %w", err)
	}

	return rendererFactory.RenderObject(fmt.Sprintf("created, underlay route: %s", iface.Spec.UnderlayRoute), os.Stdout, iface)
}
