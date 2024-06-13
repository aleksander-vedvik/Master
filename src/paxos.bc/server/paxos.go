package server

import (
	pb "github.com/aleksander-vedvik/benchmark/paxos.bc/proto"

	"github.com/relab/gorums"
)

func (srv *Server) Prepare(ctx gorums.ServerCtx, req *pb.PrepareMsg) (*pb.PromiseMsg, error) {
	return srv.acceptor.Prepare(ctx, req)
}

func (srv *Server) Accept(ctx gorums.ServerCtx, request *pb.AcceptMsg, broadcast *pb.Broadcast) {
	srv.acceptor.Accept(ctx, request, broadcast)
}

func (srv *Server) Learn(ctx gorums.ServerCtx, request *pb.LearnMsg, broadcast *pb.Broadcast) {
	srv.acceptor.Learn(ctx, request, broadcast)
}
