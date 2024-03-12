package nodeServer

import (
	pb "pbft/protos"

	"github.com/relab/gorums"
)

func (s *PBFTServer) PrePrepare(ctx gorums.ServerCtx, request *pb.PrePrepareRequest, broadcast *pb.Broadcast) {
	if !s.isInView(request.View) {
		return
	}
	if !s.sequenceNumberIsValid(request.SequenceNumber) {
		return
	}
	if s.hasAlreadyAcceptedSequenceNumber(request.SequenceNumber) {
		return
	}
	s.messageLog.add(request, s.viewNumber, request.SequenceNumber)
	go broadcast.Prepare(&pb.PrepareRequest{
		Id:             request.Id,
		View:           request.View,
		SequenceNumber: request.SequenceNumber,
		Digest:         request.Digest,
		Message:        request.Message,
		Timestamp:      request.Timestamp,
	})
}

func (s *PBFTServer) Prepare(ctx gorums.ServerCtx, request *pb.PrepareRequest, broadcast *pb.Broadcast) {
	if !s.isInView(request.View) {
		return
	}
	if !s.sequenceNumberIsValid(request.SequenceNumber) {
		return
	}
	s.messageLog.add(request, s.viewNumber, request.SequenceNumber)
	if s.prepared(request.SequenceNumber) {
		go broadcast.Commit(&pb.CommitRequest{
			Id:             request.Id,
			Timestamp:      request.Timestamp,
			View:           request.View,
			SequenceNumber: request.SequenceNumber,
			Message:        request.Message,
		})
	}
}

func (s *PBFTServer) Commit(ctx gorums.ServerCtx, request *pb.CommitRequest, broadcast *pb.Broadcast) {
	if !s.isInView(request.View) {
		return
	}
	if !s.sequenceNumberIsValid(request.SequenceNumber) {
		return
	}
	s.messageLog.add(request, s.viewNumber, request.SequenceNumber)
	if s.committed(request.SequenceNumber) {
		s.state = &pb.ClientResponse{
			Id:        request.Id,
			Result:    request.Message,
			Timestamp: request.Timestamp,
			View:      request.View,
		}
		go broadcast.SendToClient(s.state, nil)
	}
}

func (s *PBFTServer) requestIsAlreadyProcessed(req *pb.WriteRequest) (*pb.ClientResponse, bool) {
	return s.state, s.state != nil && s.state.Id == req.Id
}

func (s *PBFTServer) isLeader() bool {
	return s.leader == s.addr
}

func (s *PBFTServer) isInView(view int32) bool {
	// request is in the current view of this node
	return s.viewNumber == view
}

func (s *PBFTServer) sequenceNumberIsValid(n int32) bool {
	// the sequence number is between h and H
	return true
}

func (s *PBFTServer) hasAlreadyAcceptedSequenceNumber(n int32) bool {
	// checks if the node has already accepted a pre-prepare request
	// for this view and sequence number
	_, found := s.messageLog.getPrePrepareReq(n, s.viewNumber)
	return found
}

func (s *PBFTServer) prepared(n int32) bool {
	// the request m, a pre-prepare for m in view v with sequence number n,
	// and 2f prepares from different backups that match the pre-prepare.
	// The replicas verify whether the prepares match the pre-prepare by
	// checking that they have the same view, sequence number, and digest.
	req, found := s.messageLog.getPrePrepareReq(n, s.viewNumber)
	if !found {
		return false
	}
	reqs, found := s.messageLog.getPrepareReqs(req.Digest, req.SequenceNumber, req.View)
	return found && len(reqs) > 2*len(s.peers)/3
}

func (s *PBFTServer) committed(n int32) bool {
	/*
		We define the committed and committed-local predi-
		cates as follows: committed(m, v, n) is true if and only
		if prepared(m, v, n, i) is true for all i in some set of
		f+1 non-faulty replicas; and committed-local(m, v, n, i)
		is true if and only if prepared(m, v, n, i) is true and i has
		accepted 2f + 1 commits (possibly including its own)
		from different replicas that match the pre-prepare for m;
		a commit matches a pre-prepare if they have the same
		view, sequence number, and digest.
	*/
	req, found := s.messageLog.getPrePrepareReq(n, s.viewNumber)
	if !found {
		return false
	}
	reqs, found := s.messageLog.getPrepareReqs(req.Digest, req.SequenceNumber, req.View)
	if !found || len(reqs) <= 2*len(s.peers)/3 {
		return false
	}
	commits, found := s.messageLog.getCommitReqs(req.Digest, req.SequenceNumber, req.View)
	return found && len(commits) > 2*len(s.peers)/3
}
