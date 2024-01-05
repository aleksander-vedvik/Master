package failuredetector

import (
	"context"
	"distributedleader/src/defs"
	pbFd "distributedleader/src/proto"
	"log"
	"net"
	"time"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EvtFailureDetector struct {
	id            int            // the id of this node
	nodeIDs       []int          // node ids for every node in cluster
	alive         map[int]bool   // map of node ids considered alive
	suspected     map[int]bool   // map of node ids considered suspected
	nodeAddresses map[int]string // map of addresses of each node

	sr defs.SuspectRestorer // Provided SuspectRestorer implementation

	delay         time.Duration // the current delay for the timeout procedure
	delta         time.Duration // the delta value to be used when increasing delay
	timeoutSignal *time.Ticker  // the timeout procedure ticker

	mgr        *pbFd.Manager
	nodeConfig *pbFd.Configuration
}

func NewEvtFailureDetector(id int, nodeIDs []int, nodeAddresses map[int]string, sr defs.SuspectRestorer, delta time.Duration) *EvtFailureDetector {
	suspected := make(map[int]bool)
	alive := make(map[int]bool)

	for _, nodeId := range nodeIDs {
		alive[nodeId] = true
	}

	return &EvtFailureDetector{
		id:            id,
		nodeIDs:       nodeIDs,
		alive:         alive,
		suspected:     suspected,
		nodeAddresses: nodeAddresses,

		sr: sr,

		delay: delta,
		delta: delta,
	}
}

// Starts the failuredetector
func (e *EvtFailureDetector) Start() {
	startChan := make(chan bool, 0)
	go e.acceptIncomingHeartbeatRequests(startChan)
	time.Sleep(5 * time.Second) // Wait for the gorums server to start
	for {
		time.Sleep(e.delay)
		e.timeout()
	}
}

func (e *EvtFailureDetector) timeout() {
	if e.isAliveNonEmpty() && e.isSuspectedNonEmpty() {
		e.delay += e.delta
	}
	for _, nodeId := range e.nodeIDs {
		a, aOk := e.alive[nodeId]
		alive := a && aOk
		s, sOk := e.suspected[nodeId]
		suspected := s && sOk
		if !alive && !suspected {
			e.suspected[nodeId] = true
			e.sr.Suspect(nodeId)
		} else if alive && suspected {
			e.suspected[nodeId] = false
			delete(e.suspected, nodeId) // Optional
			e.sr.Restore(nodeId)
		}
	}
	e.alive = make(map[int]bool)
	e.sendHeartbeatRequests()
}

func (e *EvtFailureDetector) sendHeartbeatRequests() {
	e.runGorumsConfig()

	hbReq := pbFd.Heartbeat{
		From:    int32(e.id),
		To:      int32(-1),
		Request: true,
	}
	//log.Println("Sending heartbeat requests")
	ctx, cancel := context.WithCancel(context.Background())
	e.nodeConfig.HeartbeatRequest(ctx, &hbReq)
	cancel()
}

func (e *EvtFailureDetector) runGorumsConfig() {
	mgr := pbFd.NewManager(
		gorums.WithDialTimeout(1000*time.Millisecond),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	allNodesConfig, err := mgr.NewConfiguration(
		e,
		gorums.WithNodeList(e.getAddressess()),
	)
	if err != nil {
		log.Fatalln("Error creating read config:", err)
	}
	e.mgr = mgr
	e.nodeConfig = allNodesConfig
}

func (e *EvtFailureDetector) getAddressess() []string {
	addresses := make([]string, 0, len(e.nodeIDs))
	for _, address := range e.nodeAddresses {
		if e.checkAlive(address) {
			addresses = append(addresses, address)
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

func (e *EvtFailureDetector) HeartbeatRequestQF(_ *pbFd.Heartbeat, replies map[uint32]*pbFd.Heartbeat) (*pbFd.Heartbeat, bool) {
	for _, hb := range replies {
		e.alive[int(hb.From)] = true
	}
	if len(replies) < len(e.nodeIDs) {
		return nil, false
	}
	return nil, true
}

func (e *EvtFailureDetector) isAliveNonEmpty() bool {
	for _, alive := range e.alive {
		if alive {
			return true
		}
	}
	return false
}

func (e *EvtFailureDetector) isSuspectedNonEmpty() bool {
	for _, suspected := range e.suspected {
		if suspected {
			return true
		}
	}
	return false
}

func (e *EvtFailureDetector) acceptIncomingHeartbeatRequests(startChan chan<- bool) {
	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
	gorumsSrv := gorums.NewServer()
	pbFd.RegisterFailureDetectorServer(gorumsSrv, e)
	err = gorumsSrv.Serve(lis)
	if err != nil {
		log.Fatalln("Gorums server crashed, Error:", err)
	}
}

func (e *EvtFailureDetector) HeartbeatRequest(ctx gorums.ServerCtx, request *pbFd.Heartbeat) (response *pbFd.Heartbeat, err error) {
	return &pbFd.Heartbeat{
		From:    int32(e.id),
		To:      request.From,
		Request: false,
	}, nil
}
