package cmd

import (
	"context"
	"fmt"
	"github.com/onmetal/net-dpservice-go/proto"

	"time"

	"github.com/spf13/cobra"
)

// listMachineCmd represents the machine list command
var listMachineCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		client, closer := getDpClient(cmd)
		defer closer.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		msg, err := client.ListMachines(ctx, &dpdkproto.Empty{})
		if err != nil {
			panic(err)
		}
		for _, m := range msg.GetMachines() {
			fmt.Println(m.String())
		}
	},
}

func init() {
	machineCmd.AddCommand(listMachineCmd)
}
