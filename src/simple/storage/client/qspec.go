package client

import (
	pb "github.com/aleksander-vedvik/Master/protos"
)

type QSpec struct {
	quorumSize int
}

func NewQSpec(qSize int) pb.QuorumSpec {
	return &QSpec{
		quorumSize: qSize,
	}
}

func (qs *QSpec) BroadcastQF(in *pb.State, replies map[uint32]*pb.View) (*pb.View, bool) {
	if len(replies) >= qs.quorumSize {
		for _, resp := range replies {
			if resp.GetNumberOfServers() > 0 {
				return resp, true
			}
		}
	}
	return nil, false
}

func (qs *QSpec) SaveStudentQF(reqs []*pb.ClientResponse) (*pb.ClientResponse, bool) {
	if len(reqs) < qs.quorumSize {
		return nil, false
	}
	return reqs[0], true
}

func (qs *QSpec) SaveStudentsQF(reqs []*pb.ClientResponse) (*pb.ClientResponse, bool) {
	if len(reqs) < qs.quorumSize {
		return nil, false
	}
	return reqs[0], true
}
