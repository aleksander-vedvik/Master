package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/aleksander-vedvik/benchmark/pbft.plain/config"
	pb "github.com/aleksander-vedvik/benchmark/pbft.plain/protos"
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
	select {
	case res := <-respChan:
		return res, nil
	case <-reqctx.Done():
		return nil, errors.New("failed req")
	}
}

func (c *Client) Start() {
	splittedAddr := strings.Split(c.addr, ":")
	lis, err := net.Listen("tcp", ":"+splittedAddr[1])
	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf("client started. Listening on address: %s, lis=%s\n", c.addr, lis.Addr().String()))
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
	var resps []*pb.ClientResponse
	if r, ok := c.responses[req.Id]; ok {
		resps = append(r, req)
	} else {
		resps = []*pb.ClientResponse{req}
	}
	c.responses[req.Id] = resps
	respChan := c.resps[req.Id].respChan
	respCtx := c.resps[req.Id].ctx
	c.mut.Unlock()
	if len(resps) >= 2*c.config.NumNodes()/3 {
		go func() {
			select {
			case respChan <- req:
			case <-respCtx.Done():
			}
		}()
	}
	return nil, nil
}

func (c *Client) SendBenchmark(ctx context.Context) {
	c.config.Benchmark(ctx)
}

func (srv *Client) Benchmark(ctx context.Context, request *empty.Empty) (*pb.Result, error) {
	return &pb.Result{}, nil
}
