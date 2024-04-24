package server

import pb "github.com/aleksander-vedvik/benchmark/paxos/proto"

type QSpec struct {
	qsize int
}

func newQSpec(qsize int) pb.QuorumSpec {
	return &QSpec{qsize: qsize}
}

func (q *QSpec) PrepareQF(in *pb.PrepareMsg, replies map[uint32]*pb.PromiseMsg) (*pb.PromiseMsg, bool) {
	if len(replies) < q.qsize {
		return nil, false
	}
	promise := &pb.PromiseMsg{}
	allSlots := make([][]*pb.PromiseSlot, 0)
	for _, reply := range replies {
		if reply.Rnd > promise.Rnd {
			promise.Rnd = reply.Rnd
		}
		allSlots = append(allSlots, reply.Slots)
	}
	addSlots(promise, allSlots)
	return promise, true
}

func addSlots(promiseMsg *pb.PromiseMsg, allSlots [][]*pb.PromiseSlot) {
	if len(allSlots) <= 0 {
		return
	}
	//added := make(map[uint32]*pb.PromiseSlot)
	added := make([]*pb.PromiseSlot, 0, len(allSlots[0]))
	allAdded := false
	for i := 0; !allAdded; i++ {
		allAdded = true
		for _, slots := range allSlots {
			if i >= len(slots) {
				continue
			}
			p := slots[i]
			switch {
			case len(added) <= i:
				added = append(added, p)
			case p == nil:
			default:
				elem := added[i]
				if elem == nil || p.Rnd > elem.Rnd {
					added[i] = p
				}
			}
			//if elem, ok := added[p.Slot]; ok {
			//if p.Rnd > elem.Rnd {
			//added[p.Slot] = p
			//}
			//} else {
			//added[p.Slot] = p
			//}
			allAdded = false
		}
	}
	promiseMsg.Slots = added
}

// not used by the server
func (q *QSpec) WriteQF(in *pb.PaxosValue, replies []*pb.PaxosResponse) (*pb.PaxosResponse, bool) {
	return nil, true
}

// not used by the server
func (q *QSpec) BenchmarkQF(in *pb.Empty, replies map[uint32]*pb.Result) (*pb.Result, bool) {
	return nil, true
}
