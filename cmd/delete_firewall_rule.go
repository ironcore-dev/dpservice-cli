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
		Use:     "firewallrule [<ruleID> ...] <--interface-id>",
		Short:   "Delete firewall rule(s)",
		Aliases: FirewallRuleAliases,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ruleIDs := args
			return RunDeleteFirewallRule(cmd.Context(), factory, ruleIDs, opts)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type DeleteFirewallRuleOptions struct {
	InrerfaceID string
}

func (o *DeleteFirewallRuleOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.InrerfaceID, "interface-id", o.InrerfaceID, "Intreface ID where to delete firewall rule(s).")
}

func (o *DeleteFirewallRuleOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"interface-id"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunDeleteFirewallRule(ctx context.Context, factory DPDKClientFactory, ruleIDs []string, opts DeleteFirewallRuleOptions) error {
	client, cleanup, err := factory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	for _, ruleID := range ruleIDs {
		if err := client.DeleteFirewallRule(ctx, opts.InrerfaceID, ruleID); err != nil {
			fmt.Printf("Error deleting firewall rule %s: %v\n", ruleID, err)
		}

		fmt.Println("Deleted firewall rule", ruleID, "on interface", opts.InrerfaceID)
	}
	return nil
}
