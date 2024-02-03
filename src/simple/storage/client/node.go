package client

import (
	"log"
	"time"

	pb "github.com/aleksander-vedvik/Master/protos"
	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetConfig(srvAddresses []string) *pb.Configuration {
	mgr := pb.NewManager(
		gorums.WithDialTimeout(50*time.Millisecond),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	var err error
	quorum, err := mgr.NewConfiguration(
		NewQSpec(len(srvAddresses)),
		gorums.WithNodeList(srvAddresses),
	)
	if err != nil {
		log.Fatal("error creating config:", err)
	}
	return quorum
}
