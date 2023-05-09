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

	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func GetNatInfo(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts GetNatInfoOptions
	)

	cmd := &cobra.Command{
		Use:   "natinfo <NatIP> <--nat-type>",
		Short: "List all local machines that are behind this IP",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			natVIPIP, err := netip.ParseAddr(args[0])
			if err != nil {
				return fmt.Errorf("error parsing nat vip ip: %w", err)
			}

			return RunGetNatInfo(
				cmd.Context(),
				dpdkClientFactory,
				rendererFactory,
				natVIPIP,
				opts,
			)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type GetNatInfoOptions struct {
	NatType int32
}

func (o *GetNatInfoOptions) AddFlags(fs *pflag.FlagSet) {
	fs.Int32Var(&o.NatType, "nat-type", o.NatType, "NAT Info type: NATInfoTypeZero = 0/NATInfoLocal = 1/NATInfoNeigh = 2")
}

func (o *GetNatInfoOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"nat-type"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunGetNatInfo(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	natVIPIP netip.Addr,
	opts GetNatInfoOptions,
) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error getting dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	renderer, err := rendererFactory.NewRenderer("", os.Stdout)
	if err != nil {
		return fmt.Errorf("error creating renderer: %w", err)
	}

	natinfo, err := client.GetNATInfo(ctx, natVIPIP, opts.NatType)
	if err != nil {
		return fmt.Errorf("error getting nat info for interface %s: %v", natVIPIP, err)
	}

	if err := renderer.Render(natinfo); err != nil {
		return fmt.Errorf("error rendering nat info: %w", err)
	}

	return nil
}
