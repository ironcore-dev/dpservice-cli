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
	"github.com/onmetal/dpservice-cli/dpdk/api/errors"
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
		Example: "dpservice-cli add fwrule --interface-id=vm1 --action=1 --direction=1 --dst=5.5.5.0/24 --ipver=0 --priority=100 --rule-id=12 --src=1.1.1.1/32 --protocol=tcp --src-port-min=1 --src-port-max=1000 --dst-port-min=500 --dst-port-max=600",
		Aliases: FirewallRuleAliases,
		Args:    cobra.ExactArgs(0),
		// if protocol flag is set, require also additional flags
		PreRun: func(cmd *cobra.Command, args []string) {
			filter, _ := cmd.Flags().GetString("protocol")
			switch filter {
			case "icmp", "1":
				cmd.MarkFlagRequired("icmp-type")
				cmd.MarkFlagRequired("icmp-code")
			case "tcp", "6", "udp", "17":
				cmd.MarkFlagRequired("src-port-min")
				if src, _ := cmd.Flags().GetInt32("src-port-min"); src != -1 {
					cmd.MarkFlagRequired("src-port-max")
				}
				cmd.MarkFlagRequired("dst-port-min")
				if dst, _ := cmd.Flags().GetInt32("dst-port-min"); dst != -1 {
					cmd.MarkFlagRequired("dst-port-max")
				}
			}
		},
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
	TrafficDirection  string
	FirewallAction    string
	Priority          uint32
	IpVersion         string
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
	fs.StringVar(&o.TrafficDirection, "direction", o.TrafficDirection, "Traffic direction of FW Rule: Ingress = 0/Egress = 1")
	fs.StringVar(&o.FirewallAction, "action", o.FirewallAction, "Firewall action: drop/deny/0|accept/allow/1 (Can be only \"accept/allow/1\" at the moment).")
	fs.Uint32Var(&o.Priority, "priority", 1000, "Priority of FW Rule. (For future use. No effect at the moment).")
	fs.StringVar(&o.IpVersion, "ipver", o.IpVersion, "IpVersion of FW Rule: IPv4 = 0/IPv6 = 1.")
	flag.PrefixVar(fs, &o.SourcePrefix, "src", o.SourcePrefix, "Source prefix (0.0.0.0 with prefix length 0 matches all source IPs).")
	flag.PrefixVar(fs, &o.DestinationPrefix, "dst", o.DestinationPrefix, "Destination prefix (0.0.0.0 with prefix length 0 matches all destination IPs).")
	fs.StringVar(&o.ProtocolFilter, "protocol", o.ProtocolFilter, "Protocol used icmp/tcp/udp (Not defining a protocol filter matches all protocols).")
	fs.Int32Var(&o.SrcPortLower, "src-port-min", -1, "Source Ports start (-1 matches all source ports).")
	fs.Int32Var(&o.SrcPortUpper, "src-port-max", -1, "Source Ports end.")
	fs.Int32Var(&o.DstPortLower, "dst-port-min", -1, "Destination Ports start (-1 matches all destination ports).")
	fs.Int32Var(&o.DstPortUpper, "dst-port-max", -1, "Destination Ports end.")
	fs.Int32Var(&o.IcmpType, "icmp-type", -1, "ICMP type (-1 matches all ICMP Types).")
	fs.Int32Var(&o.IcmpCode, "icmp-code", -1, "ICMP code (-1 matches all ICMP Codes).")

}

func (o *AddFirewallRuleOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"interface-id", "rule-id", "direction", "action", "ipver", "src", "dst"} {
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
	defer DpdkClose(cleanup)

	srcPfx, err := netip.ParsePrefix(opts.SourcePrefix.String())
	if err != nil {
		return fmt.Errorf("error parsing src prefix: %w", err)
	}
	dstPfx, err := netip.ParsePrefix(opts.DestinationPrefix.String())
	if err != nil {
		return fmt.Errorf("error parsing dst prefix: %w", err)
	}

	var protocolFilter dpdkproto.ProtocolFilter
	switch opts.ProtocolFilter {
	case "icmp", "1":
		protocolFilter.Filter = &dpdkproto.ProtocolFilter_Icmp{Icmp: &dpdkproto.ICMPFilter{
			IcmpType: opts.IcmpType,
			IcmpCode: opts.IcmpCode}}
	case "tcp", "6":
		if opts.SrcPortLower < -1 || opts.SrcPortLower == 0 || opts.SrcPortLower > 65535 ||
			opts.SrcPortUpper < -1 || opts.SrcPortUpper == 0 || opts.SrcPortUpper > 65535 ||
			opts.DstPortLower < -1 || opts.DstPortLower == 0 || opts.DstPortLower > 65535 ||
			opts.DstPortUpper < -1 || opts.DstPortUpper == 0 || opts.DstPortUpper > 65535 {
			return fmt.Errorf("ports can only be -1 or <1,65535>")
		}
		if opts.SrcPortLower > opts.SrcPortUpper || opts.DstPortLower > opts.DstPortUpper {
			return fmt.Errorf("min port must be lower or equal to max port")
		}
		protocolFilter.Filter = &dpdkproto.ProtocolFilter_Tcp{Tcp: &dpdkproto.TCPFilter{
			SrcPortLower: opts.SrcPortLower,
			SrcPortUpper: opts.SrcPortUpper,
			DstPortLower: opts.DstPortLower,
			DstPortUpper: opts.DstPortUpper,
		}}
	case "udp", "17":
		if opts.SrcPortLower < -1 || opts.SrcPortLower == 0 || opts.SrcPortLower > 65535 ||
			opts.SrcPortUpper < -1 || opts.SrcPortUpper == 0 || opts.SrcPortUpper > 65535 ||
			opts.DstPortLower < -1 || opts.DstPortLower == 0 || opts.DstPortLower > 65535 ||
			opts.DstPortUpper < -1 || opts.DstPortUpper == 0 || opts.DstPortUpper > 65535 {
			return fmt.Errorf("ports can only be -1 or <1,65535>")
		}
		if opts.SrcPortLower > opts.SrcPortUpper || opts.DstPortLower > opts.DstPortUpper {
			return fmt.Errorf("min port must be lower or equal to max port")
		}
		protocolFilter.Filter = &dpdkproto.ProtocolFilter_Udp{Udp: &dpdkproto.UDPFilter{
			SrcPortLower: opts.SrcPortLower,
			SrcPortUpper: opts.SrcPortUpper,
			DstPortLower: opts.DstPortLower,
			DstPortUpper: opts.DstPortUpper,
		}}
	case "":
	default:
		return fmt.Errorf("protocol can be only: icmp = 1/tcp = 6/udp = 17")
	}
	if opts.Priority > 65536 {
		return fmt.Errorf("priority can be only: <0,65536")
	}

	fwrule, err := client.AddFirewallRule(ctx, &api.FirewallRule{
		TypeMeta: api.TypeMeta{Kind: api.FirewallRuleKind},
		FirewallRuleMeta: api.FirewallRuleMeta{
			InterfaceID: opts.InterfaceID,
		},
		Spec: api.FirewallRuleSpec{
			RuleID:            opts.RuleID,
			TrafficDirection:  opts.TrafficDirection,
			FirewallAction:    opts.FirewallAction,
			Priority:          opts.Priority,
			IpVersion:         opts.IpVersion,
			SourcePrefix:      &srcPfx,
			DestinationPrefix: &dstPfx,
			ProtocolFilter: &dpdkproto.ProtocolFilter{
				Filter: protocolFilter.Filter},
		},
	})
	if err != nil && err != errors.ErrServerError {
		return fmt.Errorf("error adding firewall rule: %w", err)
	}

	fwrule.TypeMeta.Kind = api.FirewallRuleKind
	fwrule.FirewallRuleMeta.InterfaceID = opts.InterfaceID
	fwrule.Spec.RuleID = opts.RuleID
	return rendererFactory.RenderObject("added", os.Stdout, fwrule)
}
