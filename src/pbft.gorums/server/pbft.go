package server

import (
	"fmt"

	pb "github.com/aleksander-vedvik/benchmark/pbft.gorums/protos"

	"github.com/relab/gorums"
)

func (s *Server) PrePrepare(ctx gorums.ServerCtx, request *pb.PrePrepareRequest, broadcast *pb.Broadcast) {
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
	broadcast.Prepare(&pb.PrepareRequest{
		Id:             request.Id,
		View:           request.View,
		SequenceNumber: request.SequenceNumber,
		Digest:         request.Digest,
		Message:        request.Message,
		Timestamp:      request.Timestamp,
	}, gorums.WithoutSelf())
}

func (s *Server) Prepare(ctx gorums.ServerCtx, request *pb.PrepareRequest, broadcast *pb.Broadcast) {
	if !s.isInView(request.View) {
		return
	}
	if !s.sequenceNumberIsValid(request.SequenceNumber) {
		return
	}
	s.messageLog.add(request, s.viewNumber, request.SequenceNumber)
	if s.prepared(request.SequenceNumber) {
		broadcast.Commit(&pb.CommitRequest{
			Id:             request.Id,
			Timestamp:      request.Timestamp,
			View:           request.View,
			Digest:         request.Digest,
			SequenceNumber: request.SequenceNumber,
			Message:        request.Message,
		}, gorums.WithoutSelf())
	}
}

func (s *Server) Commit(ctx gorums.ServerCtx, request *pb.CommitRequest, broadcast *pb.Broadcast) {
	if !s.isInView(request.View) {
		return
	}
	if !s.sequenceNumberIsValid(request.SequenceNumber) {
		return
	}
	s.messageLog.add(request, s.viewNumber, request.SequenceNumber)
	if s.committed(request.SequenceNumber) {
		state := &pb.ClientResponse{
			Id:        request.Id,
			Result:    request.Message,
			Timestamp: request.Timestamp,
			View:      request.View,
		}
		s.mut.Lock()
		s.state = state
		s.mut.Unlock()
		broadcast.SendToClient(state, nil)
	}
}

func (s *Server) requestIsAlreadyProcessed(broadcast *pb.Broadcast) (*pb.ClientResponse, bool) {
	s.mut.Lock()
	defer s.mut.Unlock()
	md := broadcast.GetMetadata()
	clientID := md.MachineID
	seqNo := md.SequenceNo
	id := fmt.Sprintf("c%v,s%v", clientID, seqNo)
	return s.state, s.state != nil && s.state.Id == id
}

func (s *Server) isLeader() bool {
	return s.leader == s.addr
}

func (s *Server) isInView(view int32) bool {
	// request is in the current view of this node
	return s.viewNumber == view
}

func (s *Server) sequenceNumberIsValid(n int32) bool {
	// the sequence number is between h and H
	return n >= 0
}

func (s *Server) hasAlreadyAcceptedSequenceNumber(n int32) bool {
	// checks if the node has already accepted a pre-prepare request
	// for this view and sequence number
	_, found := s.messageLog.getPrePrepareReq(n, s.viewNumber)
	return found
}

func (s *Server) prepared(n int32) bool {
	// the request m, a pre-prepare for m in view v with sequence number n,
	// and 2f prepares from different backups that match the pre-prepare.
	// The replicas verify whether the prepares match the pre-prepare by
	// checking that they have the same view, sequence number, and digest.
	req, found := s.messageLog.getPrePrepareReq(n, s.viewNumber)
	if !found {
		return false
	}
	reqs, found := s.messageLog.getPrepareReqs(req.Digest, req.SequenceNumber, req.View)
	return found && len(reqs) >= 2*len(s.peers)/3
}

func (s *Server) committed(n int32) bool {
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
	if !found || len(reqs) < 2*len(s.peers)/3 {
		return false
	}
	commits, found := s.messageLog.getCommitReqs(req.Digest, req.SequenceNumber, req.View)
	return found && len(commits) > 2*len(s.peers)/3
}
