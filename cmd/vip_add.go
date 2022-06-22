package cmd

import (
	"context"
	"fmt"
	"github.com/onmetal/net-dpservice-go/proto"
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

		vipIp := &dpdkproto.MachineVIPIP{}

		if ipv4 != "" {
			vipIp.IpVersion = dpdkproto.IPVersion_IPv4
			vipIp.Address = []byte(ipv4)
		} 

		machinId, err := cmd.Flags().GetString("machine_id")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}

		req := &dpdkproto.MachineVIPMsg{
			MachineID:    []byte(machinId),
			MachineVIPIP: vipIp,
		}

		msg, err := client.AddMachineVIP(ctx, req)
		if err != nil {
			panic(err)
		}
		fmt.Println("AddMachineVIP, status: %d", msg, req, msg.Error)
	},
}

func init() {
	vipCmd.AddCommand(addVipCmd)

	addVipCmd.Flags().String("machine_id", "", "")
	addVipCmd.Flags().String("ipv4", "", "")

	_ = addVipCmd.MarkFlagRequired("machine_id")
}
