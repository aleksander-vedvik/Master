package client

import (
	"log/slog"
	pb "paxos/proto"
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

func (q *QSpec) WriteQF(replies []*pb.Response) (*pb.Response, bool) {
	slog.Info("got response")
	if len(replies) < q.qsize {
		return nil, false
	}
	return nil, true
}
