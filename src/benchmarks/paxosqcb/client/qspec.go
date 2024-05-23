package client

import (
	pb "github.com/aleksander-vedvik/benchmark/paxosqcb/proto"
)

// PaxosQSpec is a quorum specification object for Paxos.
// It only holds the quorum size.
type PaxosQSpec struct {
	quorum int
}

// NewPaxosQSpec returns a quorum specification object for Paxos
// for the given configuration size n.
func NewPaxosQSpec(n int) PaxosQSpec {
	return PaxosQSpec{quorum: (n-1)/2 + 1}
}

func (qs PaxosQSpec) PrepareQF(prepare *pb.PrepareMsg, replies map[uint32]*pb.PromiseMsg) (*pb.PromiseMsg, bool) {
	return nil, true
}

func (qs PaxosQSpec) AcceptQF(accept *pb.AcceptMsg, replies map[uint32]*pb.LearnMsg) (*pb.LearnMsg, bool) {
	return nil, true
}

// ClientHandleQF is the quorum function to process the replies from the ClientHandle quorum call.
// This is where the Client handle the replies from the replicas. The quorum function should
// validate the replies against the request, and only valid replies should be considered.
// The quorum function returns true if a quorum of the replicas replied with the same response,
// and a single response is returned. Nil and false is returned if no quorum of valid replies was found.
func (qs PaxosQSpec) ClientHandleQF(request *pb.Value, replies map[uint32]*pb.Response) (*pb.Response, bool) {
	validResponses := len(replies)
	if validResponses < qs.quorum {
		return nil, false
	}
	// remove the invalid responses
	var response *pb.Response
	for _, resp := range replies {
		if !request.Match(resp) {
			validResponses--
		} else {
			response = resp
		}
	}
	if validResponses < qs.quorum {
		return nil, false
	}
	return response, true
}
