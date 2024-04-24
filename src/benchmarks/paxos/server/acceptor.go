package server

import (
	"errors"
	"sync"

	pb "github.com/aleksander-vedvik/benchmark/paxos/proto"

	"github.com/relab/gorums"
)

type Acceptor struct {
	mut         sync.Mutex
	rnd         uint32 // current round
	maxSeenSlot uint32
	slots       map[uint32]*pb.PromiseSlot // slots: is the internal data structure maintained by the acceptor to remember the slots
	numPeers    int
	senders     map[uint64]int
}

func NewAcceptor(numPeers int) *Acceptor {
	return &Acceptor{
		numPeers: numPeers,
		slots:    make(map[uint32]*pb.PromiseSlot),
		senders:  make(map[uint64]int),
	}
}

func (a *Acceptor) Prepare(ctx gorums.ServerCtx, req *pb.PrepareMsg) (*pb.PromiseMsg, error) {
	//slog.Info("received prepare", "srv", srv.addr)
	ctx.Release() // to prevent deadlocks
	a.mut.Lock()
	defer a.mut.Unlock()
	if req.Rnd < a.rnd {
		return nil, errors.New("ignored")
	}
	//if req.Slot > a.maxSeenSlot {
	//a.maxSeenSlot = req.Slot
	//}
	// req.Rnd will always be higher or equal to srv.rnd
	a.rnd = req.Rnd
	if len(a.slots) == 0 {
		return &pb.PromiseMsg{Rnd: a.rnd, Slots: nil}, nil
	}
	promiseSlots := make([]*pb.PromiseSlot, a.maxSeenSlot)
	for slotNo, promiseSlot := range a.slots {
		promiseSlots[slotNo] = promiseSlot
		//if promiseSlot.Slot >= a.maxSeenSlot {
		//promiseSlots = append(promiseSlots, promiseSlot)
		//}
	}
	return &pb.PromiseMsg{Rnd: a.rnd, Slots: promiseSlots}, nil
}

func (a *Acceptor) Accept(ctx gorums.ServerCtx, request *pb.AcceptMsg, broadcast *pb.Broadcast) {
	//slog.Info("received accept", "srv", srv.addr)
	a.mut.Lock()
	defer a.mut.Unlock()
	// do not accept any messages with a rnd less than current
	if request.Rnd < a.rnd {
		return
	}
	// set the current rnd to the highest it has seen
	a.rnd = request.Rnd
	if request.Slot > a.maxSeenSlot {
		a.maxSeenSlot = request.Slot
	}

	if slot, ok := a.slots[request.Slot]; ok {
		// return if a slot with the same rnd already exists.
		if slot.Rnd == request.Rnd {
			return
		}
		// if all servers have commited this value it is considered final.
		//if slot.Final {
		//return
		//}
	}
	a.slots[request.Slot] = &pb.PromiseSlot{
		Slot:  request.Slot,
		Rnd:   request.Rnd,
		Value: request.Val,
		//Final: false,
	}
	broadcast.Learn(&pb.LearnMsg{
		Rnd:  a.rnd,
		Slot: request.Slot,
		Val:  request.Val,
	})
}

func (a *Acceptor) Learn(ctx gorums.ServerCtx, request *pb.LearnMsg, broadcast *pb.Broadcast) {
	//slog.Info("received learn", "srv", srv.addr)
	a.mut.Lock()
	defer a.mut.Unlock()
	md := broadcast.GetMetadata()
	if a.quorum(md.BroadcastID) {
		if _, ok := a.slots[request.Slot]; ok {
			return
		} else {
			a.slots[request.Slot] = &pb.PromiseSlot{
				Slot:  request.Slot,
				Rnd:   request.Rnd,
				Value: request.Val,
			}
		}
		broadcast.SendToClient(&pb.PaxosResponse{}, nil)
		//slog.Info(fmt.Sprintf("server(%v): commited", srv.id), "val", request.Val.Val)
	}
}

// checks how many msgs the server has received for the given broadcastID.
// this method is only used in Learn, and does not require more advanced
// logic.
func (a *Acceptor) quorum(broadcastID uint64) bool {
	a.senders[broadcastID]++
	return int(a.senders[broadcastID]) > a.numPeers/2
}
