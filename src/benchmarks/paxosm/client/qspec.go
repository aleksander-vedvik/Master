package client

import (
	pb "github.com/aleksander-vedvik/benchmark/paxosm/proto"
)

type QSpec struct {
	qsize int
}

func newQSpec(qsize int) pb.QuorumSpec {
	return &QSpec{qsize: qsize}
}

// not used by the client
func (q *QSpec) PrepareQF(in *pb.PrepareMsg, replies map[uint32]*pb.PromiseMsg) (*pb.PromiseMsg, bool) {
	return nil, true
}

func (q *QSpec) WriteQF(in *pb.PaxosValue, replies map[uint32]*pb.PaxosResponse) (*pb.PaxosResponse, bool) {
	if len(replies) < q.qsize {
		return nil, false
	}
	return nil, true
}

func (q *QSpec) BenchmarkQF(in *pb.Empty, replies map[uint32]*pb.Result) (*pb.Result, bool) {
	if len(replies) < q.qsize {
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
