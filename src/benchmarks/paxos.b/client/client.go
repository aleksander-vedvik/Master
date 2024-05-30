package client

import (
	"context"
	"log"
	"log/slog"
	"net"
	"strings"

	pb "github.com/aleksander-vedvik/benchmark/paxos.b/proto"

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
	)
	//lis, err := net.Listen("tcp", addr)
	splittedAddr := strings.Split(addr, ":")
	lis, err := net.Listen("tcp", ":"+splittedAddr[1])
	if err != nil {
		panic(err)
	}
	mgr.AddClientServer(lis, gorums.WithSrvID(uint64(id)), gorums.WithSLogger(logger), gorums.WithListenAddr(addr))
	config, err := mgr.NewConfiguration(
		gorums.WithNodeList(srvAddresses),
		newQSpec(qSize),
	)
	if err != nil {
		log.Fatal("error creating config:", err)
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
	//slog.Info(fmt.Sprintf("client(%v): writing", sc.addr), "val", value)
	return sc.config.Write(ctx, &pb.PaxosValue{
		Val: value,
	})
}

func (sc *Client) Benchmark(ctx context.Context) (*pb.Result, error) {
	//slog.Info(fmt.Sprintf("client(%v): writing", sc.id), "val", value)
	return sc.config.Benchmark(ctx, &pb.Empty{})
}
