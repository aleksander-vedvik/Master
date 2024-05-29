package client

import (
	"context"
	"log/slog"
	"strconv"

	pb "github.com/aleksander-vedvik/benchmark/simple/protos"
	"github.com/golang/protobuf/ptypes/empty"
)

type Client struct {
	mgr    *pb.Manager
	config *pb.Configuration
	id     int
	addr   string
}

// Creates a new StorageClient with the provided srvAddresses as the configuration
func New(id int, addr string, srvAddresses []string, qSize int, logger *slog.Logger) *Client {
	mgr, config := getConfig(id, addr, srvAddresses, qSize, logger)
	return &Client{
		mgr:    mgr,
		config: config,
		id:     0,
		addr:   addr,
	}
}

func (c *Client) Stop() {
	c.mgr.Close()
}

func (c *Client) Write1(ctx context.Context, value string) (*pb.WriteResponse1, error) {
	req := &pb.WriteRequest1{
		Id:      strconv.Itoa(c.id),
		Message: value,
	}
	return c.config.BroadcastCall1(ctx, req)
}

func (c *Client) Write2(value string) (*pb.WriteResponse2, error) {
	ctx := context.Background()
	req := &pb.WriteRequest2{
		Id:   strconv.Itoa(c.id),
		Data: []byte(value),
	}
	return c.config.BroadcastCall2(ctx, req)
}

func (c *Client) Benchmark() (*pb.Result, error) {
	//slog.Info(fmt.Sprintf("client(%v): writing", sc.id), "val", value)
	return c.config.Benchmark(context.Background(), &empty.Empty{})
}
