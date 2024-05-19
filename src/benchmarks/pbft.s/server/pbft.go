package server

import (
	"context"

	pb "github.com/aleksander-vedvik/benchmark/pbft.s/protos"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) PrePrepare(ctx context.Context, request *pb.PrePrepareRequest) (*empty.Empty, error) {
	//slog.Info("1 server: received preprepare", "who", s.addr)
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
	//slog.Info("server: added preprepare", "who", s.addr, "v", s.viewNumber, "s", request.SequenceNumber)
	s.view.Prepare(&pb.PrepareRequest{
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
	//slog.Info("2 server: received prepare", "who", s.addr)
	if !s.isInView(request.View) {
		//slog.Error("server: not in view")
		return nil, nil
	}
	if !s.sequenceNumberIsValid(request.SequenceNumber) {
		//slog.Error("server: not valid sequence number")
		return nil, nil
	}
	s.messageLog.add(request, s.viewNumber, request.SequenceNumber)
	if s.prepared(request.SequenceNumber) {
		//slog.Error("prepared", "who", s.addr)
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
	//slog.Info("3 server: received commit")
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

func (s *Server) requestIsAlreadyProcessed(req *pb.WriteRequest) (*pb.ClientResponse, bool) {
	return s.state, s.state != nil && s.state.Id == req.Id
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
		//slog.Error("server: not found", "who", s.addr, "v", s.viewNumber, "s", n)
		return false
	}
	//slog.Error("server: before prepared", "who", s.addr)
	reqs, found := s.messageLog.getPrepareReqs(req.Digest, req.SequenceNumber, req.View)
	//slog.Error("server: prepared", "num", len(reqs), "who", s.addr)
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
	return found && len(commits) >= 2*len(s.peers)/3
}
