package storage

import pb "ui/frontend/src/protos"

type QSpec struct {
	quorumSize int
}

func NewQSpec(qSize int) pb.QuorumSpec {
	return &QSpec{
		quorumSize: qSize,
	}
}

// ReadQF is the quorum function for the Read RPC method.
func (qs *QSpec) ReadQF(in *pb.ReadRequest, replies map[uint32]*pb.State) (*pb.State, bool) {
	if len(replies) < qs.quorumSize {
		return nil, false
	}
	return newestState(replies), true
}

// WriteQF is the quorum function for the Write RPC method.
func (qs *QSpec) WriteQF(in *pb.State, replies map[uint32]*pb.WriteResponse) (*pb.WriteResponse, bool) {
	if len(replies) < qs.quorumSize {
		return nil, false
	}
	// return the first response we find
	var reply *pb.WriteResponse
	for _, r := range replies {
		reply = r
		break
	}
	return reply, true
}

func newestState(replies map[uint32]*pb.State) *pb.State {
	var newest *pb.State
	for _, s := range replies {
		if s.GetTimestamp() >= newest.GetTimestamp() {
			newest = s
		}
	}
	return newest
}
