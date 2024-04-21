package server

import (
	"context"
	"log"
	"log/slog"
	"net"
	"paxos/benchmark/metrics"
	"paxos/leaderelection"
	"sync"

	pb "paxos/proto"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type clientReq struct {
	broadcastID uint64
	message     *pb.PaxosValue
}

type PaxosServer struct {
	*pb.Server
	id                    uint32
	leaderElection        *leaderelection.MonLeader
	leader                string
	data                  []string
	addr                  string
	peers                 []string
	addedMsgs             map[string]bool
	clientReqs            []*clientReq
	mgr                   *pb.Manager
	rnd                   uint32 // current round
	maxSeenSlot           uint32
	slots                 map[uint32]*pb.PromiseSlot // slots: is the internal data structure maintained by the acceptor to remember the slots
	mu                    sync.Mutex
	proposerCtx           context.Context
	cancelProposer        context.CancelFunc
	proposer              *Proposer
	disableLeaderElection bool
	senders               map[uint64]int
	numMsgs               map[int]int
	metrics               *metrics.Metrics
}

func NewPaxosServer(id int, srvAddresses map[int]string, disableLeaderElection ...bool) *PaxosServer {
	disable := false
	if len(disableLeaderElection) > 0 {
		disable = disableLeaderElection[0]
	}
	srvAddrs := make([]string, len(srvAddresses))
	for id, addr := range srvAddresses {
		srvAddrs[id] = addr
	}
	srv := PaxosServer{
		id:                    uint32(id),
		Server:                pb.NewServer(gorums.WithMetrics()),
		data:                  make([]string, 0),
		addr:                  srvAddresses[id],
		peers:                 srvAddrs,
		addedMsgs:             make(map[string]bool),
		clientReqs:            make([]*clientReq, 0),
		slots:                 make(map[uint32]*pb.PromiseSlot),
		leader:                "",
		disableLeaderElection: disable,
		senders:               make(map[uint64]int),
		numMsgs:               make(map[int]int),
	}
	srv.configureView()
	pb.RegisterMultiPaxosServer(srv.Server, &srv)
	return &srv
}

func (srv *PaxosServer) configureView() {
	srv.mgr = pb.NewManager(
		gorums.WithGrpcDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	view, err := srv.mgr.NewConfiguration(gorums.WithNodeList(srv.peers), newQSpec(1+len(srv.peers)/2))
	if err != nil {
		panic(err)
	}
	srv.SetView(view)
}

func (srv *PaxosServer) Stop() {
	if srv.proposer != nil {
		srv.proposer.stopFunc()
	}
	srv.Server.Stop()
	srv.mgr.Close()
}

func (srv *PaxosServer) Start() {
	// create listener
	lis, err := net.Listen("tcp4", srv.addr)
	if err != nil {
		log.Fatal(err)
	}
	// add the address
	srv.addr = lis.Addr().String()
	// add the correct ID to the server
	var id uint32
	for _, node := range srv.View.Nodes() {
		if node.Address() == srv.addr {
			id = node.ID()
			break
		}
	}
	//slog.Info(fmt.Sprintf("Server started. Listening on address: %s\n\t- peers: %v\n", srv.addr, srv.peers))
	if !srv.disableLeaderElection {
		// start leader election and failure detector
		srv.leaderElection = leaderelection.New(srv.View, id)
		srv.leaderElection.StartLeaderElection()
		srv.proposerCtx, srv.cancelProposer = context.WithCancel(context.Background())
		go srv.listenForLeaderChanges()
	} else {
		srv.leader = srv.peers[0]
		if srv.leader == srv.addr {
			srv.proposer = NewProposer(srv.id, srv.peers, srv.rnd, srv.View, srv.BroadcastAccept)
			go srv.proposer.Start()
		}
	}
	// start gRPC server
	go srv.Serve(lis)
}

func (srv *PaxosServer) listenForLeaderChanges() {
	for leader := range srv.leaderElection.Leaders() {
		slog.Warn("new leader", "leader", leader)
		if leader == "" {
			leader = srv.addr
		}
		srv.mu.Lock()
		srv.leader = leader
		if srv.proposer != nil {
			srv.proposer.Stop()
			srv.proposer = nil
		}
		srv.mu.Unlock()
		if srv.isLeader() {
			srv.proposer = NewProposer(srv.id, srv.peers, srv.rnd, srv.View, srv.BroadcastAccept)
			go srv.proposer.Start()
		}
	}
}

/*func (srv *PaxosServer) Write(ctx gorums.ServerCtx, request *pb.PaxosValue, broadcast *pb.Broadcast) {
	if !srv.isLeader() {
		// alternatives:
		// 1. simply ignore request 			<- ok
		// 2. send it to the leader 			<- ok
		// 3. reply with last committed value	<- ok
		// 4. reply with error					<- not ok
		return
	}
	md := broadcast.GetMetadata()
	srv.mu.Lock()
	srv.clientReqs = append(srv.clientReqs, &clientReq{
		broadcastID: md.BroadcastID,
		message:     request,
	})
	srv.mu.Unlock()
}*/

func (srv *PaxosServer) Write(ctx gorums.ServerCtx, request *pb.PaxosValue, broadcast *pb.Broadcast) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	srv.numMsgs[10]++
	if srv.leader != srv.addr {
		// alternatives:
		// 1. simply ignore request 			<- ok
		// 2. send it to the leader 			<- ok
		// 3. reply with last committed value	<- ok
		// 4. reply with error					<- not ok
		broadcast.Forward(request, srv.leader)
		return
	}
	srv.proposer.maxSeenSlot++
	broadcast.Accept(&pb.AcceptMsg{
		Rnd:  srv.proposer.rnd,
		Slot: srv.proposer.maxSeenSlot,
		Val:  request,
	})
}

func (srv *PaxosServer) isLeader() bool {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	return srv.leader == srv.addr
}

func (srv *PaxosServer) Status() map[int]int {
	if srv == nil {
		return nil
	}
	return srv.numMsgs
}
