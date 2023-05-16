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

	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func DeleteLoadBalancer(factory DPDKClientFactory) *cobra.Command {
	var (
		opts DeleteLoadBalancerOptions
	)

	cmd := &cobra.Command{
		Use:     "loadbalancer <--id>",
		Short:   "Delete loadbalancer",
		Example: "dpservice-cli delete loadbalancer --id=1",
		Aliases: LoadBalancerAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunDeleteLoadBalancer(cmd.Context(), factory, opts)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type DeleteLoadBalancerOptions struct {
	ID string
}

func (o *DeleteLoadBalancerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ID, "id", o.ID, "LoadBalancer ID to delete.")
}

func (o *DeleteLoadBalancerOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunDeleteLoadBalancer(ctx context.Context, factory DPDKClientFactory, opts DeleteLoadBalancerOptions) error {
	client, cleanup, err := factory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	if err := client.DeleteLoadBalancer(ctx, opts.ID); err != nil {
		return fmt.Errorf("error deleting loadbalancer %s: %v", opts.ID, err)
	}

	fmt.Println("Deleted loadbalancer", opts.ID)

	return nil
}
