package cmd

import (
	"context"
	"fmt"
	"github.com/onmetal/net-dpservice-go/proto"
	"github.com/spf13/cobra"
	"os"
	"time"
)

// listRouteCmd represents the list command
var listRouteCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		client, closer := getDpClient(cmd)
		defer closer.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		vni, err := cmd.Flags().GetUint32("vni")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}

		req := &dpdkproto.VNIMsg{
			Vni: vni,
		}
		msg, err := client.ListRoutes(ctx, req)
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}
		for _, m := range msg.GetRoutes() {
			fmt.Println(m.String())
		}
	},
}

func init() {
	routeCmd.AddCommand(listRouteCmd)

	listRouteCmd.Flags().Uint32("vni", 0, "")
	_ = listRouteCmd.MarkFlagRequired("vni")
}
