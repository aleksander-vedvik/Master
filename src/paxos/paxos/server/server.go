package server

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"paxos/leaderelection"
	"sync"

	pb "paxos/proto"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type clientReq struct {
	broadcastID string
	message     *pb.PaxosValue
}

type PaxosServer struct {
	*pb.Server
	id             uint32
	leaderElection *leaderelection.MonLeader
	leader         string
	data           []string
	addr           string
	peers          []string
	addedMsgs      map[string]bool
	clientReqs     []*clientReq
	mgr            *pb.Manager
	rnd            uint32 // current round
	maxSeenSlot    uint32
	slots          map[uint32]*pb.PromiseSlot // slots: is the internal data structure maintained by the acceptor to remember the slots
	mu             sync.Mutex
	proposerCtx    context.Context
	cancelProposer context.CancelFunc
	proposerMutex  sync.Mutex
	proposer       *Proposer
}

func NewPaxosServer(id int, srvAddresses map[int]string) *PaxosServer {
	srvAddrs := make([]string, 0, len(srvAddresses))
	for _, addr := range srvAddresses {
		srvAddrs = append(srvAddrs, addr)
	}
	srv := PaxosServer{
		id:         uint32(id),
		Server:     pb.NewServer(),
		data:       make([]string, 0),
		addr:       srvAddresses[id],
		peers:      srvAddrs,
		addedMsgs:  make(map[string]bool),
		clientReqs: make([]*clientReq, 0),
		slots:      make(map[uint32]*pb.PromiseSlot),
		leader:     "",
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
	slog.Info(fmt.Sprintf("Server started. Listening on address: %s\n\t- peers: %v\n", srv.addr, srv.peers))
	// start leader election and failure detector
	srv.leaderElection = leaderelection.New(srv.View, id)
	srv.leaderElection.StartLeaderElection()
	srv.proposerCtx, srv.cancelProposer = context.WithCancel(context.Background())
	go srv.listenForLeaderChanges()
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
		}
		srv.mu.Unlock()
		if srv.isLeader() {
			srv.proposer = NewProposer(srv.id, srv.peers, srv.rnd, srv.View, srv.BroadcastAccept)
			go srv.proposer.Start()
		}
	}
}

func (srv *PaxosServer) Write(ctx gorums.ServerCtx, request *pb.PaxosValue, broadcast *pb.Broadcast) {
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
}

func (srv *PaxosServer) isLeader() bool {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	return srv.leader == srv.addr
}
