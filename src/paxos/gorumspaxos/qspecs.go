//go:build !solution

package gorumspaxos

import (
	pb "paxos/proto"
	"sort"
)

// PaxosQSpec is a quorum specification object for Paxos.
// It only holds the quorum size.
// DONE(student): If necessary add additional fields
// DO NOT remove the existing fields in the structure
type PaxosQSpec struct {
	qSize int
}

// NewPaxosQSpec returns a quorum specification object for Paxos
// for the given quorum size.
func NewPaxosQSpec(quorumSize int) pb.QuorumSpec {
	return &PaxosQSpec{
		qSize: quorumSize,
	}
}

// PrepareQF is the quorum function to process the replies of the Prepare RPC call.
// Proposer handle PromiseMsgs returned by the Acceptors, and any promise slots added
// to the PromiseMsg should be in the increasing order.
func (qs PaxosQSpec) PrepareQF(prepare *pb.PrepareMsg, replies map[uint32]*pb.PromiseMsg) (*pb.PromiseMsg, bool) {
	// DONE(student) complete the PrepareQF
	// New check for quorum, makes more sense
	if len(replies) < qs.qSize {
		return nil, false
	}

	Crnd := prepare.Crnd
	adu := prepare.Slot

	// BM: Aggregate all the slots together
	allSlots := make([]*pb.PromiseSlot, 0)
	for replyIndex := range replies {
		// BM: Only aggregate the slots with the correct Round-ID
		if replies[replyIndex].Rnd != nil && replies[replyIndex].Rnd.Id == Crnd.Id {
			allSlots = append(allSlots, replies[replyIndex].Slots...)
		}
	}

	resultSlots := sortSlots(allSlots, adu, Crnd)
	resultPromise := &pb.PromiseMsg{Rnd: Crnd, Slots: resultSlots}
	return resultPromise, true
}

/*
BM:
Helper function which sorts all the correct slots into an array of Promiseslots ready to be inserted into a PromiseMsg
Overall logic very similar to aggregateAccepts from lab4
*/
func sortSlots(allSlots []*pb.PromiseSlot, adu *pb.Slot, Crnd *pb.Round) []*pb.PromiseSlot {
	sort.SliceStable(allSlots, func(i, j int) bool {
		return allSlots[i].Slot.Id < allSlots[j].Slot.Id
	})

	if len(allSlots) <= 0 {
		return make([]*pb.PromiseSlot, 0)
	}

	minSlot := allSlots[0].Slot.Id
	maxSlot := allSlots[len(allSlots)-1].Slot.Id

	for minSlot <= adu.Id {
		minSlot += 1
	}

	resultSlotsLength := int(maxSlot) - int(minSlot) + 1
	if resultSlotsLength < 0 {
		resultSlotsLength = 0
	}
	resultSlots := make([]*pb.PromiseSlot, resultSlotsLength)

	for resultIndex := range resultSlots {

		// BM: The base values for each Slot, Vrnd and Value change if this slot maps to a reply
		resultSlots[resultIndex] = &pb.PromiseSlot{Slot: &pb.Slot{Id: minSlot + uint32(resultIndex)}, Vrnd: &pb.Round{Id: -1}}

		// BM: Initialize the base values, and update them if any fit the criteria
		thisRound := Crnd
		thisValue := &pb.Value{IsNoop: true}
		for _, slot := range allSlots {

			if slot.Slot.Id == resultSlots[resultIndex].Slot.Id && slot.Vrnd.Id > resultSlots[resultIndex].Vrnd.Id {
				thisRound = slot.Vrnd
				thisValue = slot.Value
				resultSlots[resultIndex].Vrnd.Id = slot.Vrnd.Id
			}
		}

		resultSlots[resultIndex].Vrnd = thisRound
		resultSlots[resultIndex].Value = thisValue
	}

	return resultSlots
}

// AcceptQF is the quorum function for the Accept quorum RPC call
// This is where the Proposer handle LearnMsgs to determine if a
// value has been decided by the Acceptors.
// The quorum function returns true if a value has been decided,
// and the corresponding LearnMsg holds the round number and value
// that was decided. If false is returned, no value was decided.
func (qs PaxosQSpec) AcceptQF(accMsg *pb.AcceptMsg, replies map[uint32]*pb.LearnMsg) (*pb.LearnMsg, bool) {
	// DONE (student) complete the AcceptQF
	if len(replies) < qs.qSize {
		return nil, false
	}
	validReplies := make([]*pb.LearnMsg, 0)
	for _, reply := range replies {
		if checkIsRoundEqual(accMsg.Rnd, reply.Rnd) && checkIsSlotEqual(reply.Slot, accMsg.Slot) && checkIsValueEqual(reply.Val, accMsg.Val) {
			validReplies = append(validReplies, reply)
		}
	}
	if len(validReplies) >= qs.qSize {
		return validReplies[0], true
	}

	return nil, false
}

// CommitQF is the quorum function for the Commit quorum RPC call.
// This function just waits for a quorum of empty replies,
// indicating that at least a quorum of Learners have committed
// the value decided by the Acceptors.
func (qs PaxosQSpec) CommitQF(_ *pb.LearnMsg, replies map[uint32]*pb.Empty) (*pb.Empty, bool) {
	// DONE (student) complete the CommitQF
	if len(replies) >= qs.qSize {
		return nil, true
	}
	return nil, false
}

// ClientHandleQF is the quorum function  for the ClientHandle quorum RPC call.
// This functions waits for replies from the majority of replicas. Received replies
// should be validated before returning the response
func (qs PaxosQSpec) ClientHandleQF(in *pb.Value, replies map[uint32]*pb.Response) (*pb.Response, bool) {
	// DONE (student) complete the ClientHandleQF
	if len(replies) < qs.qSize {
		return nil, false
	}
	validReplies := make([]*pb.Response, 0)
	for _, reply := range replies {
		if checkIsResponseValid(reply, in) {
			validReplies = append(validReplies, reply)
		}
	}
	if len(validReplies) >= qs.qSize {
		return validReplies[0], true
	}
	return nil, false
}
