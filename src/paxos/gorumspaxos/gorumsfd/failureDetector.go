//go:build !solution
// +build !solution

package gorumsfd

import (
	"context"
	"log"
	"net"
	"time"

	pb "paxos/fdproto"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ManagerWaitTimeOut = 5 * time.Second
)

// EvtFailureDetector represents a Eventually Perfect Failure Detector as
// described at page 53 in:
// Christian Cachin, Rachid Guerraoui, and Luís Rodrigues: "Introduction to
// Reliable and Secure Distributed Programming" Springer, 2nd edition, 2011.

// DO NOT remove the existing fields in the structure
type EvtFailureDetector struct {
	id            int               // the id of this node
	alive         map[uint32]bool   // map of node ids considered alive
	suspected     map[uint32]bool   // map of node ids  considered suspected
	sr            SuspectRestorer   // Provided SuspectRestorer implementation
	delay         time.Duration     // the current delay for the timeout procedure
	delta         time.Duration     // the delta value to be used when increasing delay
	stop          chan struct{}     // channel for signaling a stop request to the main run loop
	nodeMap       map[string]uint32 // list of addresses of the nodes.
	manager       *pb.Manager       // Manager to create the configuration
	configuration *pb.Configuration // Configuration to call the quorum calls
}

// NewEvtFailureDetector returns a new Eventual Failure Detector. It takes the
// following arguments:
// id: The id of the node running this instance of the failure detector.
// running this instance of the failure detector).
// sr: A leader detector implementing the SuspectRestorer interface.
// addrs: A list of address of the replicas.
// delay: The timeout delay after which the failure detection is performed.
// delta: The value to be used when increasing delay.
func NewEvtFailureDetector(id int, sr SuspectRestorer, nodeMap map[string]uint32,
	delay time.Duration, delta time.Duration,
) *EvtFailureDetector {
	suspected := make(map[uint32]bool)
	alive := make(map[uint32]bool)
	for _, nodeId := range nodeMap {
		alive[nodeId] = false
	}
	mgr := pb.NewManager(gorums.WithDialTimeout(ManagerWaitTimeOut),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(), // block until connections are made
			grpc.WithTransportCredentials(insecure.NewCredentials()), // disable TLS
		),
	)
	return &EvtFailureDetector{
		id:        id,
		alive:     alive,
		suspected: suspected,
		sr:        sr,
		nodeMap:   nodeMap,
		delay:     delay,
		delta:     delta,
		manager:   mgr,
		stop:      make(chan struct{}),
	}
}

// StartFailureDetector starts main run loop in a separate goroutine.
// This function should perform the following functionalities
// 1. Register FailureDetectorServer implementation
// 2. The started Go Routine, after the e.delay, it should call PerformFailureDetection
// 3. Started Go Routine, should also wait to receive the signal to stop
func (e *EvtFailureDetector) StartFailureDetector(srv *gorums.Server) error {
	// DONE(student) complete StartFailureDetector
	err := e.registerFailureDetector(srv)
	if err != nil {
		return err
	}
	time.Sleep(e.delay) // wait until server starts
	go func() {
		for {
			select {
			case <-e.stop:
				return
			default:
				time.Sleep(e.delay)
				e.PerformFailureDetection()
			}
		}
	}()
	return nil
}

func (e *EvtFailureDetector) registerFailureDetector(gorumsSrv *gorums.Server) error {
	pb.RegisterFailureDetectorServer(gorumsSrv, newEvtFailureDetectorServer(e))
	return nil
}

// PerformFailureDetection is the method used to perform ping
// operation for all nodes and report all suspected nodes.
// 1. Create configuration with all the nodes if not previously done.
// 2. call Ping rpc on the configuration
// 3. call SendStatusOfNodes to send suspect and restore notifications.
func (e *EvtFailureDetector) PerformFailureDetection() error {
	// TODO(student) complete PerformFailureDetection
	e.runGorumsConfig()
	e.configuration.Ping(context.Background(), &pb.HeartBeat{Id: int32(e.id)})
	e.SendStatusOfNodes()
	return nil
}

func (e *EvtFailureDetector) runGorumsConfig() {
	var allNodesConfig *pb.Configuration
	var err error
	for {
		allNodesConfig, err = e.manager.NewConfiguration(
			newEvtFailureDetectorQSpec(e),
			gorums.WithNodeList(e.getAddressess()),
		)
		if err != nil {
			log.Println("Error creating read config:", err)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	e.configuration = allNodesConfig
}

func (e *EvtFailureDetector) getAddressess() []string {
	addresses := make([]string, 0, len(e.nodeMap))
	for addr := range e.nodeMap {
		if e.checkAlive(addr) {
			addresses = append(addresses, addr)
		}
	}
	return addresses
}

func (e *EvtFailureDetector) checkAlive(address string) bool {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// Stop stops e's main run loop.
func (e *EvtFailureDetector) Stop() {
	e.stop <- struct{}{}
}

// SendStatusOfNodes: reports the status of the nodes to the SuspectRestorer. This
// method is called after Ping RPC which marks all the live nodes.
// If a node which is previously suspected and now is alive then increase the e.delay by e.delta
// All non reachable nodes are reported and all previously reported,
// now live nodes are restored by calling the Suspect and Restore functions of the SuspectRestorer.
func (e *EvtFailureDetector) SendStatusOfNodes() {
	// DONE(student) complete PerformFailureDetection
	for i, alive := range e.alive {
		s, ok := e.suspected[i]
		// Suspected will be true if it exist in the map and set to true
		suspected := s && ok
		if suspected && alive {
			e.sr.Restore(int(i))
			e.suspected[i] = false
			// THIS IS WRONG COMPARED TO THE IMPLEMENTATION IN THE BOOK....
			// THE DELAY SHOULD ONLY BE INCREASED ONCE PER TIMEOUT AND NOT ONCE PER NODE.
			e.delay += e.delta
		}
		if !alive && !suspected {
			e.sr.Suspect(int(i))
			e.suspected[i] = true
		}
		// Resetting the alive map one by one instead of in a new loop
		e.alive[i] = false
	}
}

// evtFailureDetectorServer implements the FailureDetector RPC
type evtFailureDetectorServer struct {
	e *EvtFailureDetector
}

// newEvtFailureDetectorServer returns the evtFailureDetectorServer
func newEvtFailureDetectorServer(e *EvtFailureDetector) evtFailureDetectorServer {
	return evtFailureDetectorServer{e}
}

// Ping handles the Ping RPC from the other replicas. Reply contains the id of the node
func (srv evtFailureDetectorServer) Ping(ctx gorums.ServerCtx, in *pb.HeartBeat) (resp *pb.HeartBeat, err error) {
	resp = &pb.HeartBeat{Id: int32(srv.e.id)}
	return resp, err
}

// evtFailureDetectorQSpec implements the QuorumSpec for the RPC
type evtFailureDetectorQSpec struct {
	e *EvtFailureDetector
}

// newEvtFailureDetectorQSpec returns the evtFailureDetectorQSpec
func newEvtFailureDetectorQSpec(e *EvtFailureDetector) evtFailureDetectorQSpec {
	return evtFailureDetectorQSpec{e}
}

// PingQF is the quorum function to handle the replies to Ping RPC call. Nodes replied to the call
// are marked live.
func (q evtFailureDetectorQSpec) PingQF(in *pb.HeartBeat, replies map[uint32]*pb.HeartBeat) (*pb.HeartBeat, bool) {
	// DONE(student) complete PerformFailureDetection
	for _, reply := range replies {
		q.e.alive[uint32(reply.Id)] = true
	}
	if len(replies) < len(q.e.nodeMap) {
		return nil, false
	}
	return nil, true
}
