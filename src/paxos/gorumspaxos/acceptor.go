//go:build !solution

package gorumspaxos

import (
	pb "paxos/proto"
)

// Acceptor represents an acceptor as defined by the Multi-Paxos algorithm.
// For current implementation alpha is assumed to be 1. In future if the alpha is
// increased, then there would be significant changes to the current implementation.
type Acceptor struct {
	rnd           *pb.Round                  // rnd: is the current round in which the acceptor is participating
	slots         map[uint32]*pb.PromiseSlot // slots: is the internal data structure maintained by the acceptor to remember the slots
	maxSeenSlotId uint32                     // maxSeenSlotId: is the highest slot for which the prepare is received
}

// NewAcceptor returns a new Multi-Paxos acceptor
func NewAcceptor() *Acceptor {
	return &Acceptor{
		rnd:           NoRound(),
		slots:         make(map[uint32]*pb.PromiseSlot),
		maxSeenSlotId: 0,
	}
}

// handlePrepare takes a prepare message and returns a promise message according to
// the Multi-Paxos algorithm. If the prepare message is invalid, nil should be returned.
func (a *Acceptor) handlePrepare(prp *pb.PrepareMsg) (prm *pb.PromiseMsg) {
	// DONE(student): Complete the handlePrepare
	if prp.Crnd.Id < a.rnd.Id {
		return
	}
	if prp.Slot != nil && prp.Slot.Id > a.maxSeenSlotId {
		a.maxSeenSlotId = prp.Slot.Id
	}
	if prp.Crnd.Id > a.rnd.Id {
		a.rnd.Id = prp.Crnd.Id
	}
	if len(a.slots) == 0 {
		return &pb.PromiseMsg{Rnd: a.rnd, Slots: nil}
	}
	promiseSlots := make([]*pb.PromiseSlot, 0, len(a.slots))
	for _, promiseSlot := range a.slots {
		if promiseSlot.Slot.Id >= a.maxSeenSlotId {
			promiseSlots = append(promiseSlots, promiseSlot)
		}
	}
	return &pb.PromiseMsg{Rnd: a.rnd, Slots: promiseSlots}
}

// handleAccept takes an accept message and returns a learn message according to
// the Multi-Paxos algorithm. If the accept message is invalid, nil should be returned.
func (a *Acceptor) handleAccept(acc *pb.AcceptMsg) (lrn *pb.LearnMsg) {
	// DONE(student): Complete the handleAccept
	if acc.Rnd.Id < a.rnd.Id {
		return
	}
	a.rnd.Id = acc.Rnd.Id
	if acc.Slot.Id >= a.maxSeenSlotId {
		a.slots[acc.Slot.Id] = &pb.PromiseSlot{
			Slot:  acc.Slot,
			Vrnd:  acc.Rnd,
			Value: acc.Val,
		}
	}
	return &pb.LearnMsg{
		Rnd:  a.rnd,
		Slot: acc.Slot,
		Val:  acc.Val,
	}
}
