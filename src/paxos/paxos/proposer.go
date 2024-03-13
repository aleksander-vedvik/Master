package main

/*
import (
	"container/list"
	"fmt"
	"paxos/leaderdetector"
	pb "paxos/proto"
	"sync"
	"time"

	"context"
)

type Operation int

const (
	NoAction Operation = iota
	Prepare
	Accept
	Commit
)

// Proposer represents a proposer as defined by the Multi-Paxos algorithm.
// DONE(student) Add additional fields if necessary
// DO NOT remove the existing fields in the structure
type Proposer struct {
	sync.RWMutex
	id                 int                           // id: id of the replica
	crndID             int32                         // crndID: holds the current round of the replica, initialized to 0
	aduSlotID          uint32                        // aduSlotID: is the slot id for which the commit is completed.
	nextSlotID         uint32                        // nextSlotID: is used in processing the requests, initialized to 0
	phaseOneDone       bool                          // phaseOneDone: indicates if the phase1 is done, initialized to false
	ld                 leaderdetector.LeaderDetector // ld: leader detector implementation module
	leader             int                           // leader: initial leader of the consensus
	nextOperation      Operation                     // nextOperation: holds the next operation to be performed
	phaseOneWaitTime   time.Duration                 // phaseOneWaitTime: duration until phase1 operation timeout
	phaseTwoWaitTime   time.Duration                 // phaseTwoWaitTime: duration until phase2 operation timeout
	acceptMsgQueue     *list.List                    // acceptMsgQueue: queue to hold the pending accept requests as part of prepare operation
	learnMsg           *pb.LearnMsg                  // learnMsg: next learnMsg to be processed. (may become a queue if alpha increased)
	clientRequestQueue []int                         // clientRequestQueue: queue used to store the pending client requests
	statusChannel      chan bool                     // statusChannel: used to pass the status of the operation to move to next operation
	stop               chan struct{}                 // stop: channel used to end the replica functionality

	ldChan <-chan int // Channel where leader updates are propagated
}

// NewProposer returns a new Multi-Paxos proposer.
func NewProposer() *Proposer {
	return &Proposer{
		stop:               make(chan struct{}),
		crndID:             0,
		nextSlotID:         0,
		clientRequestQueue: make([]int, 0),
		statusChannel:      make(chan bool, 1),
	}
}

// Internal: checks if the current replica is the leader
func (p *Proposer) isLeader() bool {
	p.RLock()
	defer p.RUnlock()
	return p.leader == p.id
}

// runPhaseOne runs the phase one (prepare->promise) of the protocol.
// This function should be called only if the replica is the leader and
// phase1 has not already completed.
// 1. Create the configuration with all replicas, call the createConfiguration method
// 2. Create a PrepareMsg with the current round ID and incremented aduSlotID
// 3. Send the prepare message to all the replicas (quorum call)
// 4. If succeeded set phaseOneDone to true and nextOperation to "accept"
// 5. If the response Promise Message contains the slots, prepare an accept
// message for each of the slot and add it to the accept queue
func (p *Proposer) runPhaseOne() error {
	if !p.isLeader() {
		return nil
	}

	prepareMSG := pb.PrepareMsg{Crnd: &pb.Round{Id: p.crndID}}
	p.IncrementAllDecidedUpTo()

	result, err := p.config.Prepare(context.Background(), &prepareMSG) // NOTE: Config now using p.qspec (DSH)
	if err != nil {
		return err
	}

	p.phaseOneDone = true
	p.nextOperation = Accept

	for _, slot := range result.Slots { // NOTE: It is here assumed that result will return an empty result.Slots so that this line does not create an error
		acceptMSG := pb.AcceptMsg{
			Rnd:  slot.Vrnd,
			Slot: slot.Slot,
			Val:  slot.Value,
		}
		p.acceptMsgQueue.PushBack(&acceptMSG) //NOTE: The list does not specify what it should receive, I thought it best to give it a pointer due to a warning about copying a lock
	}

	return nil
}

// Start starts proposer's main run loop as a separate goroutine.
// The separate goroutine should start an infinite loop and
// use a select mechanism to wait for the following events
// 1. on the status channel to conduct the paxos phases and operations
// 2. on the channel returned by leader detector subscribe method
// 3. on stop channel to stop the relpica
// default case is to process any pending client requests by adding to the
// clientRequestQueue and call runMultiPaxos
// If no requests are available sleep for RequestWaitTime
// If a signal is received on the status channel,
// move to the next phase by calling runMultiPaxos
// If the replica is the new leader, then reset the phaseOneDone,
// acceptReqQueue, clientRequestQueue, nextOperation and call runMultiPaxos
func (p *Proposer) Start() {
	go func() {
		// DONE(student): Complete the function
		// BM: DH is unsure about this method
		for {
			select {
			case newLeader := <-p.ldChan:
				p.leader = newLeader
				if p.isLeader() {
					p.phaseOneDone = false
					p.acceptMsgQueue = list.New()
					p.clientRequestQueue = list.New()
					p.nextOperation = Prepare
				}
			case <-p.stop:
				return
			case status := <-p.statusChannel:
				if status {
					time.Sleep(RequestWaitTime)
				} else {
					p.runMultiPaxos()
				}
			default:
				p.runMultiPaxos()
			}
		}
	}()
}

// Stop stops the proposer's run loop.
func (p *Proposer) Stop() {
	close(p.stop)
}

// IncrementAllDecidedUpTo increments the Proposer's adu variable by one.
func (p *Proposer) IncrementAllDecidedUpTo() {
	p.Lock()
	defer p.Unlock()
	p.aduSlotID++
}

// increaseCrnd increases the proposer's current round (crnd field)
// with the size of the Paxos configuration.
func (p *Proposer) increaseCrnd() {
	p.Lock()
	defer p.Unlock()
	p.crndID = p.crndID + int32(len(p.nodeMap))
}

// Perform accept on the servers and return error if required
//  1. Check if any pending accept requests in the acceptReqQueue to process
//  2. Check if any pending client requests in the clientRequestQueue to process
//  3. Increment the nextSlotID and prepare an accept message with the pending request,
//     crndID and nextSlotID.
//  4. Perform accept quorum call on the configuration and set the learnMsg with the
//     return value of the quorum call.
func (p *Proposer) performAccept() error {
	request := popFront(p.acceptMsgQueue)
	if request == nil {
		request = popFront(p.clientRequestQueue)
	}
	if request == nil {
		return fmt.Errorf("nothing in queues")
	}
	p.nextSlotID++
	acceptMsg := pb.AcceptMsg{
		Slot: &pb.Slot{Id: p.nextSlotID},
		Rnd:  &pb.Round{Id: p.crndID},
		Val:  request,
	}
	//4
	result, err := p.config.Accept(context.Background(), &acceptMsg)
	if err != nil {
		return err
	}
	p.learnMsg = result
	return nil
}

func (p *Proposer) performCommit() error {
	_, err := p.config.Commit(context.Background(), p.learnMsg)
	return err
}

func (p *Proposer) isPhaseOneDone() bool {
	p.RLock()
	defer p.RUnlock()
	return p.phaseOneDone
}

// ProcessRequest: processes the request from the client by putting it into the
// clientRequestQueue and calling the getResponse to get the response matching
// the request.
func (p *Proposer) AddRequestToQ(request *pb.Value) {
	if p.isLeader() {
		//log.Printf("Node id %d\t is adding the request to Queue %v", p.id, request)
		p.Lock()
		p.clientRequestQueue.PushBack(request)
		p.Unlock()
	}
}

// *list.Element
func popFront(l *list.List) *pb.Value {
	value := l.Front()
	if value == nil {
		return nil
	}
	l.Remove(value)
	val, ok := value.Value.(*pb.Value)
	if !ok {
		return nil
	}
	return val
}
*/
