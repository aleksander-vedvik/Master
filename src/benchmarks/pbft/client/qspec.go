package client

import (
	"log"
	"net"

	pb "github.com/aleksander-vedvik/benchmark/pbft/protos"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getConfig(id int, addr string, srvAddresses []string, qSize int) (*pb.Manager, *pb.Configuration) {
	mgr := pb.NewManager(
		gorums.WithGrpcDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
		gorums.WithMachineID(uint64(id)),
	)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	mgr.AddClientServer(lis)
	quorum, err := mgr.NewConfiguration(
		NewQSpec(qSize),
		gorums.WithNodeList(srvAddresses),
	)
	if err != nil {
		log.Fatal("error creating config:", err)
	}
	//quorum.AddSender(addr)
	return mgr, quorum
}

type QSpec struct {
	quorumSize int
}

func NewQSpec(qSize int) pb.QuorumSpec {
	return &QSpec{
		quorumSize: qSize,
	}
}

func (qs *QSpec) WriteQF(in *pb.WriteRequest, replies []*pb.ClientResponse) (*pb.ClientResponse, bool) {
	if len(replies) < qs.quorumSize {
		return nil, false
	}
	var val *pb.ClientResponse
	for _, resp := range replies {
		val = resp
	}
	return val, true
}

func (qs *QSpec) BenchmarkQF(in *empty.Empty, replies map[uint32]*pb.Result) (*pb.Result, bool) {
	if len(replies) < qs.quorumSize {
		return nil, false
	}
	result := &pb.Result{
		Metrics: make([]*pb.Metric, 0, len(replies)),
	}
	for _, reply := range replies {
		result.Metrics = append(result.Metrics, reply.Metrics...)
	}
	return result, true
}
