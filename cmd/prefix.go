package cmd

import (
	"github.com/spf13/cobra"
)

// PrefixCmd represents the prefix command
var prefixCmd = &cobra.Command{
	Use: "prefix",
}

func init() {
	rootCmd.AddCommand(prefixCmd)
}
