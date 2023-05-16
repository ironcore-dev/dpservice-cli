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

	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func DeleteFirewallRule(factory DPDKClientFactory) *cobra.Command {
	var (
		opts DeleteFirewallRuleOptions
	)

	cmd := &cobra.Command{
		Use:     "firewallrule <--rule-id> <--interface-id>",
		Short:   "Delete firewall rule from interface",
		Example: "dpservice-cli delete firewallrule --rule-id=1 --interface-id=vm1",
		Aliases: FirewallRuleAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunDeleteFirewallRule(cmd.Context(), factory, opts)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type DeleteFirewallRuleOptions struct {
	RuleID      string
	InterfaceID string
}

func (o *DeleteFirewallRuleOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.RuleID, "rule-id", o.RuleID, "Rule ID to delete.")
	fs.StringVar(&o.InterfaceID, "interface-id", o.InterfaceID, "Intreface ID where to delete firewall rule.")
}

func (o *DeleteFirewallRuleOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"rule-id", "interface-id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunDeleteFirewallRule(ctx context.Context, factory DPDKClientFactory, opts DeleteFirewallRuleOptions) error {
	client, cleanup, err := factory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	if err := client.DeleteFirewallRule(ctx, opts.InterfaceID, opts.RuleID); err != nil {
		return fmt.Errorf("error deleting firewall rule %s/%s: %v", opts.RuleID, opts.InterfaceID, err)
	}

	fmt.Printf("Deleted firewall rule %s on interface %s\n", opts.RuleID, opts.InterfaceID)

	return nil
}
