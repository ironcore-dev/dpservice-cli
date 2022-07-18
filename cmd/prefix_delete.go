package cmd

import (
	"context"
	"fmt"
	"github.com/onmetal/net-dpservice-go/proto"
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

		length, err := cmd.Flags().GetUint32("length")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}

		prefix := &dpdkproto.Prefix{
			PrefixLength: length,
		}

		if ipv4 != "" {
			prefix.IpVersion = dpdkproto.IPVersion_IPv4
			prefix.Address = []byte(ipv4)
		} else {
			prefix.IpVersion = dpdkproto.IPVersion_IPv6
			prefix.Address = []byte(ipv6)
		}

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
	delPrefixCmd.Flags().Uint32("length", 0, "")
	delPrefixCmd.Flags().String("ipv4", "", "")
	delPrefixCmd.Flags().String("ipv6", "", "")

	_ = delPrefixCmd.MarkFlagRequired("length")
}
