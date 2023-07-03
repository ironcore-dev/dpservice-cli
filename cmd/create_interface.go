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
	"github.com/onmetal/net-dpservice-go/errors"
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
		Example: "dpservice-cli add interface --id=vm4 --ip=10.200.1.4 --ip=2000:200:1::4 --vni=200 --device=net_tap5",
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
	IP          []netip.Addr
	Device      string
	PxeServer   string
	PxeFileName string
}

func (o *CreateInterfaceOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ID, "id", o.ID, "ID of the interface.")
	fs.Uint32Var(&o.VNI, "vni", o.VNI, "VNI to add the interface to.")
	flag.AddrSliceVar(fs, &o.IP, "ip", o.IP, "IP to assign to the interface.")
	fs.StringVar(&o.Device, "device", o.Device, "Device to allocate.")
	fs.StringVar(&o.PxeServer, "pxe-server", o.PxeServer, "PXE next server.")
	fs.StringVar(&o.PxeFileName, "pxe-file-name", o.PxeFileName, "PXE boot file name.")
}

func (o *CreateInterfaceOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"id", "vni", "ip", "device"} {
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
			IPs:    opts.IP,
			PXE:    &api.PXE{Server: opts.PxeServer, FileName: opts.PxeFileName},
		},
	})
	if err != nil && err != errors.ErrServerError {
		return fmt.Errorf("error adding interface: %w", err)
	}

	return rendererFactory.RenderObject(fmt.Sprintf("added, underlay route: %s", iface.Spec.UnderlayRoute), os.Stdout, iface)
}
