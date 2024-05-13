package server

import (
	"errors"
	"net"
	"sync"
	"time"

	pb "github.com/aleksander-vedvik/benchmark/paxosqc/proto"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//func (srv *Server) Prepare(ctx gorums.ServerCtx, req *pb.PrepareMsg) (*pb.PromiseMsg, error) {
//return srv.acceptor.Prepare(ctx, req)
//}

//func (srv *Server) Accept(ctx gorums.ServerCtx, request *pb.AcceptMsg) {
//srv.acceptor.Accept(ctx, request)
//}

//func (srv *Server) Learn(ctx gorums.ServerCtx, request *pb.LearnMsg) {
//srv.acceptor.Learn(ctx, request)
//}

const (
	// responseTimeout is the duration to wait for a response before cancelling
	responseTimeout = 1 * time.Second
	// managerDialTimeout is the default timeout for dialing a manager
	managerDialTimeout = 5 * time.Second
)

// PaxosReplica is the structure composing the Proposer and Acceptor.
type PaxosReplica struct {
	*pb.Server
	mu sync.Mutex
	*Acceptor
	*Proposer
	paxosManager     *pb.Manager // gorums paxos manager (from generated code)
	id               int         // id is the id of the node
	addr             string
	stop             chan struct{}         // channel for stopping the replica's run loop.
	learntVal        map[Slot]*pb.LearnMsg // Stores all received learn messages
	responseChannels map[uint64]chan *pb.Response
	stopped          bool
}

// NewPaxosReplica returns a new Paxos replica with a nodeMap configuration.
func New(addr string, srvAddrs []string) *PaxosReplica {
	var myID int
	nodeMap := make(map[string]uint32)
	for id, srvAddr := range srvAddrs {
		nodeMap[srvAddr] = uint32(id)
		if addr == srvAddr {
			myID = id
		}
	}

	opts := []gorums.ManagerOption{
		gorums.WithDialTimeout(managerDialTimeout),
		gorums.WithGrpcDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()), // disable TLS
		),
	}
	r := &PaxosReplica{
		Server:           pb.NewServer(),
		Acceptor:         NewAcceptor(),
		Proposer:         NewProposer(myID, 0, nodeMap),
		paxosManager:     pb.NewManager(opts...),
		id:               myID,
		addr:             addr,
		stop:             make(chan struct{}),
		learntVal:        make(map[Slot]*pb.LearnMsg),
		responseChannels: make(map[uint64]chan *pb.Response),
	}
	pb.RegisterPaxosQCServer(r.Server, r)
	r.run()
	return r
}

// Stops the failure detector, the proposer and the gorums server
// Calling this will interrupt the processing of requests
// The failure detector will no longer respond to ping requests
func (r *PaxosReplica) Stop() {
	if r.stopped {
		return
	}
	r.stop <- struct{}{} // stop the replica's run loop
	r.paxosManager.Close()
	r.Stop()
	r.stopped = true
}

func (r *PaxosReplica) Start() {
	lis, err := net.Listen("tcp", r.addr)
	if err != nil {
		panic(err)
	}
	r.Serve(lis)
}

// run starts the replica's run loop.
// It subscribes to the leader detector's trust messages and signals the proposer when a new leader is detected.
// It also starts the failure detector, which is necessary to get leader detections.
func (r *PaxosReplica) run() {
	go func() {
		qspec := NewPaxosQSpec(len(r.nodeMap))
		paxConfig, err := r.paxosManager.NewConfiguration(qspec, gorums.WithNodeMap(r.nodeMap))
		if err != nil {
			return
		}
		r.Proposer.setConfiguration(paxConfig)

		for {
			select {
			case <-r.stop:
				return
			default:
				if r.isLeader() {
					r.runMultiPaxos()
				}
			}
		}
	}()
}

// Prepare handles the prepare quorum calls from the proposer by passing the received messages to its acceptor.
// It receives prepare massages and pass them to handlePrepare method of acceptor.
// It returns promise messages back to the proposer by its acceptor.
func (r *PaxosReplica) Prepare(ctx gorums.ServerCtx, prepMsg *pb.PrepareMsg) (*pb.PromiseMsg, error) {
	ctx.Release()
	return r.handlePrepare(prepMsg), nil
}

// Accept handles the accept quorum calls from the proposer by passing the received messages to its acceptor.
// It receives Accept massages and pass them to handleAccept method of acceptor.
// It returns learn massages back to the proposer by its acceptor
func (r *PaxosReplica) Accept(ctx gorums.ServerCtx, accMsg *pb.AcceptMsg) (*pb.LearnMsg, error) {
	ctx.Release()
	return r.handleAccept(accMsg), nil
}

// Commit is invoked by the proposer as part of the commit phase of the MultiPaxos algorithm.
// It receives a learn massage representing the proposer's decided value, meaning that the
// request can be executed by the replica. (In this lab you don't need to execute the request,
// just deliver the response to the client.)
//
// Be aware that the received learn message may not be for the next slot in the sequence.
// If the received slot is less than the next slot, the message should be ignored.
// If the received slot is greater than the next slot, the message should be buffered.
// If the received slot is equal to the next slot, the message should be delivered.
//
// This method is also responsible for communicating the decided value to the ClientHandle
// method, which is responsible for returning the response to the client.
func (r *PaxosReplica) Commit(ctx gorums.ServerCtx, learn *pb.LearnMsg) {
	ctx.Release()
	r.mu.Lock()
	adu := r.adu + 1
	if prevLearn, ok := r.learntVal[learn.Slot]; !ok {
		r.learntVal[learn.Slot] = learn
	} else {
		// make sure that decided values are stored with the highest round number
		if prevLearn.Rnd < learn.Rnd {
			r.learntVal[learn.Slot] = learn
		}
	}
	r.mu.Unlock()
	switch {
	case learn.Slot == adu:
		r.execute(learn)
	case learn.Slot < adu:
	case learn.Slot > adu:
	}
}

func (r *PaxosReplica) execute(lrn *pb.LearnMsg) {
	for slot := lrn.Slot; true; slot++ {
		r.mu.Lock()
		learn, ok := r.learntVal[slot]
		r.mu.Unlock()
		if !ok {
			break
		}
		r.advanceAllDecidedUpTo()
		respCh := r.respChannel(learn)
		if respCh == nil {
			// give up if the response channel is not found
			continue
		}
		// deliver decided value to ClientHandle
		respCh <- &pb.Response{
			ClientID:      learn.Val.ClientID,
			ClientSeq:     learn.Val.ClientSeq,
			ClientCommand: learn.Val.ClientCommand,
		}
		close(respCh)
	}
}

const (
	retryDelay = 10 * time.Millisecond
	maxRetries = 5
)

func (r *PaxosReplica) respChannel(learn *pb.LearnMsg) chan *pb.Response {
	valHash := learn.Val.Hash()
	r.mu.Lock()
	respCh := r.responseChannels[valHash]
	r.mu.Unlock()
	delay := retryDelay
	// We only enter the retry loop if the response channel is not found
	for retries := 0; respCh == nil; retries++ {
		// Can happen if the client's request is committed before the response channel is created.
		// That is, this replica has not yet received the client's request, but it has received
		// the decision via the other replicas. Wait for the response channel to be created.
		// r.Logf("Replica: Commit(%v) missing response channel for client request: hash %v", learn, valHash)
		time.Sleep(delay)
		r.mu.Lock()
		respCh = r.responseChannels[valHash]
		r.mu.Unlock()
		// exponential backoff: double the delay for each retry
		delay *= 2
		if retries > maxRetries {
			return nil
		}
	}
	return respCh
}

// ClientHandle is invoked by the client to send a request to the replicas via a quorum call and get a response.
// A response is only sent back to the client when the request has been committed by the MultiPaxos replicas.
// This method will receive requests from multiple clients and must return the response to the correct client.
// If the request is not committed within a certain time, the method may return an error.
//
// Since the method is called by multiple clients, it is essential to return the matching reply to the client.
// Consider a client that sends a request M1, once M1 has been decided, the response to M1 should be returned
// to the client. However, while waiting for M1 to get committed, M2 may be proposed and committed by the replicas.
// Thus, M2 should not be returned to the client that sent M1.
func (r *PaxosReplica) ClientHandle(ctx gorums.ServerCtx, req *pb.Value) (rsp *pb.Response, err error) {
	ctx.Release()
	r.AddRequestToQ(req)
	respChannel, cleanup := r.makeResponseChan(req)
	defer cleanup()

	select {
	case resp := <-respChannel:
		return resp, nil
	case <-time.After(responseTimeout):
		return nil, errors.New("unable to get the response")
	}
}

func (r *PaxosReplica) makeResponseChan(request *pb.Value) (chan *pb.Response, func()) {
	msgID := request.Hash()
	respChannel := make(chan *pb.Response, 1)
	r.mu.Lock()
	r.responseChannels[msgID] = respChannel
	r.mu.Unlock()
	return respChannel, func() {
		r.mu.Lock()
		delete(r.responseChannels, msgID)
		r.mu.Unlock()
	}
}

// remainingResponses returns the number of responses that are still pending.
func (r *PaxosReplica) remainingResponses() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.responseChannels)
}

// responseIDs returns the IDs of the responses that are still pending.
func (r *PaxosReplica) responseIDs() []uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	ids := make([]uint64, 0, len(r.responseChannels))
	for id := range r.responseChannels {
		ids = append(ids, id)
	}
	return ids
}
