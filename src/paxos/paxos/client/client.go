package client

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "paxos/proto"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StorageClient struct {
	id     int
	config *pb.Configuration
	mgr    *pb.Manager
}

func NewStorageClient(id int, srvAddresses []string) *StorageClient {
	mgr := pb.NewManager(
		gorums.WithGrpcDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	addr := fmt.Sprintf("127.0.0.1:%v", 8080+id)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	mgr.AddClientServer(lis)
	config, err := mgr.NewConfiguration(
		gorums.WithNodeList(srvAddresses),
		newQSpec(1+len(srvAddresses)/2),
	)
	if err != nil {
		log.Fatal("error creating config:", err)
	}
	return &StorageClient{
		id:     id,
		config: config,
		mgr:    mgr,
	}
}

func (sc *StorageClient) Stop() {
	sc.mgr.Close()
}

func (sc *StorageClient) Write(value string) {
	//slog.Info(fmt.Sprintf("client(%v): writing", sc.id), "val", value)
	ctx := context.Background()
	_, _ = sc.config.Write(ctx, &pb.PaxosValue{
		Val: value,
	})
}
