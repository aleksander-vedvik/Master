package server

import (
	"errors"
	"sort"
	"sync"

	pb "github.com/aleksander-vedvik/benchmark/paxos.bc/proto"

	"github.com/relab/gorums"
)

type Acceptor struct {
	mut         sync.Mutex
	addr        string
	rnd         uint32 // current round
	maxSeenSlot uint32
	slots       map[uint32]*pb.PromiseSlot // slots: is the internal data structure maintained by the acceptor to remember the slots
	learntVals  map[uint32]*pb.LearnMsg    // slots: is the internal data structure maintained by the acceptor to remember the slots
	numPeers    int
	senders     map[uint64]int
	adu         uint32
	cache       map[uint32]struct {
		slot         uint32
		sendResponse func()
	}
}

func NewAcceptor(addr string, numPeers int) *Acceptor {
	return &Acceptor{
		addr:       addr,
		numPeers:   numPeers,
		slots:      make(map[uint32]*pb.PromiseSlot),
		learntVals: make(map[uint32]*pb.LearnMsg),
		senders:    make(map[uint64]int),
		cache: make(map[uint32]struct {
			slot         uint32
			sendResponse func()
		}),
	}
}

func (a *Acceptor) Prepare(ctx gorums.ServerCtx, req *pb.PrepareMsg) (*pb.PromiseMsg, error) {
	a.mut.Lock()
	defer a.mut.Unlock()
	if req.Rnd < a.rnd {
		return nil, errors.New("ignored")
	}
	a.rnd = req.Rnd
	if len(a.slots) == 0 {
		return &pb.PromiseMsg{Rnd: a.rnd, Slots: nil}, nil
	}
	promiseSlots := make([]*pb.PromiseSlot, a.maxSeenSlot)
	for slotNo, promiseSlot := range a.slots {
		promiseSlots[slotNo] = promiseSlot
	}
	return &pb.PromiseMsg{Rnd: a.rnd, Slots: promiseSlots}, nil
}

func (a *Acceptor) Accept(ctx gorums.ServerCtx, request *pb.AcceptMsg, broadcast *pb.Broadcast) {
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
	}
	a.slots[request.Slot] = &pb.PromiseSlot{
		Slot:  request.Slot,
		Rnd:   request.Rnd,
		Value: request.Val,
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
		// old message
		return
	}
	if a.adu < slot {
		a.cache[slot] = struct {
			slot         uint32
			sendResponse func()
		}{slot, func() {
			broadcast.SendToClient(resp, nil)
		}}
		return
	}
	if a.adu == slot {
		broadcast.SendToClient(resp, nil)
		a.adu++
	}
	send := make([]struct {
		slot         uint32
		sendResponse func()
	}, 0, len(a.cache))
	for _, c := range a.cache {
		send = append(send, c)
	}
	sort.Slice(send, func(i, j int) bool {
		return send[i].slot < send[j].slot
	})
	for _, cachedSlot := range send {
		if cachedSlot.slot <= a.adu {
			cachedSlot.sendResponse()
			a.adu++
			delete(a.cache, cachedSlot.slot)
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
