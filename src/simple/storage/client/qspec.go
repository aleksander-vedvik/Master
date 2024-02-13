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

// ReadQF is the quorum function for the Read RPC method.
func (qs *QSpec) BroadcastQF(in *pb.State, replies map[uint32]*pb.ClientResponse) (*pb.ClientResponse, bool) {
	if len(replies) >= qs.quorumSize {
		for _, resp := range replies {
			if resp.GetSuccess() {
				return resp, true
			}
		}
	}
	return nil, false
}

// WriteQF is the quorum function for the Write RPC method.
func (qs *QSpec) DeliverQF(in *pb.State, replies map[uint32]*pb.Empty) (*pb.Empty, bool) {
	return &pb.Empty{}, false
}
