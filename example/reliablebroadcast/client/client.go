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
	mgr, config := createConfiguration(addr, srvAddresses, qSize)
	return &Client{
		mgr:    mgr,
		config: config,
	}
}

func createConfiguration(addr string, srvAddresses []string, qSize int) (*pb.Manager, *pb.Configuration) {
	mgr := pb.NewManager(
		gorums.WithGrpcDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	lis, _ := net.Listen("tcp", addr)
	mgr.AddClientServer(lis, lis.Addr())
	configuration, _ := mgr.NewConfiguration(
		NewQSpec(qSize),
		gorums.WithNodeList(srvAddresses),
	)
	return mgr, configuration
}

func (c *Client) Stop() {
	c.mgr.Close()
}

func (c *Client) Broadcast(value string) (*pb.Message, error) {
	ctx := context.Background()
	req := &pb.Message{
		Data: value,
	}
	return c.config.Broadcast(ctx, req)
}
