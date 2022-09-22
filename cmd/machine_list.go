package cmd

import (
	"context"
	"fmt"

	dpdkproto "github.com/onmetal/net-dpservice-go/proto"

	"time"

	"github.com/spf13/cobra"
)

// listInterfaceCmd represents the machine list command
var listInterfaceCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		client, closer := getDpClient(cmd)
		defer closer.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		msg, err := client.ListInterfaces(ctx, &dpdkproto.Empty{})
		if err != nil {
			panic(err)
		}
		for _, m := range msg.GetInterfaces() {
			fmt.Println(m.String())
		}
	},
}

func init() {
	machineCmd.AddCommand(listInterfaceCmd)
}
