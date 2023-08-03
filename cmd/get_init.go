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

	"github.com/spf13/cobra"
)

func GetInit(factory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Indicates if the DPDK app has been initialized already",
		Example: "dpservice-cli get init",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunGetInit(
				cmd.Context(),
				factory,
				rendererFactory,
			)
		},
	}

	return cmd
}

func RunGetInit(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
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

	init, err := client.CheckInitialized(ctx)
	if err != nil && init.Status.Code == 0 {
		return fmt.Errorf("error checking initialization status: %w", err)
	}

	return rendererFactory.RenderObject("", os.Stdout, init)
}
