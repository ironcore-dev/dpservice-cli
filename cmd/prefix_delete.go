package cmd

import (
	"context"
	"fmt"
	"github.com/onmetal/net-dpservice-go/proto"
	"net/netip"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// delPrefixCmd represents the prefix del command
var delPrefixCmd = &cobra.Command{
	Use: "del",
	Run: func(cmd *cobra.Command, args []string) {
		client, closer := getDpClient(cmd)
		defer closer.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		machinId, err := cmd.Flags().GetString("machine_id")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}

		ipv4, err := cmd.Flags().GetString("ipv4")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}
		ipv6, err := cmd.Flags().GetString("ipv6")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}

		prefix := &dpdkproto.Prefix{}
		var ipPrefix netip.Prefix
		if ipv4 != "" {
			ipPrefix, err = netip.ParsePrefix(ipv4)
			prefix.IpVersion = dpdkproto.IPVersion_IPv4
		} else {
			ipPrefix, err = netip.ParsePrefix(ipv6)
			prefix.IpVersion = dpdkproto.IPVersion_IPv6
		}

		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}

		prefix.Address = []byte(ipPrefix.Addr().String())
		prefix.PrefixLength = uint32(ipPrefix.Bits())

		req := &dpdkproto.MachinePrefixMsg{
			MachineId: &dpdkproto.MachineIDMsg{
				MachineID: []byte(machinId),
			},
			Prefix: prefix,
		}

		msg, err := client.DeleteMachinePrefix(ctx, req)
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}
		fmt.Println("DeleteMachinePrefix", msg)
	},
}

func init() {
	prefixCmd.AddCommand(delPrefixCmd)
	delPrefixCmd.Flags().String("ipv4", "", "192.168.1.1/32")
	delPrefixCmd.Flags().String("ipv6", "", "")

	delPrefixCmd.Flags().StringP("machine_id", "m", "", "")
	_ = delPrefixCmd.MarkFlagRequired("machine_id")
}
