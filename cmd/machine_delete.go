package cmd

import (
	"context"
	"fmt"
	"github.com/onmetal/dpservice-go-library/pkg/dpdkproto"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// delMachineCmd represents the machine del command
var delMachineCmd = &cobra.Command{
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

		msg, err := client.DeleteMachine(ctx, req)
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}
		fmt.Println("DeleteMachine", msg)
	},
}

func init() {
	machineCmd.AddCommand(delMachineCmd)
	delMachineCmd.Flags().String("machine_id", "", "")
}
