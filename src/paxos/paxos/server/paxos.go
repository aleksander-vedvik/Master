package server

import (
	"context"
	pb "paxos/proto"

	"github.com/relab/gorums"
)

func (srv *PaxosServer) runPhaseOne() {
	srv.rnd = srv.pickNext()
	_, err := srv.View.Prepare(context.Background(), &pb.PrepareMsg{
		Rnd:  srv.rnd,
		Slot: srv.maxSeenSlot,
	})
	if err != nil {
		return
	}
}

func (srv *PaxosServer) runPhaseTwo() {

}

func (srv *PaxosServer) pickNext() uint32 {
	return 0
}

func (srv *PaxosServer) Prepare(ctx gorums.ServerCtx, req *pb.PrepareMsg) (*pb.PromiseMsg, error) {
	if req.Rnd < srv.rnd {
		return nil, nil
	}
	if req.Slot > srv.maxSeenSlot {
		srv.maxSeenSlot = req.Slot
	}
	if req.Rnd > srv.rnd {
		srv.rnd = req.Rnd
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
	md := broadcast.GetMetadata()
	if srv.quorum(md.Count) {
		broadcast.SendToClient(&pb.Response{}, nil)
	}
}

func (srv *PaxosServer) quorum(count uint64) bool {
	return int(count) > len(srv.peers)/2
}

func (srv *PaxosServer) Ping(ctx gorums.ServerCtx, request *pb.Heartbeat) {
	srv.leaderElection.Ping(request)
}
