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

func DeleteLoadBalancerPrefix(factory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts DeleteLoadBalancerPrefixOptions
	)

	cmd := &cobra.Command{
		Use:     "lbprefix <--prefix> <--interface-id>",
		Short:   "Delete a loadbalancer prefix",
		Example: "dpservice-cli delete lbprefix --prefix=ff80::1/64 --interface-id=vm1",
		Aliases: LoadBalancerPrefixAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunDeleteLoadBalancerPrefix(
				cmd.Context(),
				factory,
				rendererFactory,
				opts,
			)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type DeleteLoadBalancerPrefixOptions struct {
	Prefix      netip.Prefix
	InterfaceID string
}

func (o *DeleteLoadBalancerPrefixOptions) AddFlags(fs *pflag.FlagSet) {
	flag.PrefixVar(fs, &o.Prefix, "prefix", o.Prefix, "Loadbalancer prefix to delete.")
	fs.StringVar(&o.InterfaceID, "interface-id", o.InterfaceID, "Interface ID of the loadbalancer prefix.")
}

func (o *DeleteLoadBalancerPrefixOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"prefix", "interface-id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunDeleteLoadBalancerPrefix(ctx context.Context, factory DPDKClientFactory, rendererFactory RendererFactory, opts DeleteLoadBalancerPrefixOptions) error {
	client, cleanup, err := factory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error deleting dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	if err := client.DeleteLoadBalancerPrefix(ctx, opts.InterfaceID, opts.Prefix); err != nil {
		return fmt.Errorf("error deleting loadbalancer prefix %s/%v: %v", opts.InterfaceID, opts.Prefix, err)
	}

	renderer, err := rendererFactory.NewRenderer("deleted", os.Stdout)
	if err != nil {
		return fmt.Errorf("error creating renderer: %w", err)
	}
	lbprefix := api.Prefix{
		TypeMeta: api.TypeMeta{Kind: "LoadBalancerPrefix"},
		PrefixMeta: api.PrefixMeta{
			InterfaceID: opts.InterfaceID,
			Prefix:      opts.Prefix,
		},
		Status: api.Status{
			Message: "Deleted",
		},
	}
	if err := renderer.Render(&lbprefix); err != nil {
		return fmt.Errorf("error rendering loadbalancer prefix: %w", err)
	}
	return nil
}
