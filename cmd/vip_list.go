package cmd

import (
	"context"
	"fmt"
	"github.com/onmetal/dpservice-go-library/pkg/dpdkproto"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// listVipCmd represents the list command
var listVipCmd = &cobra.Command{
	Use: "list",
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

		msg, err := client.GetMachineVIP(ctx, req)
		if err != nil {
			panic(err)
		}
		fmt.Println("GetMachineVIP", msg, req)
	},
}

func init() {
	vipCmd.AddCommand(listVipCmd)
	listVipCmd.Flags().String("machine_id", "", "")
	_ = listVipCmd.MarkFlagRequired("machine_id")
}
