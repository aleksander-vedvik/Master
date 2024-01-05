//go:build !solution
// +build !solution

package gorumspaxos

import (
	"container/list"
	"fmt"
	"log"
	"net"
	"os"
	"paxos/leaderdetector"
	"sync"
	"time"

	fd "paxos/gorumsfd"
	pb "paxos/proto"

	"github.com/relab/gorums"
)

// constants used in implementation
const (
	// NoRoundId indicates invalid value for round
	NoRoundId int32 = -1
	// Ignore indicates invalid slotID
	Ignore uint32 = 0
	// waitTimeForPhaseOne, is the time which the phase1 quorum functions wait for reply
	waitTimeForPhaseOne time.Duration = 5 * time.Second
	// waitTimeForPhaseTwo, is the time which the phase2 quorum functions wait for reply
	waitTimeForPhaseTwo time.Duration = 5 * time.Second
	// failureDetectordDelay is the default value of the delay parameter in failure detector
	failureDetectordDelay time.Duration = 5 * time.Second
	// failureDetectordDelta is the default value of the delta parameter in failure detector
	failureDetectordDelta time.Duration = 1 * time.Second
	// BM: Time to wait before throughput dump
	waitBeforeDump int = 10
)

// PaxosReplica is the structure composing the Proposer and Acceptor.
// failureDetector implementation of the failure detector
// DO NOT remove the existing fields in the structure
type PaxosReplica struct {
	pb.MultiPaxos
	sync.Mutex
	*Acceptor
	*Proposer
	failureDetector fd.FailureDetector
	id              int        // id is the id of the node
	localAddr       string     // localAddr is the local address of the replica
	responseList    *list.List // responseChan is the channel used by replica to deliver the response to proposer

	reqCounter     int
	idleCounter    int // amount of times its 0
	throughputList []int
}

// NewPaxosReplica returns a new Paxos replica with a configuration as provided
// by the input addrs. This replica will run on the given port.
func NewPaxosReplica(args NewPaxosReplicaArgs) *PaxosReplica {
	acceptor := NewAcceptor()
	quorumSize := (len(args.NodeMap)-1)/2 + 1
	qspec := NewPaxosQSpec(quorumSize)
	nodeIds := make([]int, 0)
	for _, id := range args.NodeMap {
		nodeIds = append(nodeIds, int(id))
	}
	ld := leaderdetector.NewMonLeaderDetector(nodeIds)
	proposerArgs := NewProposerArgs{
		id:               args.Id,
		aduSlotID:        Ignore,
		leaderDetector:   ld,
		qspec:            qspec,
		nodeMap:          args.NodeMap,
		phaseOneWaitTime: waitTimeForPhaseOne,
		phaseTwoWaitTime: waitTimeForPhaseTwo,
	}
	proposer := NewProposer(proposerArgs)
	failureDetector := fd.NewEvtFailureDetector(args.Id, ld, args.NodeMap,
		failureDetectordDelay, failureDetectordDelta)
	responseList := list.New()
	return &PaxosReplica{
		localAddr:       args.LocalAddr,
		Acceptor:        acceptor,
		Proposer:        proposer,
		failureDetector: failureDetector,
		id:              args.Id,
		responseList:    responseList,
		reqCounter:      0,
		idleCounter:     0,
		throughputList:  make([]int, 0),
	}
}

// ServerStart starts the replica
// 1. Invokes the start function of the proposer
// 2. Create a new gorums server
// 3. Register MultiPaxos server
// 4. Start failure detector
// 5. Call Serve on gorums server
func (replica *PaxosReplica) ServerStart(lis net.Listener) {
	go func() {
		for {
			<-time.After(1 * time.Second)
			if replica.reqCounter != 0 {
				replica.throughputList = append(replica.throughputList, replica.reqCounter)
				replica.idleCounter = 0
			}
			if replica.reqCounter == 0 {
				replica.idleCounter += 1
			}
			if replica.idleCounter > waitBeforeDump {
				log.Printf("Throughput: %d ", replica.throughputList)
				replica.dumpThroughput()
				break
			}

			replica.reqCounter = 0
		}
	}()
	// DONE(student) Implement the function
	//1
	replica.Proposer.Start()
	//2
	gorumsServer := gorums.NewServer()
	//3
	pb.RegisterMultiPaxosServer(gorumsServer, replica)
	// 4
	replica.failureDetector.StartFailureDetector(gorumsServer)
	//5
	gorumsServer.Serve(lis)
}

// Prepare handles the prepare quorum calls from the proposer by passing the received messages to its acceptor.
// It receives prepare massages and pass them to handlePrepare method of acceptor.
// It returns promise messages back to the proposer by its acceptor.
func (r *PaxosReplica) Prepare(ctx gorums.ServerCtx, prepMsg *pb.PrepareMsg) (*pb.PromiseMsg, error) {
	//log.Printf("Node id %d \t Acceptor: Prepare(%v) received", r.id, prepMsg)
	prm := r.handlePrepare(prepMsg)
	return prm, nil
}

// Accept handles the accept quorum calls from the proposer by passing the received messages to its acceptor.
// It receives Accept massages and pass them to handleAccept method of acceptor.
// It returns learn massages back to the proposer by its acceptor
func (r *PaxosReplica) Accept(ctx gorums.ServerCtx, accMsg *pb.AcceptMsg) (*pb.LearnMsg, error) {
	//log.Printf("Node id %d \t Acceptor: Accept(%v) received", r.id, accMsg)
	lrn := r.handleAccept(accMsg)
	return lrn, nil
}

// Commit is invoked when the proposer calls the commit RPC on the configuration.
// It receives a learn massage from proposer, this means the request is decided and
// the replica can commit the request. It returns an empty massage back.
func (r *PaxosReplica) Commit(ctx gorums.ServerCtx, lrnMsg *pb.LearnMsg) (*pb.Empty, error) {
	//log.Printf("Node id %d \t Learner: Commit(%v) received", r.id, lrnMsg)
	r.IncrementAllDecidedUpTo()
	r.Lock()
	r.responseList.PushBack(lrnMsg.Val)
	r.Unlock()
	return &pb.Empty{}, nil
}

// ClientHandle method is invoked when the client calls the ClientHandle RPC on all the replicas.
// This method may be called multiple times from different clients.
// If the client request is committed on the replica then it should send the response.
// P.S. Since the method is called by multiple clients, do remember to return the matching reply to
// the client. Example If Client C send the replica the request M1 then return to the client when M1 is decided.
// While waiting for M1 to get committed, M2 may be proposed and committed by the replicas.
// getResponse method helps you to match the request to the response.
func (r *PaxosReplica) ClientHandle(ctx gorums.ServerCtx, req *pb.Value) (rsp *pb.Response, err error) {
	//log.Printf("Node id %d\t Replica: ClientHandle(%v) received", r.id, req)
	r.reqCounter++
	r.AddRequestToQ(req)
	return r.getResponse(req)
}

// getResponse: is called after adding the client request to the queue.
// This function waits on the responseList for the response
// matching the request. This involves creating a new go routine and making it
// repeatedly check the responseList for the matched response. If a response
// is not present within the phaseTwoWaitTime, clean up the created goroutine
// and return error.
// P.S. Channels can also be used to achieve this functionality, if required you have complete
// freedom to change the definition of this method.
func (r *PaxosReplica) getResponse(request *pb.Value) (rsp *pb.Response, err error) {
	// DONE(student) Complete the method
	// Buffering the channels to prevent deadlock
	responseChan := make(chan *pb.Response, 1)
	stopChan := make(chan struct{}, 1)
	// This go routine periodically checks the reponseList for any responses.
	// If there exists a response corresponding to the request, it will be sent on the responseChan.
	go func() {
		for {
			select {
			// The goroutine did not find any response matching the request in the given time (waitTimeForPhaseTwo)
			case <-stopChan:
				return
			// Check if there are any responses matching the request
			default:
				//time.Sleep(10 * time.Millisecond)
				r.Lock()
				val := checkResponseList(r.responseList, request)
				// Unlocking and exiting the goroutine if a response is found
				if val != nil {
					responseChan <- val
					r.Unlock()
					return
				}
				r.Unlock()
			}
		}
	}()
	// Uses timer to terminate the goroutine and a channel to get the appropiate response
	select {
	case rsp = <-responseChan:
		return rsp, nil
	case <-time.After(waitTimeForPhaseTwo):
		stopChan <- struct{}{}
		return &pb.Response{}, fmt.Errorf("no response")
	}
}

func (r *PaxosReplica) dumpThroughput() {
	f, err := os.Create(fmt.Sprintf("server=%v.csv", r.id))
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err = f.WriteString("Throughput (s)\n")

	if err != nil {
		log.Fatal(err)
	}

	for _, data := range r.throughputList {
		_, err = f.WriteString(fmt.Sprintf("%v\n", data))

		if err != nil {
			log.Fatal(err)
		}
	}
}

// Goes through the responseList to see if there are any responses matching the request.
// Not sure if it is correct to convert to *pb.Response..
// DONE: Check if this works as expected. It should return a response corresponding to the given request.
func checkResponseList(responseList *list.List, request *pb.Value) *pb.Response {
	for response := responseList.Front(); response != nil; response = response.Next() {
		val, ok := response.Value.(*pb.Value)
		if !ok {
			continue
		}
		if checkReqAndResIsEqual(request, val) {
			resp := &pb.Response{
				ClientID:      val.ClientID,
				ClientSeq:     val.ClientSeq,
				ClientCommand: val.ClientCommand,
			}
			return resp
		}
	}
	return nil
}

func checkReqAndResIsEqual(req *pb.Value, rsp *pb.Value) bool {
	if req == nil && rsp == nil {
		return true
	}
	if req == nil || rsp == nil {
		return false
	}
	if req.ClientCommand != rsp.ClientCommand {
		return false
	}
	if req.ClientID != rsp.ClientID {
		return false
	}
	if req.ClientSeq != rsp.ClientSeq {
		return false
	}
	return true
}
