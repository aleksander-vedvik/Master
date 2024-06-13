package server

import (
	"context"
	"errors"
	"sync"

	pb "github.com/aleksander-vedvik/benchmark/paxosqcb/proto"
)

// Proposer represents a proposer as defined by the Multi-Paxos algorithm.
type Proposer struct {
	mu                 sync.RWMutex
	id                 int               // replica's id.
	leader             int               // current Paxos leader.
	crnd               Round             // replica's current round; initially, this replica's id.
	adu                Slot              // all-decided-up-to is the highest consecutive slot that has been committed.
	nextSlot           Slot              // slot for the next request, initially 0.
	phaseOneDone       bool              // indicates if the phase1 is done, initially false.
	config             MultiPaxosConfig  // configuration used for multipaxos.
	nodeMap            map[string]uint32 // map of the address to the node id.
	acceptMsgQueue     []*pb.AcceptMsg   // queue of pending accept messages as part of prepare operation.
	clientRequestQueue []*pb.AcceptMsg   // queue of pending client requests.
	msgQueue           chan *msg
}

type msg struct {
	accept    *pb.AcceptMsg
	broadcast *pb.Broadcast
}

// NewProposer returns a new Multi-Paxos proposer with the specified
// replica id, initial leader, and nodeMap.
func NewProposer(myID, leader int, nodeMap map[string]uint32) *Proposer {
	propIdx := myIndex(myID, nodeMap)
	return &Proposer{
		id:                 myID,
		leader:             leader,
		nodeMap:            nodeMap,
		phaseOneDone:       true, // do not use leader election
		crnd:               Round(propIdx),
		adu:                0,
		acceptMsgQueue:     make([]*pb.AcceptMsg, 0),
		clientRequestQueue: make([]*pb.AcceptMsg, 0),
		msgQueue:           make(chan *msg, 50),
	}
}

// isLeader returns true if this replica is the leader.
func (p *Proposer) isLeader() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.leader == p.id
}

// Perform the accept quorum call on the replicas.
//
//  1. Check if any pending accept requests in the acceptReqQueue to process
//  2. Check if any pending client requests in the clientRequestQueue to process
//  3. Increment the nextSlot and prepare an accept message for the pending request,
//     using crnd and nextSlot.
//  4. Perform accept quorum call on the configuration and return the learnMsg.
func (p *Proposer) performAccept(accept *pb.AcceptMsg) (*pb.LearnMsg, error) {
	ctx := context.Background()
	// all quorum calls should happen without holding locks, otherwise
	// leader may not be able to process the its own RPC call.
	return p.config.Accept(ctx, accept)
}

// Perform the commit operation using a multicast call.
func (p *Proposer) performCommit(learn *pb.LearnMsg, broadcast *pb.Broadcast) error {
	if learn == nil {
		return errors.New("no learn message to send")
	}
	broadcast.Commit(learn)
	return nil
}

// setConfiguration set the configuration for the proposer, allowing it to
// communicate with the other replicas.
func (p *Proposer) setConfiguration(config MultiPaxosConfig) {
	p.mu.Lock()
	p.config = config
	p.mu.Unlock()
}

// ProcessRequest: processes the request from the client by putting it into the
// clientRequestQueue and calling the getResponse to get the response matching
// the request.
func (p *Proposer) AddRequestToQ(request *pb.Value, broadcast *pb.Broadcast) {
	if p.isLeader() {
		accept := &pb.AcceptMsg{Val: request}
		p.msgQueue <- &msg{
			accept:    accept,
			broadcast: broadcast,
		}
	}
}
