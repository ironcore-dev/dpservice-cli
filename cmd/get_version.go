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
	"strings"

	"github.com/onmetal/dpservice-cli/util"
	"github.com/onmetal/net-dpservice-go/api"
	"github.com/onmetal/net-dpservice-go/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func GetVersion(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts GetVersionOptions
	)

	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Get version of dpservice and protobuf.",
		Example: "dpservice-cli get version",
		Aliases: NatAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunGetVersion(
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

type GetVersionOptions struct {
}

func (o *GetVersionOptions) AddFlags(fs *pflag.FlagSet) {
}

func (o *GetVersionOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	return nil
}

func RunGetVersion(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	opts GetVersionOptions,
) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer DpdkClose(cleanup)

	svcVersion, err := client.GetVersion(ctx, &api.Version{
		TypeMeta: api.TypeMeta{Kind: api.VersionKind},
		VersionMeta: api.VersionMeta{
			ClientName: "dpservice-cli",
			ClientVer:  util.BuildVersion,
		},
	})
	if err != nil && !strings.Contains(err.Error(), errors.StatusErrorString) {
		return fmt.Errorf("error getting version: %w", err)
	}
	return rendererFactory.RenderObject("", os.Stdout, svcVersion)
}
