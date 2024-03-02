package client

import (
	pb "github.com/aleksander-vedvik/Master/protos"
)

type QSpec struct {
	quorumSize    int
	broadcastSize int
}

func NewQSpec(qSize, broadcastSize int) pb.QuorumSpec {
	return &QSpec{
		quorumSize:    qSize,
		broadcastSize: broadcastSize,
	}
}

func (qs *QSpec) BroadcastQF(in *pb.State, replies map[uint32]*pb.View) (*pb.View, bool) {
	if len(replies) >= qs.quorumSize {
		for _, resp := range replies {
			return resp, true
		}
	}
	return nil, false
}

func (qs *QSpec) SaveStudentQF(reqs []*pb.ClientResponse) (*pb.ClientResponse, bool) {
	if len(reqs) < qs.broadcastSize {
		return nil, false
	}
	return reqs[0], true
}

func (qs *QSpec) SaveStudentsQF(reqs []*pb.ClientResponse) (*pb.ClientResponse, bool) {
	if len(reqs) < qs.broadcastSize {
		return nil, false
	}
	return reqs[0], true
}
