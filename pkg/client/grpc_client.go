package client

import (
	"context"
	"github.com/onmetal/dpservice-go-library/pkg/dpdkproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"time"
)

func New(server string) (dpdkproto.DPDKonmetalClient, io.Closer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, server, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, nil, err
	}
	client := dpdkproto.NewDPDKonmetalClient(conn)
	return client, conn, nil
}
