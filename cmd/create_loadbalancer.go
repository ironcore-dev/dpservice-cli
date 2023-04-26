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
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func CreateLoadBalancer(dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory) *cobra.Command {
	var (
		opts CreateLoadBalancerOptions
	)

	cmd := &cobra.Command{
		Use:     "loadbalancer <id> --vni <vni> --vip <vip> --lbports [ports]",
		Short:   "Create a loadbalancer",
		Aliases: LoadBalancerAliases,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			loadbalancerID := args[0]
			return RunCreateLoadBalancer(
				cmd.Context(),
				dpdkClientFactory,
				rendererFactory,
				loadbalancerID,
				opts,
			)
		},
	}

	opts.AddFlags(cmd.Flags())

	util.Must(opts.MarkRequiredFlags(cmd))

	return cmd
}

type CreateLoadBalancerOptions struct {
	VNI     uint32
	LbVipIP netip.Addr
	Lbports []string
}

func (o *CreateLoadBalancerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.Uint32Var(&o.VNI, "vni", o.VNI, "VNI to add the loadbalancer to.")
	//fs.IP("vip", net.IP(o.LbVipIP.AsSlice()), "VIP to assign to the loadbalancer.")
	flag.AddrVar(fs, &o.LbVipIP, "vip", o.LbVipIP, "VIP to assign to the loadbalancer.")
	fs.StringSliceVar(&o.Lbports, "lbports", o.Lbports, "LB ports to assign to the loadbalancer")

}

func (o *CreateLoadBalancerOptions) MarkRequiredFlags(cmd *cobra.Command) error {
	for _, name := range []string{"vni", "vip", "lbports"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return err
		}
	}
	return nil
}

func RunCreateLoadBalancer(ctx context.Context, dpdkClientFactory DPDKClientFactory, rendererFactory RendererFactory, loadbalancerID string, opts CreateLoadBalancerOptions) error {
	client, cleanup, err := dpdkClientFactory.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating dpdk client: %w", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			fmt.Printf("Error cleaning up client: %v\n", err)
		}
	}()

	renderer, err := rendererFactory.NewRenderer("created", os.Stdout)
	if err != nil {
		return fmt.Errorf("error creating renderer: %w", err)
	}

	var ports = make([]api.LBPort, 0)
	for _, p := range opts.Lbports {
		port, err := api.StringLbportToLbport(p)
		if err != nil {
			return err
		}
		ports = append(ports, port)
	}

	lb, err := client.CreateLoadBalancer(ctx, &api.LoadBalancer{
		LoadBalancerMeta: api.LoadBalancerMeta{
			ID: loadbalancerID,
		},
		Spec: api.LoadBalancerSpec{
			VNI:     opts.VNI,
			LbVipIP: opts.LbVipIP,
			Lbports: ports,
		},
	})
	fmt.Println(lb)
	if err != nil {
		return fmt.Errorf("error creating loadbalancer: %w", err)
	}

	if err := renderer.Render(lb); err != nil {
		return fmt.Errorf("error rendering loadbalancer: %w", err)
	}
	return nil
}
