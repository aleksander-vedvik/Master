package client

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/aleksander-vedvik/benchmark/pbft.s/config"
	pb "github.com/aleksander-vedvik/benchmark/pbft.s/protos"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Client struct {
	pb.PBFTNodeServer
	mut       sync.Mutex
	id        int
	addr      string
	config    *config.Config
	srv       *grpc.Server
	responses map[string][]*pb.ClientResponse
	respChans map[string]chan *pb.ClientResponse
}

// Creates a new StorageClient with the provided srvAddresses as the configuration
func New(addr string, srvAddresses []string, qSize int) *Client {
	client := &Client{
		id:        0,
		addr:      addr,
		config:    config.NewConfig(srvAddresses),
		responses: make(map[string][]*pb.ClientResponse),
		respChans: make(map[string]chan *pb.ClientResponse),
	}
	pb.RegisterPBFTNodeServer(client.srv, client)
	return client
}

func (c *Client) WriteVal(ctx context.Context, value string) (*pb.ClientResponse, error) {
	id := uuid.NewString()
	respChan := make(chan *pb.ClientResponse, c.config.NumNodes())
	c.mut.Lock()
	c.respChans[id] = respChan
	c.mut.Unlock()
	req := &pb.WriteRequest{
		Id:        id,
		From:      c.addr,
		Message:   value,
		Timestamp: time.Now().Unix(),
	}
	_, _ = c.config.Write(ctx, req)
	select {
	case res := <-respChan:
		return res, nil
	case <-ctx.Done():
		return nil, errors.New("failed req")
	}
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
	c.mut.Lock()
	defer c.mut.Unlock()
	var resps []*pb.ClientResponse
	if r, ok := c.responses[req.Id]; ok {
		resps = append(r, req)
	} else {
		resps = []*pb.ClientResponse{req}
	}
	c.responses[req.Id] = resps
	if len(resps) > 2*c.config.NumNodes()/3 {
		c.respChans[req.Id] <- req
	}
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
