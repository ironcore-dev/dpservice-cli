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
	"strconv"

	"github.com/onmetal/dpservice-cli/dpdk/api"
	"github.com/onmetal/dpservice-cli/dpdk/api/errors"
	"github.com/onmetal/dpservice-cli/renderer"
	"github.com/onmetal/dpservice-cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func DeleteFirewallRule(factory DPDKClientFactory, rendererFactory RendererFactory, errorRenderer renderer.Renderer) *cobra.Command {
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

			return RunDeleteFirewallRule(
				cmd.Context(),
				factory,
				rendererFactory,
				opts,
				errorRenderer,
			)
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

func RunDeleteFirewallRule(ctx context.Context, factory DPDKClientFactory, rendererFactory RendererFactory, opts DeleteFirewallRuleOptions, errorRenderer renderer.Renderer) error {
	client, cleanup, err := factory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	if status, err := client.DeleteFirewallRule(ctx, opts.InterfaceID, opts.RuleID); err != nil {
		if rendErr := errorRenderer.Render(&status); rendErr != nil {
			return fmt.Errorf("error rendering status: %w", rendErr)
		}
		return fmt.Errorf(strconv.Itoa(errors.SERVER_ERROR))
	}

	fwrule := api.FirewallRule{
		TypeMeta: api.TypeMeta{Kind: api.FirewallRuleKind},
		FirewallRuleMeta: api.FirewallRuleMeta{
			InterfaceID: opts.InterfaceID,
			RuleID:      opts.RuleID,
		},
		Status: api.Status{
			Message: "Deleted",
		},
	}

	renderer, err := rendererFactory.NewRenderer("deleted", os.Stdout)
	if err != nil {
		return fmt.Errorf("error creating renderer: %w", err)
	}
	if err := renderer.Render(&fwrule); err != nil {
		return fmt.Errorf("error rendering prefix: %w", err)
	}

	return nil
}
