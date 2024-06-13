package client

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strings"

	pb "github.com/aleksander-vedvik/benchmark/paxos.bc/proto"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	config *pb.Configuration
	mgr    *pb.Manager
	addr   string
}

func New(id int, addr string, srvAddresses []string, qSize int, logger *slog.Logger) *Client {
	mgr := pb.NewManager(
		gorums.WithGrpcDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
		gorums.WithMachineID(uint64(id)),
		gorums.WithLogger(logger),
	)
	splittedAddr := strings.Split(addr, ":")
	lis, err := net.Listen("tcp", ":"+splittedAddr[1])
	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf("ClientServer started. Listening on address: %s, lis=%s\n", addr, lis.Addr().String()))
	address, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}
	mgr.AddClientServer(lis, address, gorums.WithSrvID(uint64(id)), gorums.WithSLogger(logger))
	config, err := mgr.NewConfiguration(
		gorums.WithNodeList(srvAddresses),
		newQSpec(qSize),
	)
	if err != nil {
		panic("error creating config")
	}
	return &Client{
		config: config,
		mgr:    mgr,
		addr:   addr,
	}
}

func (sc *Client) Stop() {
	sc.mgr.Close()
}

func (sc *Client) Write(ctx context.Context, value string) (*pb.PaxosResponse, error) {
	return sc.config.Write(ctx, &pb.PaxosValue{
		Val: value,
	})
}

func (sc *Client) Benchmark(ctx context.Context) (*pb.Result, error) {
	return sc.config.Benchmark(ctx, &pb.Empty{})
}
