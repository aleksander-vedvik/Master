package server

import (
	pb "paxos/proto"

	"github.com/relab/gorums"
)

func (srv *PaxosServer) Prepare(ctx gorums.ServerCtx, req *pb.PrepareMsg) (*pb.PromiseMsg, error) {
	if req.Crnd < srv.rnd {
		return nil, nil
	}
	if req.Slot > srv.maxSeenSlot {
		srv.maxSeenSlot = req.Slot
	}
	if req.Crnd > srv.rnd {
		srv.rnd = req.Crnd
	}
	if len(srv.slots) == 0 {
		return &pb.PromiseMsg{Rnd: srv.rnd, Slots: nil}, nil
	}
	promiseSlots := make([]*pb.PromiseSlot, 0, len(srv.slots))
	for _, promiseSlot := range srv.slots {
		if promiseSlot.Slot >= srv.maxSeenSlot {
			promiseSlots = append(promiseSlots, promiseSlot)
		}
	}
	return &pb.PromiseMsg{Rnd: srv.rnd, Slots: promiseSlots}, nil
}

func (srv *PaxosServer) Accept(ctx gorums.ServerCtx, request *pb.AcceptMsg, broadcast *pb.Broadcast) {
	if request.Rnd < srv.rnd {
		return
	}
	srv.rnd = request.Rnd

	if request.Slot >= srv.maxSeenSlot {
		srv.slots[request.Slot] = &pb.PromiseSlot{
			Slot:  request.Slot,
			Vrnd:  request.Rnd,
			Value: request.Val,
		}
	}
	broadcast.Learn(&pb.LearnMsg{
		Rnd:  srv.rnd,
		Slot: request.Slot,
		Val:  request.Val,
	})
}

func (srv *PaxosServer) Learn(ctx gorums.ServerCtx, request *pb.LearnMsg, broadcast *pb.Broadcast) {
	if srv.quorum() {
		broadcast.SendToClient(&pb.Response{}, nil)
	}
}

func (srv *PaxosServer) quorum() bool {
	return false
}

func (srv *PaxosServer) Ping(ctx gorums.ServerCtx, request *pb.Heartbeat) {
	srv.leaderElection.Ping(request)
}
