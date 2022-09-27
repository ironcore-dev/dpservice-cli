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

	"github.com/onmetal/dpservice-go-library/dpdk/api"
	"github.com/onmetal/dpservice-go-library/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func CreateVirtualIP(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts CreateVirtualIPOptions
	)

	cmd := &cobra.Command{
		Use:     "virtualip <ip>",
		Short:   "Create a virtual ip",
		Aliases: []string{"vip"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ip, err := netip.ParseAddr(args[0])
			if err != nil {
				return fmt.Errorf("error parsing ip: %w", err)
			}

			return RunCreateVirtualIP(
				cmd.Context(),
				dpdkClientFactory,
				rendererFactory,
				ip,
				opts,
			)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type CreateVirtualIPOptions struct {
	InterfaceID string
}

func (o *CreateVirtualIPOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.InterfaceID, "interface-id", o.InterfaceID, "Interface ID to create the virtual ip for.")
}

func (o *CreateVirtualIPOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"interface-id"} {
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
	ip netip.Addr,
	opts CreateVirtualIPOptions,
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

	renderer, err := rendererFactory.NewRenderer(os.Stdout)
	if err != nil {
		return fmt.Errorf("error creating renderer: %w", err)
	}

	virtualIP, err := client.CreateVirtualIP(ctx, &api.VirtualIP{
		VirtualIPMeta: api.VirtualIPMeta{
			InterfaceID: opts.InterfaceID,
			IP:          ip,
		},
		Spec: api.VirtualIPSpec{},
	})
	if err != nil {
		return fmt.Errorf("error creating virtual ip: %w", err)
	}

	if err := renderer.Render(virtualIP); err != nil {
		return fmt.Errorf("error rendering virtual ip: %w", err)
	}
	return nil
}
