package client

import (
	"context"
	"net"
	"time"

	"github.com/aleksander-vedvik/benchmark/pbft.s/config"
	pb "github.com/aleksander-vedvik/benchmark/pbft.s/protos"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Client struct {
	pb.PBFTNodeServer
	id     int
	addr   string
	config *config.Config
	srv    *grpc.Server
}

// Creates a new StorageClient with the provided srvAddresses as the configuration
func New(addr string, srvAddresses []string, qSize int) *Client {
	client := &Client{
		id:     0,
		addr:   addr,
		config: config.NewConfig(srvAddresses),
	}
	pb.RegisterPBFTNodeServer(client.srv, client)
	return client
}

func (c *Client) WriteVal(value string) (*pb.ClientResponse, error) {
	c.id++
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pb.WriteRequest{
		Id:        uuid.NewString(),
		From:      c.addr,
		Message:   value,
		Timestamp: time.Now().Unix(),
	}
	return c.config.Write(ctx, req)
}

func (c *Client) Start() {
	lis, err := net.Listen("tcp", c.addr)
	if err != nil {
		panic(err)
	}
	go c.srv.Serve(lis)
}

func (c *Client) Stop() {
	c.srv.Stop()
}

func (c *Client) Write(ctx context.Context, req *pb.WriteRequest) (*empty.Empty, error) {
	return nil, nil
}

func (c *Client) PrePrepare(ctx context.Context, req *pb.PrePrepareRequest) (*empty.Empty, error) {
	return nil, nil
}

func (c *Client) Prepare(ctx context.Context, req *pb.PrepareRequest) (*empty.Empty, error) {
	return nil, nil
}

func (c *Client) Commit(ctx context.Context, req *pb.CommitRequest) (*empty.Empty, error) {
	return nil, nil
}

func (c *Client) ClientHandler(ctx context.Context, req *pb.ClientResponse) (*empty.Empty, error) {
	return nil, nil
}

func (c *Client) Benchmark(ctx context.Context, req *empty.Empty) (*pb.Result, error) {
	return nil, nil
}

/*func (c *Client) Benchmark() (*pb.Result, error) {
	// slog.Info(fmt.Sprintf("client(%v): writing", sc.id), "val", value)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return c.config.Benchmark(ctx, &empty.Empty{})
}*/
