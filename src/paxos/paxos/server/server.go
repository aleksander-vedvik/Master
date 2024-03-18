package server

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"paxos/leaderelection"
	"time"

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
	viewNumber     int32
	clientReqs     []*clientReq
	//requestQueue   []*pb.PrePrepareRequest
	//maxLimitOfReqs int
	sequenceNumber int32
	mgr            *pb.Manager
	rnd            uint32
	maxSeenSlot    uint32
	slots          map[uint32]*pb.PromiseSlot // slots: is the internal data structure maintained by the acceptor to remember the slots
}

func NewPaxosServer(addr string, srvAddresses []string) *PaxosServer {
	srv := PaxosServer{
		Server:         pb.NewServer(),
		data:           make([]string, 0),
		addr:           addr,
		peers:          srvAddresses,
		addedMsgs:      make(map[string]bool),
		clientReqs:     make([]*clientReq, 0),
		leader:         "",
		sequenceNumber: 1,
		viewNumber:     1,
	}
	srv.configureView()
	pb.RegisterMultiPaxosServer(srv.Server, &srv)
	return &srv
}

func (srv *PaxosServer) configureView() {
	srv.mgr = pb.NewManager(
		gorums.WithDialTimeout(50*time.Millisecond),
		gorums.WithGrpcDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	view, err := srv.mgr.NewConfiguration(gorums.WithNodeList(srv.peers))
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
	}
}

func (srv *PaxosServer) Write(ctx gorums.ServerCtx, request *pb.Value, broadcast *pb.Broadcast) {
	md := broadcast.GetMetadata()
	srv.clientReqs = append(srv.clientReqs, &clientReq{
		broadcastID: md.BroadcastID,
		message:     request,
	})
	broadcast.Accept(&pb.AcceptMsg{
		Rnd:  srv.rnd,
		Slot: srv.maxSeenSlot,
		Val:  request,
	})
}

func (srv *PaxosServer) sendWrong() {
	for _, req := range srv.clientReqs {
		srv.View.Accept(context.Background(), &pb.AcceptMsg{
			Rnd:  srv.rnd,
			Slot: srv.maxSeenSlot,
			Val:  req.message,
		})
	}
}

func (srv *PaxosServer) sendCorrect() {
	for _, req := range srv.clientReqs {
		srv.BroadcastAccept(&pb.AcceptMsg{
			Rnd:  srv.rnd,
			Slot: srv.maxSeenSlot,
			Val:  req.message,
		}, req.broadcastID)
	}
}
