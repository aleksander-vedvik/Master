package server

import (
	pb "github.com/aleksander-vedvik/benchmark/pbft/protos"

	"github.com/relab/gorums"
)

func (s *Server) Ping(ctx gorums.ServerCtx, request *pb.Heartbeat) {
	s.leaderElection.Ping(request.GetId())
}
