package cmd

import (
	"context"
	"fmt"
	"github.com/onmetal/net-dpservice-go/proto"
	"os"

	"time"

	"github.com/spf13/cobra"
)

// listPrefixCmd represents the machine list command
var listPrefixCmd = &cobra.Command{
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

		msg, err := client.ListMachinePrefixes(ctx, req)
		if err != nil {
			panic(err)
		}
		for _, p := range msg.GetPrefixes() {
			fmt.Println(p.String())
		}
	},
}

func init() {
	prefixCmd.AddCommand(listPrefixCmd)

	listPrefixCmd.Flags().StringP("machine_id", "m", "", "")
	_ = listPrefixCmd.MarkFlagRequired("machine_id")
}
