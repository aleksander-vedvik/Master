package storage

import (
	"log"
	"time"

	pb "github.com/aleksander-vedvik/Master/protos"
	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Node struct {
	node *gorums.GorumsNode
}

func NewNode() *Node {
	node := gorums.NewGorumsNode()

	node.RegisterClient(NewStorageClient([]string{}))
	node.RegisterServer(NewStorageServer())

	node.NewConfiguration([]string{})

	return &Node{node}
}

func getConfig(srvAddresses []string) *pb.Configuration {
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
