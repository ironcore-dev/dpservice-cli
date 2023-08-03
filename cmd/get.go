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
	"fmt"

	"github.com/spf13/cobra"
)

func Get(factory DPDKClientFactory) *cobra.Command {
	rendererOptions := &RendererOptions{Output: "table"}

	cmd := &cobra.Command{
		Use:  "get",
		Args: cobra.NoArgs,
		RunE: SubcommandRequired,
	}

	rendererOptions.AddFlags(cmd.PersistentFlags())

	subcommands := []*cobra.Command{
		GetInterface(factory, rendererOptions),
		GetVirtualIP(factory, rendererOptions),
		GetLoadBalancer(factory, rendererOptions),
		GetNat(factory, rendererOptions),
		GetFirewallRule(factory, rendererOptions),
		GetVni(factory, rendererOptions),
		GetVersion(factory, rendererOptions),
		GetInit(factory, rendererOptions),
	}

	cmd.Short = fmt.Sprintf("Gets one of %v", CommandNames(subcommands))
	cmd.Long = fmt.Sprintf("Gets one of %v", CommandNames(subcommands))

	cmd.AddCommand(
		subcommands...,
	)

	return cmd
}
