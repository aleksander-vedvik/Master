package client

import (
	pb "reliablebroadcast/proto"
)

type QSpec struct {
	quorumSize int
}

func NewQSpec(qSize int) pb.QuorumSpec {
	return &QSpec{
		quorumSize: qSize,
	}
}

func (qs *QSpec) BroadcastQF(in *pb.Message, replies []*pb.Message) (*pb.Message, bool) {
	if len(replies) < qs.quorumSize {
		return nil, false
	}
	var val *pb.Message
	for _, resp := range replies {
		val = resp
	}
	return val, true
}
