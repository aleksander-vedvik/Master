package node

import (
	"log"
	pb "pbft/protos"
	"time"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getConfig(addr string, srvAddresses []string) *pb.Configuration {
	mgr := pb.NewManager(
		gorums.WithDialTimeout(50*time.Millisecond),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	var err error
	quorum, err := mgr.NewConfiguration(
		NewQSpec(1),
		gorums.WithNodeList(srvAddresses),
	)
	if err != nil {
		log.Fatal("error creating config:", err)
	}
	//quorum.AddSender(addr)
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

func (qs *QSpec) WriteQF(in *pb.WriteRequest, replies map[uint32]*pb.ClientResponse) (*pb.ClientResponse, bool) {
	/*if len(replies) < qs.quorumSize {
		return nil, false
	}*/
	var val *pb.ClientResponse
	for _, resp := range replies {
		val = resp
	}
	return val, true
}
