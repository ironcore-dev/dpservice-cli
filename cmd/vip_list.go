package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	dpdkproto "github.com/onmetal/net-dpservice-go/proto"
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

		machinId, err := cmd.Flags().GetString("interface_id")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}

		req := &dpdkproto.InterfaceIDMsg{
			InterfaceID: []byte(machinId),
		}

		msg, err := client.GetInterfaceVIP(ctx, req)
		if err != nil {
			panic(err)
		}
		fmt.Println("GetInterfaceVIP", msg, req)
	},
}

func init() {
	vipCmd.AddCommand(listVipCmd)
	listVipCmd.Flags().String("interface_id", "", "")
	_ = listVipCmd.MarkFlagRequired("interface_id")
}
