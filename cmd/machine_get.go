package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	dpdkproto "github.com/onmetal/net-dpservice-go/proto"

	"github.com/spf13/cobra"
)

// getInterfaceCmd represents the interface get command
var getInterfaceCmd = &cobra.Command{
	Use: "get",
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

		msg, err := client.GetInterface(ctx, req)
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}
		fmt.Println(msg.String())
	},
}

func init() {
	machineCmd.AddCommand(getInterfaceCmd)
	getInterfaceCmd.Flags().StringP("interface_id", "i", "", "")
}
