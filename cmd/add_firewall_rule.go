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
	"net/netip"
	"os"

	"github.com/onmetal/dpservice-cli/dpdk/api"
	"github.com/onmetal/dpservice-cli/flag"
	"github.com/onmetal/dpservice-cli/util"
	dpdkproto "github.com/onmetal/net-dpservice-go/proto"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func AddFirewallRule(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts AddFirewallRuleOptions
	)

	cmd := &cobra.Command{
		Use:     "firewallrule <--interface-id> [flags]",
		Short:   "Add a FirewallRule to interface",
		Example: "dpservice-cli add fwrule --interface-id=vm1 --action=1 --direction=1 --dst=5.5.5.0/24 --ipv=0 --priority=100 --rule-id=12 --src=1.1.1.1/32 --protocol=tcp --srcPortLower=1 --srcPortUpper=1000 --dstPortLower=500 --dstPortUpper=600",
		Aliases: FirewallRuleAliases,
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			return RunAddFirewallRule(
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

type AddFirewallRuleOptions struct {
	InterfaceID       string
	RuleID            string
	TrafficDirection  uint8
	FirewallAction    uint8
	Priority          uint32
	IpVersion         uint8
	SourcePrefix      netip.Prefix
	DestinationPrefix netip.Prefix
	ProtocolFilter    string
	SrcPortLower      int32
	SrcPortUpper      int32
	DstPortLower      int32
	DstPortUpper      int32
	IcmpType          int32
	IcmpCode          int32
}

func (o *AddFirewallRuleOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.InterfaceID, "interface-id", o.InterfaceID, "InterfaceID of FW Rule.")
	fs.StringVar(&o.RuleID, "rule-id", o.RuleID, "RuleID of FW Rule.")
	fs.Uint8Var(&o.TrafficDirection, "direction", o.TrafficDirection, "Traffic direction of FW Rule: Ingress = 0/Egress = 1")
	fs.Uint8Var(&o.FirewallAction, "action", o.FirewallAction, "Firewall action: Drop = 0/Accept = 1 // Can be only \"accept\" at the moment.")
	fs.Uint32Var(&o.Priority, "priority", o.Priority, "Priority of FW Rule. // For future use. No effect at the moment.")
	fs.Uint8Var(&o.IpVersion, "ipv", o.IpVersion, "IpVersion of FW Rule IPv4 = 0/IPv6 = 1.")
	flag.PrefixVar(fs, &o.SourcePrefix, "src", o.SourcePrefix, "Source prefix // 0.0.0.0 with prefix length 0 matches all source IPs.")
	flag.PrefixVar(fs, &o.DestinationPrefix, "dst", o.DestinationPrefix, "Destination prefix // 0.0.0.0 with prefix length 0 matches all destination IPs.")
	fs.StringVar(&o.ProtocolFilter, "protocol", o.ProtocolFilter, "Protocol used icmp/tcp/udp // Not defining a protocol filter matches all protocols.")
	fs.Int32Var(&o.SrcPortLower, "srcPortLower", o.SrcPortLower, "Source Ports start // -1 matches all source ports.")
	fs.Int32Var(&o.SrcPortUpper, "srcPortUpper", o.SrcPortUpper, "Source Ports end.")
	fs.Int32Var(&o.DstPortLower, "dstPortLower", o.DstPortLower, "Destination Ports start // -1 matches all destination ports.")
	fs.Int32Var(&o.DstPortUpper, "dstPortUpper", o.DstPortUpper, "Destination Ports end.")
	fs.Int32Var(&o.IcmpType, "icmpType", o.IcmpType, "ICMP type // -1 matches all ICMP Types.")
	fs.Int32Var(&o.IcmpCode, "icmpCode", o.IcmpCode, "ICMP code // -1 matches all ICMP Codes.")

}

func (o *AddFirewallRuleOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	// TODO if protocol is not specified it should match all protocols
	for _, name := range []string{"interface-id", "rule-id", "direction", "action", "ipv", "src", "dst", "protocol"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunAddFirewallRule(ctx context.Context, dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory, opts AddFirewallRuleOptions) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	renderer, err := rendererFactory.NewRenderer("added", os.Stdout)
	if err != nil {
		return fmt.Errorf("error creating renderer: %w", err)
	}

	srcPfx, err := netip.ParsePrefix(opts.SourcePrefix.String())
	if err != nil {
		return fmt.Errorf("error parsing src prefix: %w", err)
	}
	dstPfx, err := netip.ParsePrefix(opts.DestinationPrefix.String())
	if err != nil {
		return fmt.Errorf("error parsing dst prefix: %w", err)
	}

	// TODO add cases if icmp type or code is -1
	var protocolFilter dpdkproto.ProtocolFilter
	switch opts.ProtocolFilter {
	case "icmp":
		protocolFilter.Filter = &dpdkproto.ProtocolFilter_Icmp{Icmp: &dpdkproto.ICMPFilter{
			IcmpType: opts.IcmpType,
			IcmpCode: opts.IcmpCode}}
	case "tcp":
		if opts.SrcPortLower == -1 {
			opts.SrcPortLower = 1
			opts.SrcPortUpper = 65535
		}
		if opts.DstPortLower == -1 {
			opts.DstPortLower = 1
			opts.DstPortUpper = 65535
		}
		protocolFilter.Filter = &dpdkproto.ProtocolFilter_Tcp{Tcp: &dpdkproto.TCPFilter{
			SrcPortLower: opts.SrcPortLower,
			SrcPortUpper: opts.SrcPortUpper,
			DstPortLower: opts.DstPortLower,
			DstPortUpper: opts.DstPortUpper,
		}}
	case "udp":
		if opts.SrcPortLower == -1 {
			opts.SrcPortLower = 1
			opts.SrcPortUpper = 65535
		}
		if opts.DstPortLower == -1 {
			opts.DstPortLower = 1
			opts.DstPortUpper = 65535
		}
		protocolFilter.Filter = &dpdkproto.ProtocolFilter_Udp{Udp: &dpdkproto.UDPFilter{
			SrcPortLower: opts.SrcPortLower,
			SrcPortUpper: opts.SrcPortUpper,
			DstPortLower: opts.DstPortLower,
			DstPortUpper: opts.DstPortUpper,
		}}
	}

	fwrule, err := client.AddFirewallRule(ctx, &api.FirewallRule{
		TypeMeta: api.TypeMeta{Kind: api.FirewallRuleKind},
		FirewallRuleMeta: api.FirewallRuleMeta{
			RuleID:      opts.RuleID,
			InterfaceID: opts.InterfaceID,
		},
		Spec: api.FirewallRuleSpec{
			TrafficDirection:  opts.TrafficDirection,
			FirewallAction:    opts.FirewallAction,
			Priority:          opts.Priority,
			SourcePrefix:      &srcPfx,
			DestinationPrefix: &dstPfx,
			ProtocolFilter: &dpdkproto.ProtocolFilter{
				Filter: protocolFilter.Filter},
		},
	},
	)
	if err != nil {
		return fmt.Errorf("error adding firewall rule: %w", err)
	}

	if err := renderer.Render(fwrule); err != nil {
		return fmt.Errorf("error rendering firewall rule: %w", err)
	}
	return nil
}
