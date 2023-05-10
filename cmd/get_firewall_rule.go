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

	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func GetFirewallRule(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts GetFirewallRuleOptions
	)

	cmd := &cobra.Command{
		Use:     "firewallrule <interfaceID> <--rule-id>",
		Short:   "Get firewall rule",
		Example: "dpservice-cli get fwrule vm1 --rule-id=1",
		Aliases: FirewallRuleAliases,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			interfaceID := args[0]
			return RunGetFirewallRule(
				cmd.Context(),
				dpdkClientFactory,
				rendererFactory,
				interfaceID,
				opts,
			)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type GetFirewallRuleOptions struct {
	RuleID string
}

func (o *GetFirewallRuleOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.RuleID, "rule-id", o.RuleID, "Rule ID of Firewall Rule.")
}

func (o *GetFirewallRuleOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"rule-id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunGetFirewallRule(
	ctx context.Context,
	dpdkClientFactory DPDKClientFactory,
	rendererFactory RendererFactory,
	interfaceID string,
	opts GetFirewallRuleOptions,
) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error getting dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	renderer, err := rendererFactory.NewRenderer("", os.Stdout)
	if err != nil {
		return fmt.Errorf("error creating renderer: %w", err)
	}

	if len(interfaceID) == 0 {
		return fmt.Errorf("need to specify interface id")
	}

	fwrule, err := client.GetFirewallRule(ctx, interfaceID, opts.RuleID)
	if err != nil {
		return fmt.Errorf("error getting firewall rule: %w", err)
	}

	if err := renderer.Render(fwrule); err != nil {
		return fmt.Errorf("error rendering firewall rule %s/%s: %w", interfaceID, opts.RuleID, err)
	}
	return nil
}
