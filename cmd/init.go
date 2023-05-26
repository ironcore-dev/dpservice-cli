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

	"github.com/onmetal/dpservice-cli/dpdk/api/errors"
	"github.com/onmetal/dpservice-cli/util"
	dpdkproto "github.com/onmetal/net-dpservice-go/proto"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// func Init is not up to dpdk.proto spec, but is implemented to comply with current dpservice implementation
func Init(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts InitOptions
	)

	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Initial set up of the DPDK app",
		Example: "dpservice-cli init",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunInit(
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

type InitOptions struct {
}

func (o *InitOptions) AddFlags(fs *pflag.FlagSet) {
}

func (o *InitOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	return nil
}

func RunInit(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
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
	if err == nil && uuid != "" {
		return fmt.Errorf("error dp-service already initialized, uuid: %s", uuid)
	}

	init, err := client.Init(ctx, dpdkproto.InitConfig{})
	if err != nil && err != errors.ErrServerError {
		return fmt.Errorf("error init: %w", err)
	}

	return rendererFactory.RenderObject("", os.Stdout, init)
}
