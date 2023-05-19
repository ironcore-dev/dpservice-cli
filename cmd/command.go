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
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	dpdkClientOptions := &DPDKClientOptions{}
	rendererOptions := &RendererOptions{}

	cmd := &cobra.Command{
		Use:           "dpservice-cli",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          SubcommandRequired,
	}

	rendererOptions.AddFlags(cmd.PersistentFlags())
	dpdkClientOptions.AddFlags(cmd.PersistentFlags())

	cmd.AddCommand(
		Add(dpdkClientOptions),
		Get(dpdkClientOptions),
		List(dpdkClientOptions),
		Delete(dpdkClientOptions),
		Initialized(dpdkClientOptions),
		Init(dpdkClientOptions),
		completionCmd,
	)

	return cmd
}
