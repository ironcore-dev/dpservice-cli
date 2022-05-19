package cmd

import (
	"context"
	"fmt"
	"github.com/onmetal/net-dpservice-go/proto"
	"github.com/spf13/cobra"
	"os"
	"time"
)

// delRouteCmd represents the del command
var delVipCmd = &cobra.Command{
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

		req := &dpdkproto.MachineIDMsg{
			MachineID: []byte(machinId),
		}

		msg, err := client.DelMachineVIP(ctx, req)
		if err != nil {
			panic(err)
		}
		fmt.Println("DelMachineVIP", msg, req)
	},
}

func init() {
	vipCmd.AddCommand(delVipCmd)
	delVipCmd.Flags().String("machine_id", "", "")
	_ = delVipCmd.MarkFlagRequired("machine_id")
}
