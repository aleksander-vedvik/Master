package server

import (
	pb "github.com/aleksander-vedvik/benchmark/paxosqcb/proto"
)

// Acceptor represents an acceptor as defined by the Multi-Paxos algorithm.
type Acceptor struct {
	rnd         Round               // highest round the acceptor has promised in.
	accepted    map[Slot]*pb.PValue // map of accepted values for each slot.
	highestSeen Slot                // highest slot for which a prepare has been received.
}

// NewAcceptor returns a new Multi-Paxos acceptor
func NewAcceptor() *Acceptor {
	return &Acceptor{
		rnd:      NoRound,
		accepted: make(map[Slot]*pb.PValue),
	}
}

// handlePrepare processes the prepare according to the Multi-Paxos algorithm,
// returning a promise, or nil if the prepare should be ignored.
func (a *Acceptor) handlePrepare(prepare *pb.PrepareMsg) (prm *pb.PromiseMsg) {
	if prepare.Crnd < a.rnd {
		return nil
	}
	if prepare.Slot == NoSlot {
		return &pb.PromiseMsg{Rnd: prepare.Crnd}
	}
	a.rnd = prepare.Crnd
	var accepted []*pb.PValue
	for slot := prepare.Slot; slot <= a.highestSeen; slot++ {
		pval := a.getPValue(slot)
		if pval.Vrnd > NoRound {
			accepted = append(accepted, pval)
		}
	}
	return &pb.PromiseMsg{
		Rnd:      a.rnd,
		Accepted: accepted,
	}
}

// handleAccept processes the accept according to the Multi-Paxos algorithm,
// returning a learn, or nil if the accept should be ignored.
func (a *Acceptor) handleAccept(accept *pb.AcceptMsg) (lrn *pb.LearnMsg) {
	pval := a.getPValue(accept.Slot)

	if accept.Rnd < a.rnd || accept.Rnd == pval.Vrnd {
		return nil
	}

	if accept.Rnd > a.rnd {
		a.rnd = accept.Rnd
	}

	pval.Vrnd = accept.Rnd
	pval.Vval = accept.Val

	return &pb.LearnMsg{
		Slot: pval.Slot,
		Rnd:  pval.Vrnd,
		Val:  pval.Vval,
	}
}

// getPValue returns the accepted PValue for the given slot;
// if pval has been accepted for the slot, it creates a new one.
func (a *Acceptor) getPValue(slot Slot) *pb.PValue {
	pval, found := a.accepted[slot]
	if found {
		return pval
	}
	pval = &pb.PValue{Slot: slot, Vrnd: NoRound}
	a.accepted[slot] = pval
	if pval.Slot > a.highestSeen {
		a.highestSeen = pval.Slot
	}
	return pval
}
