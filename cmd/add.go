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
	"os"
	"reflect"

	"github.com/onmetal/dpservice-cli/dpdk/api"
	"github.com/onmetal/dpservice-cli/dpdk/api/errors"
	"github.com/onmetal/dpservice-cli/dpdk/client/dynamic"
	"github.com/onmetal/dpservice-cli/sources"
	"github.com/spf13/cobra"
)

func Add(factory DPDKClientFactory) *cobra.Command {
	rendererOptions := &RendererOptions{Output: "name"}
	sourcesOptions := &SourcesOptions{}

	cmd := &cobra.Command{
		Use:  "add",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return RunAdd(ctx, factory, rendererOptions, sourcesOptions)
		},
	}

	rendererOptions.AddFlags(cmd.PersistentFlags())

	sourcesOptions.AddFlags(cmd.Flags())

	subcommands := []*cobra.Command{
		CreateInterface(factory, rendererOptions),
		AddPrefix(factory, rendererOptions),
		AddRoute(factory, rendererOptions),
		AddVirtualIP(factory, rendererOptions),
		CreateLoadBalancer(factory, rendererOptions),
		CreateLoadBalancerPrefix(factory, rendererOptions),
		AddLoadBalancerTarget(factory, rendererOptions),
		AddNat(factory, rendererOptions),
		AddNeighborNat(factory, rendererOptions),
		AddFirewallRule(factory, rendererOptions),
	}

	cmd.Short = fmt.Sprintf("Creates one of %v", CommandNames(subcommands))
	cmd.Long = fmt.Sprintf("Creates one of %v", CommandNames(subcommands))

	cmd.AddCommand(
		subcommands...,
	)

	return cmd
}

func RunAdd(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	sourcesReaderFactory SourcesReaderFactory,
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

	dc := dynamic.NewFromStructured(client)

	renderer, err := rendererFactory.NewRenderer("created", os.Stdout)
	if err != nil {
		return fmt.Errorf("error creating renderer: %w", err)
	}

	iterator, err := sourcesReaderFactory.NewIterator()
	if err != nil {
		return fmt.Errorf("error creating sources iterator: %w", err)
	}

	objs, err := sources.CollectObjects(iterator, api.DefaultScheme)
	if err != nil {
		return fmt.Errorf("error collecting objects: %w", err)
	}

	for _, obj := range objs {
		res, err := dc.Create(ctx, obj)
		if err == errors.ErrServerError {
			r := reflect.ValueOf(res)
			err := reflect.Indirect(r).FieldByName("Status").FieldByName("Error")
			msg := reflect.Indirect(r).FieldByName("Status").FieldByName("Message")
			fmt.Printf("Error adding %T: Server error: %v %v\n", res, err, msg)
			continue
		}
		if err != nil {
			fmt.Printf("Error adding %T: %v\n", obj, err)
			continue
		}

		if err := renderer.Render(res); err != nil {
			return fmt.Errorf("error rendering %T: %w", obj, err)
		}
	}

	return nil
}
