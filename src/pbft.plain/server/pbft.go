package server

import (
	"context"

	pb "github.com/aleksander-vedvik/benchmark/pbft.plain/protos"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) PrePrepare(ctx context.Context, request *pb.PrePrepareRequest) (*empty.Empty, error) {
	if !s.isInView(request.View) {
		return nil, nil
	}
	if !s.sequenceNumberIsValid(request.SequenceNumber) {
		return nil, nil
	}
	if s.hasAlreadyAcceptedSequenceNumber(request.SequenceNumber) {
		return nil, nil
	}
	s.messageLog.add(request, s.viewNumber, request.SequenceNumber)
	s.progressState(request.Id)
	go s.handlePrepares(request.Id)
	go s.view.Prepare(&pb.PrepareRequest{
		Id:             request.Id,
		View:           request.View,
		SequenceNumber: request.SequenceNumber,
		Digest:         request.Digest,
		Message:        request.Message,
		Timestamp:      request.Timestamp,
		From:           request.From,
	})
	return nil, nil
}

func (s *Server) Prepare(ctx context.Context, request *pb.PrepareRequest) (*empty.Empty, error) {
	s.mut.Lock()
	if s.order[request.Id] < Prepare {
		if s.tmpPrepares[request.Id] == nil || len(s.tmpPrepares[request.Id]) <= 0 {
			s.tmpPrepares[request.Id] = []*pb.PrepareRequest{request}
		} else {
			s.tmpPrepares[request.Id] = append(s.tmpPrepares[request.Id], request)
		}
		s.mut.Unlock()
		return nil, nil
	}
	s.mut.Unlock()
	if !s.isInView(request.View) {
		return nil, nil
	}
	if !s.sequenceNumberIsValid(request.SequenceNumber) {
		return nil, nil
	}
	s.messageLog.add(request, s.viewNumber, request.SequenceNumber)
	if s.prepared(request.SequenceNumber) {
		s.progressState(request.Id)
		go s.handleCommits(request.Id)
		go s.view.Commit(&pb.CommitRequest{
			Id:             request.Id,
			Timestamp:      request.Timestamp,
			View:           request.View,
			Digest:         request.Digest,
			SequenceNumber: request.SequenceNumber,
			Message:        request.Message,
			From:           request.From,
		})
	}
	return nil, nil
}

func (s *Server) Commit(ctx context.Context, request *pb.CommitRequest) (*empty.Empty, error) {
	s.mut.Lock()
	if s.order[request.Id] < Commit {
		if s.tmpCommits[request.Id] == nil || len(s.tmpCommits[request.Id]) <= 0 {
			s.tmpCommits[request.Id] = []*pb.CommitRequest{request}
		} else {
			s.tmpCommits[request.Id] = append(s.tmpCommits[request.Id], request)
		}
		s.mut.Unlock()
		return nil, nil
	}
	s.mut.Unlock()
	if !s.isInView(request.View) {
		return nil, nil
	}
	if !s.sequenceNumberIsValid(request.SequenceNumber) {
		return nil, nil
	}
	s.messageLog.add(request, s.viewNumber, request.SequenceNumber)
	if s.committed(request.SequenceNumber) {
		s.mut.Lock()
		state := &pb.ClientResponse{
			Id:        request.Id,
			Result:    request.Message,
			Timestamp: request.Timestamp,
			View:      request.View,
			From:      request.From,
		}
		s.state = state
		s.mut.Unlock()
		go s.view.ClientHandler(state)
	}
	return nil, nil
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

func (s *Server) progressState(requestId string) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.order[requestId]++
}

func (s *Server) handlePrepares(requestId string) {
	s.mut.Lock()
	defer s.mut.Unlock()
	for _, prepare := range s.tmpPrepares[requestId] {
		go s.Prepare(context.Background(), prepare)
	}
}

func (s *Server) handleCommits(requestId string) {
	s.mut.Lock()
	defer s.mut.Unlock()
	for _, commit := range s.tmpCommits[requestId] {
		go s.Commit(context.Background(), commit)
	}
}
