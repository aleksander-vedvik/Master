package client

import (
	"context"
	"net"

	pb "reliablebroadcast/proto"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	mgr    *pb.Manager
	config *pb.Configuration
}

func New(addr string, srvAddresses []string, qSize int) *Client {
	mgr := pb.NewManager(
		gorums.WithGrpcDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	lis, _ := net.Listen("tcp", addr)
	mgr.AddClientServer(lis, lis.Addr())
	config, _ := mgr.NewConfiguration(
		NewQSpec(qSize),
		gorums.WithNodeList(srvAddresses),
	)
	return &Client{
		mgr:    mgr,
		config: config,
	}
}

func (c *Client) Broadcast(ctx context.Context, value string) (*pb.Message, error) {
	req := &pb.Message{
		Data: value,
	}
	return c.config.Broadcast(ctx, req)
}
