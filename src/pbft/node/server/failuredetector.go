package nodeServer

import (
	pb "pbft/protos"

	"github.com/relab/gorums"
)

func (s *PBFTServer) Ping(ctx gorums.ServerCtx, request *pb.Heartbeat) {
	s.leaderElection.Ping(request.GetId())
}
