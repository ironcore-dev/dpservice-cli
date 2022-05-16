package cmd

import (
	"github.com/spf13/cobra"
)

// vipCmd represents the vip command
var vipCmd = &cobra.Command{
	Use: "vip",
}

func init() {
	rootCmd.AddCommand(vipCmd)
}
