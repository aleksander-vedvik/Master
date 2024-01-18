package node

import (
	"log"
	pb "pbft/protos"
	"time"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

type QSpec struct {
	quorumSize int
}

func NewQSpec(qSize int) pb.QuorumSpec {
	return &QSpec{
		quorumSize: qSize,
	}
}

func (qs *QSpec) PrePrepareQF(in *pb.PrePrepareRequest, replies map[uint32]*pb.Empty) (*pb.Empty, bool) {
	if len(replies) < qs.quorumSize {
		return nil, false
	}
	return replies[0], true
}

func (qs *QSpec) PrepareQF(in *pb.PrepareRequest, replies map[uint32]*pb.Empty) (*pb.Empty, bool) {
	if len(replies) < qs.quorumSize {
		return nil, false
	}
	return replies[0], true
}

func (qs *QSpec) CommitQF(in *pb.CommitRequest, replies map[uint32]*pb.Empty) (*pb.Empty, bool) {
	if len(replies) < qs.quorumSize {
		return nil, false
	}
	return replies[0], true
}
