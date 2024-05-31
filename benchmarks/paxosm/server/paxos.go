package server

import (
	pb "github.com/aleksander-vedvik/benchmark/paxosm/proto"

	"github.com/relab/gorums"
)

func (srv *Server) Prepare(ctx gorums.ServerCtx, req *pb.PrepareMsg) (*pb.PromiseMsg, error) {
	return srv.acceptor.Prepare(ctx, req)
}

func (srv *Server) Accept(ctx gorums.ServerCtx, request *pb.AcceptMsg) {
	srv.acceptor.Accept(ctx, request)
}

func (srv *Server) Learn(ctx gorums.ServerCtx, request *pb.LearnMsg) {
	srv.acceptor.Learn(ctx, request)
}
