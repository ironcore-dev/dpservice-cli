package cmd

import (
	"context"
	"fmt"

	dpdkproto "github.com/onmetal/net-dpservice-go/proto"

	// "net"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// addVipCmd represents the add command
var addVipCmd = &cobra.Command{
	Use: "add",
	Run: func(cmd *cobra.Command, args []string) {
		client, closer := getDpClient(cmd)
		defer closer.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		ipv4, err := cmd.Flags().GetString("ipv4")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}

		vipIp := &dpdkproto.InterfaceVIPIP{}

		if ipv4 != "" {
			vipIp.IpVersion = dpdkproto.IPVersion_IPv4
			vipIp.Address = []byte(ipv4)
		}

		machinId, err := cmd.Flags().GetString("interface_id")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}

		req := &dpdkproto.InterfaceVIPMsg{
			InterfaceID:    []byte(machinId),
			InterfaceVIPIP: vipIp,
		}

		msg, err := client.AddInterfaceVIP(ctx, req)
		if err != nil {
			panic(err)
		}
		fmt.Println("AddInterfaceVIP: ", msg)
	},
}

func init() {
	vipCmd.AddCommand(addVipCmd)

	addVipCmd.Flags().String("interface_id", "", "")
	addVipCmd.Flags().String("ipv4", "", "")

	_ = addVipCmd.MarkFlagRequired("interface_id")
}
