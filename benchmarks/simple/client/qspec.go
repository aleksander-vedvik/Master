package client

import (
	"log"
	"log/slog"
	"net"

	pb "github.com/aleksander-vedvik/benchmark/simple/protos"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getConfig(id int, addr string, srvAddresses []string, qSize int, logger *slog.Logger) (*pb.Manager, *pb.Configuration) {
	mgr := pb.NewManager(
		gorums.WithGrpcDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
		gorums.WithMachineID(uint64(id)),
		gorums.WithLogger(logger),
	)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	mgr.AddClientServer(lis, gorums.WithSrvID(uint64(id)), gorums.WithSLogger(logger))
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

func (qs *QSpec) BroadcastCall1QF(in *pb.WriteRequest1, replies []*pb.WriteResponse1) (*pb.WriteResponse1, bool) {
	if len(replies) < qs.quorumSize {
		return nil, false
	}
	var val *pb.WriteResponse1
	for _, resp := range replies {
		val = resp
		break
	}
	return val, true
}

func (qs *QSpec) BroadcastCall2QF(in *pb.WriteRequest2, replies []*pb.WriteResponse2) (*pb.WriteResponse2, bool) {
	if len(replies) < qs.quorumSize {
		return nil, false
	}
	var val *pb.WriteResponse2
	for _, resp := range replies {
		val = resp
		break
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
