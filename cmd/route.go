package cmd

import (
	"github.com/spf13/cobra"
)

// routeCmd represents the route command
var routeCmd = &cobra.Command{
	Use: "route",
}

func init() {
	rootCmd.AddCommand(routeCmd)
}
