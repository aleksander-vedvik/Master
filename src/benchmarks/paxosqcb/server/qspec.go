package server

import (
	"slices"

	pb "github.com/aleksander-vedvik/benchmark/paxosqcb/proto"
)

// PaxosQSpec is a quorum specification object for Paxos.
// It only holds the quorum size.
type PaxosQSpec struct {
	quorum int
}

// NewPaxosQSpec returns a quorum specification object for Paxos
// for the given configuration size n.
func NewPaxosQSpec(n int) PaxosQSpec {
	return PaxosQSpec{quorum: (n-1)/2 + 1}
}

// PrepareQF is the quorum function to process the replies from the Prepare quorum call.
// This is where the Proposer handle PromiseMsgs returned by the Acceptors, and any
// Accepted values in the promise replies should be combined into the returned PromiseMsg
// in the increasing slot order with any gaps filled with no-ops. Invalid promise messages
// should be ignored. For example, quorum function should only process promise replies that
// belong to the current round, as indicated by the prepare message. The quorum function
// returns true if a quorum of valid promises was found, and the combined PromiseMsg.
// Nil and false is returned if no quorum of valid promises was found.
func (qs PaxosQSpec) PrepareQF(prepare *pb.PrepareMsg, replies map[uint32]*pb.PromiseMsg) (*pb.PromiseMsg, bool) {
	if len(replies) < qs.quorum {
		return nil, false
	}
	// keep only valid validPromises that match the prepare message
	validPromises := make([]*pb.PromiseMsg, 0, len(replies))
	for _, prm := range replies {
		if prepare.IsValid(prm) {
			validPromises = append(validPromises, prm)
		}
	}
	if len(validPromises) < qs.quorum {
		return nil, false
	}

	pvalMap := make(map[Slot]*pb.PValue)
	for _, prm := range validPromises {
		for _, pval := range prm.Accepted {
			if pval.Slot > prepare.Slot {
				alreadySeenPVal, found := pvalMap[pval.Slot]
				if !found {
					pvalMap[pval.Slot] = pval
					continue
				}
				if alreadySeenPVal.Vrnd < pval.Vrnd {
					pvalMap[pval.Slot] = pval
				}
			}
		}
	}

	accepted := make([]*pb.PValue, 0)
	if len(pvalMap) > 1 {
		// sort the slots in increasing order to ensure that gaps
		// in the slot sequence can be filled with no-ops.
		sortedSlots := make([]Slot, 0, len(pvalMap))
		for slot := range pvalMap {
			sortedSlots = append(sortedSlots, slot)
		}
		slices.Sort(sortedSlots)

		// fill any gaps in the slot sequence by adding no-ops.
		firstSlot, lastSeenSlot := sortedSlots[0], sortedSlots[len(sortedSlots)-1]
		for slot := firstSlot; slot <= lastSeenSlot; slot++ {
			pval, found := pvalMap[slot]
			if !found {
				// fill a gap in the slot sequence by adding a no-op.
				pval = &pb.PValue{
					Vrnd: prepare.Crnd,
					Vval: &pb.Value{IsNoop: true},
					Slot: slot,
				}
			}
			accepted = append(accepted, pval)
		}
	} else {
		for _, pval := range pvalMap {
			accepted = append(accepted, pval)
		}
	}
	return &pb.PromiseMsg{
		Rnd:      prepare.Crnd,
		Accepted: accepted,
	}, true
}

// AcceptQF is the quorum function to process the replies from the Accept quorum call.
// This is where the Proposer handle LearnMsgs to determine if a value has been decided
// by the Acceptors. The quorum function returns true if a value has been decided, and
// the corresponding LearnMsg holds the slot, round number and value that was decided.
// Nil and false is returned if no value was decided.
func (qs PaxosQSpec) AcceptQF(accept *pb.AcceptMsg, replies map[uint32]*pb.LearnMsg) (*pb.LearnMsg, bool) {
	validLearns := len(replies)
	if validLearns < qs.quorum {
		return nil, false
	}
	// remove the invalid learns
	var learn *pb.LearnMsg
	for _, lrn := range replies {
		if !accept.Match(lrn) {
			validLearns--
		} else {
			learn = lrn
		}
	}
	if validLearns < qs.quorum {
		return nil, false
	}
	return learn, true
}

// ClientHandleQF is the quorum function to process the replies from the ClientHandle quorum call.
// This is where the Client handle the replies from the replicas. The quorum function should
// validate the replies against the request, and only valid replies should be considered.
// The quorum function returns true if a quorum of the replicas replied with the same response,
// and a single response is returned. Nil and false is returned if no quorum of valid replies was found.
func (qs PaxosQSpec) ClientHandleQF(request *pb.Value, replies map[uint32]*pb.Response) (*pb.Response, bool) {
	validResponses := len(replies)
	if validResponses < qs.quorum {
		return nil, false
	}
	// remove the invalid responses
	var response *pb.Response
	for _, resp := range replies {
		if !request.Match(resp) {
			validResponses--
		} else {
			response = resp
		}
	}
	if validResponses < qs.quorum {
		return nil, false
	}
	return response, true
}
