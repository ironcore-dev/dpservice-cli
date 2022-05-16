package cmd

import (
	"context"
	"fmt"
	"github.com/onmetal/dpservice-go-library/pkg/dpdkproto"
	"github.com/spf13/cobra"
	"net"
	"os"
	"time"
)

// delRouteCmd represents the del command
var delRouteCmd = &cobra.Command{
	Use: "del",
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

		msg, err := client.DeleteRoute(ctx, req)
		if err != nil {
			panic(err)
		}
		fmt.Println("DeleteRoute", msg, req)
	},
}

func init() {
	routeCmd.AddCommand(delRouteCmd)

	delRouteCmd.Flags().Uint32("vni", 0, "")
	delRouteCmd.Flags().Uint32("length", 0, "")
	delRouteCmd.Flags().IP("ipv4", net.IP{}, "")
	delRouteCmd.Flags().IP("ipv6", net.IP{}, "")

	delRouteCmd.Flags().Uint32("weight", 100, "")
	delRouteCmd.Flags().IP("t_vni", net.IP{}, "")
	delRouteCmd.Flags().IP("t_ipv6", net.IP{}, "")

	_ = delRouteCmd.MarkFlagRequired("vni")
	_ = delRouteCmd.MarkFlagRequired("length")
	_ = delRouteCmd.MarkFlagRequired("t_vni")
	_ = delRouteCmd.MarkFlagRequired("t_ipv6")

}
