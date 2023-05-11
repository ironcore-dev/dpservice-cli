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
	"strings"

	"github.com/onmetal/dpservice-cli/dpdk/api"
	"github.com/onmetal/dpservice-cli/util"
	dpdkproto "github.com/onmetal/net-dpservice-go/proto"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Init(factory DPDKClientFactory) *cobra.Command {
	var (
		opts InitOptions
	)

	cmd := &cobra.Command{
		Use:     "init <underlayIPv6Prefix> [flags]",
		Short:   "Initial set up of the DPDK app",
		Long:    "To add multiple values to flags, use \",\" (comma) between values",
		Example: "dpservice-cli init ff80::1/64 --pfnames=a,b,c --uplink-ports=e,f,g",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			prefix, err := netip.ParsePrefix(args[0])
			if err != nil {
				return fmt.Errorf("error parsing prefix: %w", err)
			}

			return RunInit(
				cmd.Context(),
				factory,
				prefix,
				opts,
			)
		},
	}
	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type InitOptions struct {
	UplinkPorts string
	PfNames     string
}

func (o *InitOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.UplinkPorts, "uplink-ports", o.UplinkPorts, "Linux name of the NICs that are connected to the Leaf Switches.")
	fs.StringVar(&o.PfNames, "pfnames", o.PfNames, "Linux name of the Physical Functions, that Virtual Functions will be derived from.")
}

func (o *InitOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"uplink-ports", "pfnames"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunInit(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	prefix netip.Prefix,
	opts InitOptions,
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

	uuid, err := client.Initialized(ctx)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	if uuid != "" {
		return fmt.Errorf("error dp-service already initialized, uuid: %s", uuid)
	}

	uplinkPorts := strings.Split(opts.UplinkPorts, ",")
	pfNames := strings.Split(opts.PfNames, ",")

	err = client.Init(ctx, dpdkproto.InitConfig{
		UnderlayIPv6Prefix: &dpdkproto.Prefix{
			IpVersion:    api.NetIPAddrToProtoIPVersion(prefix.Addr()),
			Address:      []byte(prefix.Addr().String()),
			PrefixLength: uint32(prefix.Bits()),
		},
		UplinkPorts: uplinkPorts,
		PfNames:     pfNames,
	})
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	return nil
}
