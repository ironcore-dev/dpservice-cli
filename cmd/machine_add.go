package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	dpdkproto "github.com/onmetal/net-dpservice-go/proto"
	"github.com/spf13/cobra"
)

// addInterfaceCmd represents the machine add  command
var addInterfaceCmd = &cobra.Command{
	Use: "create",
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

		pci_name, err := cmd.Flags().GetString("pci_name")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}

		machinId, err := cmd.Flags().GetString("interface_id")
		if err != nil {
			fmt.Println("Err:", err)
			os.Exit(1)
		}

		req := &dpdkproto.CreateInterfaceRequest{
			InterfaceType: dpdkproto.InterfaceType_VirtualInterface,
			InterfaceID:   []byte(machinId),
			DeviceName:    pci_name,
			Vni:           vni,
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

		msg, err := client.CreateInterface(ctx, req)
		if err != nil {
			panic(err)
		}
		fmt.Println("createinterface", msg, req)

	},
}

func init() {
	machineCmd.AddCommand(addInterfaceCmd)

	addInterfaceCmd.Flags().Uint32("vni", 0, "")
	addInterfaceCmd.Flags().StringP("interface_id", "i", "", "")
	addInterfaceCmd.Flags().String("ipv4", "", "")
	addInterfaceCmd.Flags().String("ipv6", "", "")
	addInterfaceCmd.Flags().String("pci_name", "", "")

	_ = addInterfaceCmd.MarkFlagRequired("pci_name")
	_ = addInterfaceCmd.MarkFlagRequired("vni")
	_ = addInterfaceCmd.MarkFlagRequired("interface_id")
}
