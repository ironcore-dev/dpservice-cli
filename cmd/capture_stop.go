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

func CaptureStop(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "stop",
		Short:   "Stop capturing packets for all interfaces",
		Example: "dpservice-cli capture stop",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunCaptureStop(
				cmd.Context(),
				dpdkClientFactory,
				rendererFactory,
			)
		},
	}
	return cmd
}

func RunCaptureStop(ctx context.Context, dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) error {

	dpdkClient, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}

	defer DpdkClose(cleanup)

	captureStop, err := dpdkClient.CaptureStop(ctx)

	if err != nil && captureStop.Status.Code == 0 {
		return fmt.Errorf("error stopping capturing: %w", err)
	}

	return rendererFactory.RenderObject("Packet capturing stopped \n", os.Stdout, captureStop)
}
