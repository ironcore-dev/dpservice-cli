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

	"github.com/onmetal/dpservice-go-library/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func GetPrefix(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts GetPrefixOptions
	)

	cmd := &cobra.Command{
		Use:   "prefix [<prefix>...]",
		Short: "Get or list prefix(es)",
		RunE: func(cmd *cobra.Command, args []string) error {
			prefixes, err := ParsePrefixArgs(args)
			if err != nil {
				return err
			}

			return RunGetPrefix(
				cmd.Context(),
				dpdkClientFactory,
				rendererFactory,
				prefixes,
				opts,
			)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type GetPrefixOptions struct {
	InterfaceID string
}

func (o *GetPrefixOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.InterfaceID, "interface-id", o.InterfaceID, "Interface ID of the prefix.")
}

func (o *GetPrefixOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"interface-id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunGetPrefix(
	ctx context.Context,
	factory DPDKClientFactory,
	rendererFactory RendererFactory,
	prefixes []netip.Prefix,
	opts GetPrefixOptions,
) error {
	client, cleanup, err := factory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
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

	if len(prefixes) == 0 {
		prefixList, err := client.ListPrefixes(ctx, opts.InterfaceID)
		if err != nil {
			return fmt.Errorf("error listing prefixes: %w", err)
		}

		if err := renderer.Render(prefixList); err != nil {
			return fmt.Errorf("error rendering list: %w", err)
		}
		return nil
	}

	return fmt.Errorf("getting individual prefixes is not implemented")
}
