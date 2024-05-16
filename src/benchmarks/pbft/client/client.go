package client

import (
	"context"
	"strconv"
	"time"

	pb "github.com/aleksander-vedvik/benchmark/pbft/protos"
)

type Client struct {
	mgr  *pb.Manager
	view *pb.Configuration
	id   int
	addr string
}

// Creates a new StorageClient with the provided srvAddresses as the configuration
func New(addr string, srvAddresses []string, qSize int) *Client {
	mgr, config := getConfig(addr, srvAddresses, qSize)
	return &Client{
		mgr:  mgr,
		view: config,
		id:   0,
		addr: addr,
	}
}

func (c *Client) Stop() {
	c.mgr.Close()
}

func (c *Client) Write(ctx context.Context, value string) (*pb.ClientResponse, error) {
	c.id++
	req := &pb.WriteRequest{
		Id:        strconv.Itoa(c.id),
		From:      "client",
		Message:   value,
		Timestamp: time.Now().Unix(),
	}
	return c.view.Write(ctx, req)
}

//func (c *Client) Benchmark() (*pb.Result, error) {
////slog.Info(fmt.Sprintf("client(%v): writing", sc.id), "val", value)
//ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
//defer cancel()
//return c.config.Benchmark(ctx, &empty.Empty{})
//}
