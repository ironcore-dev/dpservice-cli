package cmd

import (
	"github.com/spf13/cobra"
)

// MachineCmd represents the machine command
var machineCmd = &cobra.Command{
	Use: "machine",
}

func init() {
	rootCmd.AddCommand(machineCmd)
}