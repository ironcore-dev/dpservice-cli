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
	"os"

	"github.com/ironcore-dev/dpservice-cli/util"
	"github.com/ironcore-dev/dpservice-go/api"
	"github.com/ironcore-dev/dpservice-go/errors"
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

	if opts.InterfaceID == "" {
		ifaces, err := client.ListInterfaces(ctx)
		if err != nil && ifaces.Status.Code == 0 {
			return fmt.Errorf("error listing interfaces: %w", err)
		}
		virtualIPs := make([]*api.VirtualIP, 0, len(ifaces.Items))
		for _, iface := range ifaces.Items {
			vip, err := client.GetVirtualIP(ctx, iface.ID, errors.Ignore(errors.SNAT_NO_DATA))
			if err != nil && vip.Status.Code == 0 {
				return fmt.Errorf("error getting virtual ip: %w", err)
			}
			if vip.Status.Code == 0 {
				virtualIPs = append(virtualIPs, vip)
			}
		}
		if len(virtualIPs) == 0 {
			noVipFound := api.VirtualIP{
				TypeMeta: api.TypeMeta{
					Kind: api.VirtualIPKind,
				},
				Status: api.Status{
					Code:    errors.SNAT_NO_DATA,
					Message: "SNAT_NO_DATA",
				},
			}
			return rendererFactory.RenderObject("no interface has virtual ip configured", os.Stdout, &noVipFound)
		}
		for _, vip := range virtualIPs {
			err = rendererFactory.RenderObject("", os.Stdout, vip)
			if err != nil {
				return err
			}
		}
		return nil
	}

	virtualIP, err := client.GetVirtualIP(ctx, opts.InterfaceID)
	if err != nil && virtualIP.Status.Code == 0 {
		return fmt.Errorf("error getting virtual ip: %w", err)
	}

	return rendererFactory.RenderObject("", os.Stdout, virtualIP)
}
