package server

import (
	pb "github.com/aleksander-vedvik/benchmark/paxos.b/proto"

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

/*import (
	pb "github.com/aleksander-vedvik/benchmark/paxos/proto"

	"github.com/relab/gorums"
)

func (srv *Server) Prepare(ctx gorums.ServerCtx, req *pb.PrepareMsg) (*pb.PromiseMsg, error) {
	//slog.Info("received prepare", "srv", srv.addr)
	srv.mut.Lock()
	defer srv.mut.Unlock()
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

func (srv *Server) Accept(ctx gorums.ServerCtx, request *pb.AcceptMsg, broadcast *pb.Broadcast) {
	//slog.Info("received accept", "srv", srv.addr)
	srv.mut.Lock()
	defer srv.mut.Unlock()
	// do not accept any messages with a rnd less than current
	if request.Rnd < srv.rnd {
		return
	}
	// set the current rnd to the highest it has seen
	srv.rnd = request.Rnd

	if slot, ok := srv.slots[request.Slot]; ok {
		// return if a slot with the same rnd already exists.
		if slot.Rnd == request.Rnd {
			return
		}
		// if all servers have commited this value it is considered final.
		//if slot.Final {
		//return
		//}
	}
	srv.slots[request.Slot] = &pb.PromiseSlot{
		Slot:  request.Slot,
		Rnd:   request.Rnd,
		Value: request.Val,
		//Final: false,
	}
	broadcast.Learn(&pb.LearnMsg{
		Rnd:  srv.rnd,
		Slot: request.Slot,
		Val:  request.Val,
	})
}

func (srv *Server) Learn(ctx gorums.ServerCtx, request *pb.LearnMsg, broadcast *pb.Broadcast) {
	//slog.Info("received learn", "srv", srv.addr)
	srv.mut.Lock()
	defer srv.mut.Unlock()
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
func (srv *Server) quorum(broadcastID uint64) bool {
	srv.senders[broadcastID]++
	return int(srv.senders[broadcastID]) > len(srv.peers)/2
}
*/
