package client

import (
	"log"
	"net"

	pb "github.com/aleksander-vedvik/Master/protos"
	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetBConfig(srvAddresses []string, numSrvs int) *pb.Configuration {
	mgr := pb.NewManager(
		gorums.WithGrpcDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	lis, err := net.Listen("tcp4", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	mgr.AddClientServer(lis)
	quorum, err := mgr.NewConfiguration(
		gorums.WithNodeList(srvAddresses),
		NewQSpec(len(srvAddresses), numSrvs),
	)
	if err != nil {
		log.Fatal("error creating config:", err)
	}
	//quorum.RegisterClientServer("localhost:8080")
	return quorum
}

func GetQConfig(srvAddresses []string) *pb.Configuration {
	mgr := pb.NewManager(
		gorums.WithGrpcDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	quorum, err := mgr.NewConfiguration(
		gorums.WithNodeList(srvAddresses),
		NewQSpec(len(srvAddresses), len(srvAddresses)),
	)
	if err != nil {
		log.Fatal("error creating config:", err)
	}
	//quorum.RegisterClientServer("localhost:8080")
	return quorum
}
