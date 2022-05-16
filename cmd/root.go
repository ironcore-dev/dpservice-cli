/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	client2 "github.com/onmetal/dpservice-go-library/pkg/client"
	"github.com/onmetal/dpservice-go-library/pkg/dpdkproto"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "dpservice-go-library",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("server", "localhost:1337", "")
	rootCmd.MarkFlagRequired("server")
}

func getDpClient(cmd *cobra.Command) (dpdkproto.DPDKonmetalClient, io.Closer) {
	server, err := cmd.Flags().GetString("server")
	if err != nil {
		fmt.Println("Err:", err)
		os.Exit(1)
	}

	client, closer, err := client2.New(server)
	if err != nil {
		fmt.Println("Err:", err)
		os.Exit(1)
	}
	return client, closer
}
