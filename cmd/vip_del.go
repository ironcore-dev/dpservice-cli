package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	dpdkproto "github.com/onmetal/net-dpservice-go/proto"
	"github.com/spf13/cobra"
)

// delRouteCmd represents the del command
var delVipCmd = &cobra.Command{
	Use: "del",
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

		msg, err := client.DeleteInterfaceVIP(ctx, req)
		if err != nil {
			panic(err)
		}
		fmt.Println("DelInterfaceVIP", msg, req)
	},
}

func init() {
	vipCmd.AddCommand(delVipCmd)
	delVipCmd.Flags().String("interface_id", "", "")
	_ = delVipCmd.MarkFlagRequired("interface_id")
}
