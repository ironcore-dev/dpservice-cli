package cmd

import (
	"context"
	"fmt"
	"github.com/onmetal/net-dpservice-go/proto"
	"github.com/spf13/cobra"
	"os"
	"time"
)

// addMachineCmd represents the machine add  command
var addMachineCmd = &cobra.Command{
	Use: "add",
	Run: func(cmd *cobra.Command, args []string) {
		client, closer := getDpClient(cmd)
		defer closer.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		ipv4, err := cmd.Flags().GetString("ipv4")
		if err != nil && cmd.Flags().HasFlags() {
			fmt.Println("Err:", err)
			os.Exit(1)
		}
		ipv6, err := cmd.Flags().GetString("ipv6")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}
		vni, err := cmd.Flags().GetUint32("vni")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}

		machinId, err := cmd.Flags().GetString("machine_id")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}

		req := &dpdkproto.AddMachineRequest{
			MachineType: dpdkproto.MachineType_VirtualMachine,
			MachineID:   []byte(machinId),
			Vni:         vni,
		}
		if ipv4 != "" {
			req.Ipv4Config = &dpdkproto.IPConfig{
				IpVersion:      dpdkproto.IPVersion_IPv4,
				PrimaryAddress: []byte(ipv4),
			}
		}

		if ipv6 != "" {
			req.Ipv6Config = &dpdkproto.IPConfig{
				IpVersion:      dpdkproto.IPVersion_IPv6,
				PrimaryAddress: []byte(ipv6),
			}
		}

		msg, err := client.AddMachine(ctx, req)
		if err != nil {
			panic(err)
		}
		fmt.Println("addmachine", msg, req)

	},
}

func init() {
	machineCmd.AddCommand(addMachineCmd)

	addMachineCmd.Flags().Uint32("vni", 0, "")
	addMachineCmd.Flags().StringP("machine_id", "m", "", "")
	addMachineCmd.Flags().String("ipv4", "", "")
	addMachineCmd.Flags().String("ipv6", "", "")

	_ = addMachineCmd.MarkFlagRequired("vni")
	_ = addMachineCmd.MarkFlagRequired("machine_id")
}
