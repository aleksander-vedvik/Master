package client

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	pb "paxos/proto"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StorageClient struct {
	config *pb.Configuration
}

func NewStorageClient(id int, srvAddresses []string) *StorageClient {
	mgr := pb.NewManager(
		gorums.WithGrpcDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	config, err := mgr.NewConfiguration(
		gorums.WithNodeList(srvAddresses),
		newQSpec(1+len(srvAddresses)/2),
		gorums.WithListener(fmt.Sprintf("127.0.0.1:%v", 8080+id)),
	)
	if err != nil {
		log.Fatal("error creating config:", err)
	}
	return &StorageClient{
		config: config,
	}
}

func (sc *StorageClient) Write(value string) {
	slog.Info("client: writing value", "val", value)
	ctx := context.Background()
	_, _ = sc.config.Write(ctx, &pb.PaxosValue{
		Val: value,
	})
}
