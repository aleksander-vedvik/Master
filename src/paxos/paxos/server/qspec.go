package server

import pb "paxos/proto"

type QSpec struct {
	qsize int
}

func newQSpec(qsize int) *QSpec {
	return &QSpec{qsize: qsize}
}

func (q *QSpec) PrepareQF(in *pb.PrepareMsg, replies map[uint32]*pb.PromiseMsg) {
	if len(replies) < q.qsize {
		return
	}
}
