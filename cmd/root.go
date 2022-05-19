package cmd

import (
	"context"
	"fmt"
	"github.com/onmetal/net-dpservice-go/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"
	"time"
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

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, server, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		fmt.Println("Err:", err)
		os.Exit(1)
	}
	client := dpdkproto.NewDPDKonmetalClient(conn)

	return client, conn
}
