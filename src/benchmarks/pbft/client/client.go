package client

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	pb "github.com/aleksander-vedvik/benchmark/pbft/protos"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	mgr  *pb.Manager
	view *pb.Configuration
	id   int
	addr string
}

// Creates a new StorageClient with the provided srvAddresses as the configuration
func New(id int, addr string, srvAddresses []string, qSize int, logger *slog.Logger) *Client {
	mgr, config := getConfig(id, addr, srvAddresses, qSize, logger)
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

func (c *Client) Benchmark(ctx context.Context) (*pb.Result, error) {
	return c.view.Benchmark(ctx, &emptypb.Empty{})
}

//func (c *Client) Benchmark() (*pb.Result, error) {
////slog.Info(fmt.Sprintf("client(%v): writing", sc.id), "val", value)
//ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
//defer cancel()
//return c.config.Benchmark(ctx, &empty.Empty{})
//}
