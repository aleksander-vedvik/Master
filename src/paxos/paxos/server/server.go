package server

import (
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
	message     *pb.Value
}

type PaxosServer struct {
	*pb.Server
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
}

func NewPaxosServer(addr string, srvAddresses []string) *PaxosServer {
	srv := PaxosServer{
		Server:     pb.NewServer(),
		data:       make([]string, 0),
		addr:       addr,
		peers:      srvAddresses,
		addedMsgs:  make(map[string]bool),
		clientReqs: make([]*clientReq, 0),
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
	view, err := srv.mgr.NewConfiguration(gorums.WithNodeList(srv.peers), newQSpec(len(srv.peers)))
	if err != nil {
		panic(err)
	}
	srv.SetView(view)
}

func (srv *PaxosServer) Start() {
	lis, err := net.Listen("tcp4", srv.addr)
	if err != nil {
		log.Fatal(err)
	}
	//go s.status()
	go srv.Serve(lis)
	srv.addr = lis.Addr().String()
	slog.Info(fmt.Sprintf("Server started. Listening on address: %s\n\t- peers: %v\n", srv.addr, srv.peers))
	srv.leaderElection = leaderelection.New(srv.View)
	srv.leaderElection.StartLeaderElection()
	go srv.listenForLeaderChanges()
}

func (srv *PaxosServer) listenForLeaderChanges() {
	for leader := range srv.leaderElection.Leaders() {
		slog.Warn("leader changed", "leader", leader)
		srv.leader = leader
		if srv.isLeader() {
			srv.runPhaseOne()
		}
	}
}

func (srv *PaxosServer) Write(ctx gorums.ServerCtx, request *pb.Value, broadcast *pb.Broadcast) {
	if !srv.isLeader() {
		// alternatives:
		// 1. simply ignore request 			<- ok
		// 2. send it to the leader 			<- ok
		// 3. reply with last committed value	<- ok
		// 3. reply with error					<- not ok
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

func (srv *PaxosServer) dispatch() {
	srv.mu.Lock()
	for _, req := range srv.clientReqs {
		srv.BroadcastAccept(&pb.AcceptMsg{
			Rnd:  srv.rnd,
			Slot: srv.maxSeenSlot,
			Val:  req.message,
		}, req.broadcastID)
	}
	srv.clientReqs = make([]*clientReq, 0)
	srv.mu.Unlock()
}

func (srv *PaxosServer) isLeader() bool {
	return srv.leader == srv.addr
}
