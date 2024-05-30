package client

import (
	"context"
	"log"
	"log/slog"

	pb "github.com/aleksander-vedvik/benchmark/paxosqc/proto"
	"github.com/google/uuid"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	config *pb.Configuration
	mgr    *pb.Manager
	seq    uint32
}

func New(id int, addr string, srvAddresses []string, qSize int, logger *slog.Logger) *Client {
	mgr := pb.NewManager(
		gorums.WithGrpcDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
		gorums.WithMachineID(uint64(id)),
		gorums.WithLogger(logger),
	)
	config, err := mgr.NewConfiguration(
		gorums.WithNodeList(srvAddresses),
		NewPaxosQSpec(qSize),
	)
	if err != nil {
		log.Fatal("error creating config:", err)
	}
	return &Client{
		config: config,
		mgr:    mgr,
	}
}

func (sc *Client) Stop() {
	sc.mgr.Close()
}

func (sc *Client) Write(ctx context.Context, value string) (*pb.Response, error) {
	sc.seq++
	return sc.config.ClientHandle(ctx, &pb.Value{
		ClientCommand: value,
		ClientSeq:     sc.seq,
		ID:            uuid.NewString(),
	})
}

func (sc *Client) Benchmark(ctx context.Context) (*pb.Empty, error) {
	// slog.Info(fmt.Sprintf("client(%v): writing", sc.id), "val", value)
	return sc.config.Benchmark(ctx, &pb.Empty{})
}
