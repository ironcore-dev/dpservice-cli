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
	"github.com/onmetal/dpservice-cli/dpdk/api/errors"
	"github.com/onmetal/dpservice-cli/flag"
	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func AddPrefix(
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
) *cobra.Command {
	var (
		opts AddPrefixOptions
	)

	cmd := &cobra.Command{
		Use:     "prefix <--prefix> <--interface-id>",
		Short:   "Add a prefix to interface.",
		Example: "dpservice-cli add prefix --prefix=10.20.30.0/24 --interface-id=vm1",
		Args:    cobra.ExactArgs(0),
		Aliases: PrefixAliases,
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunAddPrefix(
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

type AddPrefixOptions struct {
	Prefix      netip.Prefix
	InterfaceID string
}

func (o *AddPrefixOptions) AddFlags(fs *pflag.FlagSet) {
	flag.PrefixVar(fs, &o.Prefix, "prefix", o.Prefix, "Prefix to add to the interface.")
	fs.StringVar(&o.InterfaceID, "interface-id", o.InterfaceID, "ID of the interface to add the prefix for.")
}

func (o *AddPrefixOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"prefix", "interface-id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunAddPrefix(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	opts AddPrefixOptions,
) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer DpdkClose(cleanup)

	prefix, err := client.AddPrefix(ctx, &api.Prefix{
		PrefixMeta: api.PrefixMeta{
			InterfaceID: opts.InterfaceID,
		},
		Spec: api.PrefixSpec{
			Prefix: opts.Prefix,
		},
	})
	if err != nil && err != errors.ErrServerError {
		return fmt.Errorf("error adding prefix: %w", err)
	}

	return rendererFactory.RenderObject(fmt.Sprintf("added, underlay route: %s", prefix.Spec.UnderlayRoute), os.Stdout, prefix)
}
