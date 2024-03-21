package server

import (
	"context"
	"log/slog"
	pb "paxos/proto"
	"time"

	"github.com/relab/gorums"
)

func (srv *PaxosServer) runPhaseOne() {
	srv.proposerMutex.Lock()
	defer srv.proposerMutex.Unlock()
	srv.mu.Lock()
	srv.proposerCtx, srv.cancelProposer = context.WithCancel(context.Background())
	srv.mu.Unlock()
start:
	select {
	case <-srv.proposerCtx.Done():
		return
	default:
	}
	if !srv.isLeader() {
		return
	}
	slog.Info("phase one: starting...")
	srv.setNewRnd()
	promiseMsg, err := srv.View.Prepare(context.Background(), &pb.PrepareMsg{
		Rnd:  srv.rnd,
		Slot: srv.maxSeenSlot,
	})
	if err != nil {
		select {
		case <-time.After(5 * time.Second):
			goto start
		case <-srv.proposerCtx.Done():
			return
		}
	}
	maxSlot := srv.maxSeenSlot
	for _, slot := range promiseMsg.Slots {
		if slot.Slot > maxSlot {
			maxSlot = slot.Slot
		}
		if s, ok := srv.slots[slot.Slot]; ok {
			if s.Final {
				continue
			}
		}
		srv.slots[slot.Slot] = slot
	}
	srv.maxSeenSlot = maxSlot
	slog.Info("phase one: finished")
	srv.runPhaseTwo()
}

func (srv *PaxosServer) runPhaseTwo() {
	slog.Info("phase two: started...")
	defer slog.Info("phase two: finished")
	for {
		select {
		case <-srv.proposerCtx.Done():
			return
		default:
		}
		srv.mu.Lock()
		for _, req := range srv.clientReqs {
			srv.maxSeenSlot++
			srv.BroadcastAccept(&pb.AcceptMsg{
				Rnd:  srv.rnd,
				Slot: srv.maxSeenSlot,
				Val:  req.message,
			}, req.broadcastID)
		}
		srv.clientReqs = make([]*clientReq, 0)
		srv.mu.Unlock()
	}
}

func (srv *PaxosServer) setNewRnd() {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	numSrvs := uint32(len(srv.peers))
	srv.rnd -= srv.rnd % numSrvs
	srv.rnd += srv.id + numSrvs
}

func (srv *PaxosServer) Prepare(ctx gorums.ServerCtx, req *pb.PrepareMsg) (*pb.PromiseMsg, error) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if req.Rnd < srv.rnd {
		return nil, nil
	}
	if req.Slot > srv.maxSeenSlot {
		srv.maxSeenSlot = req.Slot
	}
	// req.Rnd will always be higher or equal to srv.rnd
	srv.rnd = req.Rnd
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
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if request.Rnd < srv.rnd {
		return
	}
	srv.rnd = request.Rnd

	if slot, ok := srv.slots[request.Slot]; ok {
		if slot.Final {
			return
		}
	}
	srv.slots[request.Slot] = &pb.PromiseSlot{
		Slot:  request.Slot,
		Rnd:   request.Rnd,
		Value: request.Val,
		Final: false,
	}
	broadcast.Learn(&pb.LearnMsg{
		Rnd:  srv.rnd,
		Slot: request.Slot,
		Val:  request.Val,
	})
}

func (srv *PaxosServer) Learn(ctx gorums.ServerCtx, request *pb.LearnMsg, broadcast *pb.Broadcast) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	md := broadcast.GetMetadata()
	if srv.quorum(md.Count) {
		if s, ok := srv.slots[request.Slot]; ok {
			if s.Final {
				return
			}
			s.Final = true
		} else {
			srv.slots[request.Slot] = &pb.PromiseSlot{
				Slot:  request.Slot,
				Rnd:   request.Rnd,
				Value: request.Val,
				Final: true,
			}
		}
		broadcast.SendToClient(&pb.PaxosResponse{}, nil)
		slog.Info("commited", "val", request.Val.Val)
	}
}

func (srv *PaxosServer) quorum(count uint64) bool {
	return int(count) > len(srv.peers)/2
}

func (srv *PaxosServer) Ping(ctx gorums.ServerCtx, request *pb.Heartbeat) {
	srv.leaderElection.Ping(request.GetId())
}
