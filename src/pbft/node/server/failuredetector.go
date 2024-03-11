package nodeServer

import (
	pb "pbft/protos"

	"github.com/relab/gorums"
)

func (s *PBFTServer) Heartbeat(ctx gorums.ServerCtx, request *pb.HeartbeatReq) {
	s.leaderElection.Heartbeat(request)
}
