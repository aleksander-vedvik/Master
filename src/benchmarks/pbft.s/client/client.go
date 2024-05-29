package client

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/aleksander-vedvik/benchmark/pbft.s/config"
	pb "github.com/aleksander-vedvik/benchmark/pbft.s/protos"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type resp struct {
	ctx      context.Context
	respChan chan *pb.ClientResponse
}

type Client struct {
	pb.PBFTNodeServer
	mut       sync.Mutex
	id        int
	addr      string
	config    *config.Config
	srv       *grpc.Server
	responses map[string][]*pb.ClientResponse
	resps     map[string]*resp
}

// Creates a new StorageClient with the provided srvAddresses as the configuration
func New(addr string, srvAddresses []string, qSize int) *Client {
	client := &Client{
		id:        0,
		addr:      addr,
		config:    config.NewConfig(addr, srvAddresses),
		responses: make(map[string][]*pb.ClientResponse),
		resps:     make(map[string]*resp),
		srv:       grpc.NewServer(),
	}
	pb.RegisterPBFTNodeServer(client.srv, client)
	return client
}

func (c *Client) WriteVal(ctx context.Context, value string) (*pb.ClientResponse, error) {
	//slog.Info("writing value", "val", value)
	id := uuid.NewString()
	respChan := make(chan *pb.ClientResponse, c.config.NumNodes())
	reqctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	c.mut.Lock()
	c.resps[id] = &resp{
		ctx:      reqctx,
		respChan: respChan,
	}
	c.mut.Unlock()
	req := &pb.WriteRequest{
		Id:        id,
		From:      c.addr,
		Message:   value,
		Timestamp: time.Now().Unix(),
	}
	_, _ = c.config.Write(ctx, req)
	//return <-respChan, nil
	select {
	case res := <-respChan:
		return res, nil
	case <-reqctx.Done():
		return nil, errors.New("failed req")
		//case <-ctx.Done():
		//return nil, errors.New("failed req")
	}
}

func (c *Client) Start() {
	lis, err := net.Listen("tcp", c.addr)
	if err != nil {
		panic(err)
	}
	go c.srv.Serve(lis)
	slog.Info("client started", "addr", c.addr)
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
	//slog.Info("received reply")
	c.mut.Lock()
	defer c.mut.Unlock()
	var resps []*pb.ClientResponse
	if r, ok := c.responses[req.Id]; ok {
		resps = append(r, req)
	} else {
		resps = []*pb.ClientResponse{req}
	}
	c.responses[req.Id] = resps
	if len(resps) >= 2*c.config.NumNodes()/3 {
		select {
		case c.resps[req.Id].respChan <- req:
		case <-c.resps[req.Id].ctx.Done():
		}
	}
	return nil, nil
}

func (c *Client) SendBenchmark(ctx context.Context) {
	c.config.Benchmark(ctx)
}

func (srv *Client) Benchmark(ctx context.Context, request *empty.Empty) (*pb.Result, error) {
	return &pb.Result{}, nil
}
