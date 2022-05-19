package cmd

import (
	"context"
	"fmt"
	"github.com/onmetal/net-dpservice-go/proto"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// addRouteCmd represents the add command
var addRouteCmd = &cobra.Command{
	Use: "add",
	Run: func(cmd *cobra.Command, args []string) {
		client, closer := getDpClient(cmd)
		defer closer.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		ipv4, err := cmd.Flags().GetIP("ipv4")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}
		ipv6, err := cmd.Flags().GetIP("ipv6")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}
		vni, err := cmd.Flags().GetUint32("vni")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}

		t_vni, err := cmd.Flags().GetUint32("t_vni")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}
		t_ipv6, err := cmd.Flags().GetIP("t_ipv6")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}

		weight, err := cmd.Flags().GetUint32("weight")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}

		length, err := cmd.Flags().GetUint32("length")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}

		prefix := &dpdkproto.Prefix{
			PrefixLength: length,
		}

		if ipv4.String() != "" {
			prefix.IpVersion = dpdkproto.IPVersion_IPv4
			prefix.Address = ipv4
		} else {
			prefix.IpVersion = dpdkproto.IPVersion_IPv6
			prefix.Address = ipv6
		}
		req := &dpdkproto.VNIRouteMsg{
			Vni: &dpdkproto.VNIMsg{Vni: vni},
			Route: &dpdkproto.Route{
				IpVersion:      dpdkproto.IPVersion_IPv6,
				Weight:         weight,
				Prefix:         prefix,
				NexthopVNI:     t_vni,
				NexthopAddress: t_ipv6,
			},
		}

		msg, err := client.AddRoute(ctx, req)
		if err != nil {
			panic(err)
		}
		fmt.Println("AddRoute", msg, req)
	},
}

func init() {
	routeCmd.AddCommand(addRouteCmd)

	addRouteCmd.Flags().Uint32("vni", 0, "")
	addRouteCmd.Flags().Uint32("length", 0, "")
	addRouteCmd.Flags().IP("ipv4", net.IP{}, "")
	addRouteCmd.Flags().IP("ipv6", net.IP{}, "")

	addRouteCmd.Flags().Uint32("weight", 100, "")
	addRouteCmd.Flags().IP("t_vni", net.IP{}, "")
	addRouteCmd.Flags().IP("t_ipv6", net.IP{}, "")

	_ = addRouteCmd.MarkFlagRequired("vni")
	_ = addRouteCmd.MarkFlagRequired("length")
	_ = addRouteCmd.MarkFlagRequired("t_vni")
	_ = addRouteCmd.MarkFlagRequired("t_ipv6")
}