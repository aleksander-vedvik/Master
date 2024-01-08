package test

import (
	"log"
	"time"

	pb "github.com/aleksander-vedvik/Master/protos"
	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GorumsManager struct {
	*pb.Configuration
}
type GorumsQSpec struct {
	QuorumSize int
	Nodes      gorums.NodeListOption
	QSpec      pb.QuorumSpec
}

func NewGorumsManager(conf func(q *GorumsQSpec)) *GorumsManager {
	mgr := pb.NewManager(
		gorums.WithDialTimeout(50*time.Millisecond),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	qSpec := newQSpecWithDefaults()
	conf(qSpec)

	allNodesConfig, err := mgr.NewConfiguration(
		qSpec.QSpec,
		qSpec.Nodes,
	)
	if err != nil {
		log.Fatal("error creating config:", err)
		return nil
	}
	return &GorumsManager{
		allNodesConfig,
	}
}

func newQSpecWithDefaults() *GorumsQSpec {
	return &GorumsQSpec{}
}

func (g *GorumsManager) Write(val string) (string, error) {
	return "", nil
}
