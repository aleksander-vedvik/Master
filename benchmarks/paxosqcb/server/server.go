package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	pb "github.com/aleksander-vedvik/benchmark/paxosqcb/proto"

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
	responseTimeout = 120 * time.Second
	// managerDialTimeout is the default timeout for dialing a manager
	managerDialTimeout = 5 * time.Second
)

type resp struct {
	respChan chan *pb.Response
	ctx      context.Context
}

type learnMsg struct {
	send func()
	msg  *pb.LearnMsg
}

// PaxosReplica is the structure composing the Proposer and Acceptor.
type PaxosReplica struct {
	*pb.Server
	mu sync.Mutex
	*Acceptor
	*Proposer
	paxosManager     *pb.Manager // gorums paxos manager (from generated code)
	id               int         // id is the id of the node
	addr             string
	stop             chan struct{}      // channel for stopping the replica's run loop.
	learntVal        map[Slot]*learnMsg // Stores all received learn messages
	responseChannels map[string]*resp
	stopped          bool
}

// NewPaxosReplica returns a new Paxos replica with a nodeMap configuration.
func New(addr string, srvAddrs []string, logger *slog.Logger) *PaxosReplica {
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
		gorums.WithLogger(logger),
	}
	r := &PaxosReplica{
		Server:           pb.NewServer(gorums.WithSLogger(logger)),
		Acceptor:         NewAcceptor(),
		Proposer:         NewProposer(myID, 0, nodeMap),
		paxosManager:     pb.NewManager(opts...),
		id:               myID,
		addr:             addr,
		stop:             make(chan struct{}),
		learntVal:        make(map[Slot]*learnMsg, 100),
		responseChannels: make(map[string]*resp),
	}
	pb.RegisterPaxosQCBServer(r.Server, r)
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
	close(r.stop) // stop the replica's run loop
	r.paxosManager.Close()
	r.Server.Stop()
	r.stopped = true
}

func (r *PaxosReplica) Start(local bool) {
	lis, err := net.Listen("tcp", r.addr)
	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf("Server started. Listening on address: %s\n", r.addr))
	if local {
		go r.Serve(lis)
		return
	}
	r.Serve(lis)
}

// run starts the replica's run loop.
// It subscribes to the leader detector's trust messages and signals the proposer when a new leader is detected.
// It also starts the failure detector, which is necessary to get leader detections.
func (r *PaxosReplica) run() {
	qspec := NewPaxosQSpec(len(r.nodeMap))
	paxConfig, err := r.paxosManager.NewConfiguration(qspec, gorums.WithNodeMap(r.nodeMap))
	if err != nil {
		return
	}
	r.Proposer.setConfiguration(paxConfig)
	r.SetView(paxConfig)
	go func() {
		for {
			select {
			case <-r.stop:
				return
			case msg := <-r.Proposer.msgQueue:
				if !r.isLeader() {
					continue
				}
				r.mu.Lock()
				r.Proposer.nextSlot++
				msg.accept.Rnd = r.Proposer.crnd
				msg.accept.Slot = r.Proposer.nextSlot
				r.mu.Unlock()
				lrn, err := r.Proposer.performAccept(msg.accept)
				if err != nil {
					continue
				}
				select {
				case <-r.stop:
					return
				default:
				}
				r.Proposer.performCommit(lrn, msg.broadcast)
				/*default:
				if r.isLeader() {
					r.runMultiPaxos()
				}*/
			}
		}
	}()
}

// Prepare handles the prepare quorum calls from the proposer by passing the received messages to its acceptor.
// It receives prepare massages and pass them to handlePrepare method of acceptor.
// It returns promise messages back to the proposer by its acceptor.
func (r *PaxosReplica) Prepare(ctx gorums.ServerCtx, prepMsg *pb.PrepareMsg) (*pb.PromiseMsg, error) {
	//ctx.Release()
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.handlePrepare(prepMsg), nil
}

// Accept handles the accept quorum calls from the proposer by passing the received messages to its acceptor.
// It receives Accept massages and pass them to handleAccept method of acceptor.
// It returns learn massages back to the proposer by its acceptor
func (r *PaxosReplica) Accept(ctx gorums.ServerCtx, accMsg *pb.AcceptMsg) (*pb.LearnMsg, error) {
	//ctx.Release()
	r.mu.Lock()
	defer r.mu.Unlock()
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
func (r *PaxosReplica) Commit(ctx gorums.ServerCtx, learn *pb.LearnMsg, broadcast *pb.Broadcast) {
	//slog.Info("received commit", "replica", r.addr, "slot", learn.Slot)
	//ctx.Release()
	lrn := &learnMsg{
		msg: learn,
		send: func() {
			broadcast.SendToClient(&pb.Response{
				ClientID:      learn.Val.ClientID,
				ClientSeq:     learn.Val.ClientSeq,
				ClientCommand: learn.Val.ClientCommand,
			}, nil)
		},
	}
	r.mu.Lock()
	adu := r.adu + 1
	if prevLearn, ok := r.learntVal[learn.Slot]; !ok {
		r.learntVal[learn.Slot] = lrn
	} else {
		// make sure that decided values are stored with the highest round number
		if prevLearn.msg.Rnd < learn.Rnd {
			r.learntVal[learn.Slot] = lrn
		}
	}
	r.mu.Unlock()

	switch {
	case learn.Slot == adu:
		r.execute(lrn)
	case learn.Slot < adu:
	case learn.Slot > adu:
	}
}

func (r *PaxosReplica) execute(lrn *learnMsg) {
	for slot := lrn.msg.Slot; true; slot++ {
		r.mu.Lock()
		learn, ok := r.learntVal[slot]
		r.mu.Unlock()
		if !ok {
			break
		}
		r.advanceAllDecidedUpTo()
		learn.send()
	}
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
func (r *PaxosReplica) ClientHandle(ctx gorums.ServerCtx, req *pb.Value, broadcast *pb.Broadcast) {
	//ctx.Release()
	r.AddRequestToQ(req, broadcast)
}

func (r *PaxosReplica) Benchmark(ctx gorums.ServerCtx, request *pb.Empty) (*pb.Empty, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	slog.Info("purging state")
	close(r.stop)
	r.Acceptor = NewAcceptor()
	r.Proposer = NewProposer(r.id, 0, r.nodeMap)
	r.learntVal = make(map[Slot]*learnMsg, 100)
	r.responseChannels = make(map[string]*resp)
	r.stop = make(chan struct{})
	r.run()
	return &pb.Empty{}, nil
}
