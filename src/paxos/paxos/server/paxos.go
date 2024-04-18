package server

import (
	pb "paxos/proto"

	"github.com/relab/gorums"
)

func (srv *PaxosServer) Prepare(ctx gorums.ServerCtx, req *pb.PrepareMsg) (*pb.PromiseMsg, error) {
	//slog.Info("received prepare", "srv", srv.addr)
	srv.mu.Lock()
	defer srv.mu.Unlock()
	srv.numMsgs[1]++
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
	//slog.Info("received accept", "srv", srv.addr)
	srv.mu.Lock()
	defer srv.mu.Unlock()
	srv.numMsgs[2]++
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
	//slog.Info("received learn", "srv", srv.addr)
	srv.mu.Lock()
	defer srv.mu.Unlock()
	srv.numMsgs[3]++
	md := broadcast.GetMetadata()
	if srv.quorum(md.BroadcastID) {
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
		//slog.Info(fmt.Sprintf("server(%v): commited", srv.id), "val", request.Val.Val)
	}
}

// checks how many msgs the server has received for the given broadcastID.
// this method is only used in Learn, and does not require more advanced
// logic.
func (srv *PaxosServer) quorum(broadcastID uint64) bool {
	srv.senders[broadcastID]++
	return int(srv.senders[broadcastID]) > len(srv.peers)/2
}

func (srv *PaxosServer) Ping(ctx gorums.ServerCtx, request *pb.Heartbeat) {
	srv.leaderElection.Ping(request.GetId())
}
