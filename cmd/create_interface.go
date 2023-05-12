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

func CreateInterface(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts CreateInterfaceOptions
	)

	cmd := &cobra.Command{
		Use:     "interface <id>",
		Short:   "Create an interface",
		Example: "dpservice-cli create interface vm4 --ips=10.200.1.4 --ips=2000:200:1::4 --vni=200 --device=net_tap5",
		Aliases: InterfaceAliases,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			interfaceID := args[0]
			return RunCreateInterface(
				cmd.Context(),
				dpdkClientFactory,
				rendererFactory,
				interfaceID,
				opts,
			)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type CreateInterfaceOptions struct {
	VNI    uint32
	IPs    []netip.Addr
	Device string
}

func (o *CreateInterfaceOptions) AddFlags(fs *pflag.FlagSet) {
	fs.Uint32Var(&o.VNI, "vni", o.VNI, "VNI to add the interface to.")
	flag.AddrSliceVar(fs, &o.IPs, "ips", o.IPs, "IPs to assign to the interface.")
	fs.StringVar(&o.Device, "device", o.Device, "Device to allocate.")
}

func (o *CreateInterfaceOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"vni", "ips"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunCreateInterface(ctx context.Context, dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory, interfaceID string, opts CreateInterfaceOptions) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	renderer, err := rendererFactory.NewRenderer("created", os.Stdout)
	if err != nil {
		return fmt.Errorf("error creating renderer: %w", err)
	}

	iface, err := client.CreateInterface(ctx, &api.Interface{
		InterfaceMeta: api.InterfaceMeta{
			ID: interfaceID,
		},
		Spec: api.InterfaceSpec{
			VNI:    opts.VNI,
			Device: opts.Device,
			IPs:    opts.IPs,
		},
	})
	if err != nil {
		return fmt.Errorf("error creating interface: %w", err)
	}

	if err := renderer.Render(iface); err != nil {
		return fmt.Errorf("error rendering interface: %w", err)
	}
	return nil
}