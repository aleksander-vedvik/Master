package client

import (
	"context"
	"strconv"
	"time"

	"github.com/aleksander-vedvik/benchmark/pbft.s/config"
	pb "github.com/aleksander-vedvik/benchmark/pbft.s/protos"
)

type Client struct {
	id     int
	addr   string
	config *config.Config
}

// Creates a new StorageClient with the provided srvAddresses as the configuration
func New(addr string, srvAddresses []string, qSize int) *Client {
	return &Client{
		id:     0,
		addr:   addr,
		config: config.NewConfig(srvAddresses),
	}
}

func (c *Client) Write(value string) (*pb.ClientResponse, error) {
	c.id++
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pb.WriteRequest{
		Id:        strconv.Itoa(c.id),
		From:      "client",
		Message:   value,
		Timestamp: time.Now().Unix(),
	}
	return c.config.Write(ctx, req)
}

//func (c *Client) Benchmark() (*pb.Result, error) {
////slog.Info(fmt.Sprintf("client(%v): writing", sc.id), "val", value)
//ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
//defer cancel()
//return c.config.Benchmark(ctx, &empty.Empty{})
//}
