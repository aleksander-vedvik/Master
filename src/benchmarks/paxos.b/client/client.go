package client

import (
	"context"
	"log"
	"net"
	"time"

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

func New(addr string, srvAddresses []string, qSize int) *Client {
	mgr := pb.NewManager(
		gorums.WithGrpcDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	mgr.AddClientServer(lis)
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

func (sc *Client) Benchmark() (*pb.Result, error) {
	//slog.Info(fmt.Sprintf("client(%v): writing", sc.id), "val", value)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return sc.config.Benchmark(ctx, &pb.Empty{})
}
