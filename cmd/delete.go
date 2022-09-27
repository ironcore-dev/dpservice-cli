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

	"github.com/onmetal/dpservice-go-library/dpdk/api"
	"github.com/onmetal/dpservice-go-library/dpdk/client/dynamic"
	"github.com/onmetal/dpservice-go-library/sources"
	"github.com/spf13/cobra"
)

func Delete(factory DPDKClientFactory) *cobra.Command {
	sourcesOptions := &SourcesOptions{}

	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return RunDelete(ctx, factory, sourcesOptions)
		},
	}

	sourcesOptions.AddFlags(cmd.Flags())

	subcommands := []*cobra.Command{
		DeleteInterface(factory),
		DeletePrefix(factory),
		DeleteRoute(factory),
		DeleteVirtualIP(factory),
	}

	cmd.Short = fmt.Sprintf("Deletes one of %v", CommandNames(subcommands))
	cmd.Long = fmt.Sprintf("Deletes one of %v", CommandNames(subcommands))

	cmd.AddCommand(subcommands...)

	return cmd
}

func RunDelete(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
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

	iterator, err := sourcesReaderFactory.NewIterator()
	if err != nil {
		return fmt.Errorf("error creating sources iterator: %w", err)
	}

	objs, err := sources.CollectObjects(iterator, api.DefaultScheme)
	if err != nil {
		return fmt.Errorf("error collecting objects: %w", err)
	}

	for _, obj := range objs {
		key := dynamic.ObjectKeyFromObject(obj)

		if err := dc.Delete(ctx, obj); err != nil {
			fmt.Printf("Error deleting %T %s: %v\n", obj, key, err)
		}

		fmt.Printf("Deleted %T %s\n", obj, key)
	}

	return nil
}
