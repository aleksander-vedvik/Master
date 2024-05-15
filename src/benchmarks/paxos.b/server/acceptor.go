package server

import (
	"errors"
	"log/slog"
	"sync"

	pb "github.com/aleksander-vedvik/benchmark/paxos.b/proto"

	"github.com/relab/gorums"
)

type Acceptor struct {
	mut         sync.Mutex
	rnd         uint32 // current round
	maxSeenSlot uint32
	slots       map[uint32]*pb.PromiseSlot // slots: is the internal data structure maintained by the acceptor to remember the slots
	learntVals  map[uint32]*pb.LearnMsg    // slots: is the internal data structure maintained by the acceptor to remember the slots
	numPeers    int
	senders     map[uint64]int
	adu         uint32
	cache       []struct {
		slot         uint32
		sendResponse func()
	}
}

func NewAcceptor(numPeers int) *Acceptor {
	return &Acceptor{
		numPeers:   numPeers,
		slots:      make(map[uint32]*pb.PromiseSlot),
		learntVals: make(map[uint32]*pb.LearnMsg),
		senders:    make(map[uint64]int),
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
		slog.Info("accept: failed", "reason", "rnd", "a.rnd", a.rnd, "request.Rnd", request.Rnd)
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
			slog.Info("accept: failed", "reason", "slot rnd")
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
	a.mut.Lock()
	defer a.mut.Unlock()
	md := broadcast.GetMetadata()
	if a.quorum(md.BroadcastID) {
		if prev, ok := a.learntVals[request.Slot]; !ok {
			a.learntVals[request.Slot] = request
		} else {
			if prev.Rnd < request.Rnd {
				a.learntVals[request.Slot] = request
			}
		}
		a.execute(request.Slot, broadcast, &pb.PaxosResponse{})
	}
}

func (a *Acceptor) execute(slot uint32, broadcast *pb.Broadcast, resp *pb.PaxosResponse) {
	if a.adu > slot {
		slog.Info("commit: failed", "reason", "adu")
		// old message
		return
	}
	if a.adu < slot {
		for i, c := range a.cache {
			if c.slot < slot {
				continue
			}
			tmp := append(a.cache[:i], struct {
				slot         uint32
				sendResponse func()
			}{slot, func() {
				broadcast.SendToClient(resp, nil)
			}})
			if i+1 >= len(a.cache) {
				a.cache = tmp
				return
			}
			a.cache = append(tmp, a.cache[i+1:]...)
		}
		return
	}
	if a.adu == slot {
		broadcast.SendToClient(resp, nil)
		a.adu++
	}
	for i, c := range a.cache {
		if c.slot <= a.adu {
			c.sendResponse()
			a.adu++
		} else {
			if i < len(a.cache) {
				a.cache = a.cache[i:]
			}
			return
		}
	}
}

// checks how many msgs the server has received for the given broadcastID.
// this method is only used in Learn, and does not require more advanced
// logic.
func (a *Acceptor) quorum(broadcastID uint64) bool {
	a.senders[broadcastID]++
	return int(a.senders[broadcastID]) > a.numPeers/2
}
