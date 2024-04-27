package client

import (
	"context"
	"strconv"
	"time"

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
func New(addr string, srvAddresses []string, qSize int) *Client {
	mgr, config := getConfig(addr, srvAddresses, qSize)
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return c.config.Benchmark(ctx, &empty.Empty{})
}
