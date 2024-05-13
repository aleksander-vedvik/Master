package server

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net"
	"sync"
	"time"

	ld "github.com/aleksander-vedvik/benchmark/leaderelection"
	pb "github.com/aleksander-vedvik/benchmark/paxosm/proto"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	*pb.Server
	acceptor              *Acceptor
	id                    uint32
	leaderElection        *ld.MonLeader[*pb.Node, *pb.Configuration, *pb.Heartbeat]
	leader                string
	data                  []string
	addr                  string
	peers                 []string
	mgr                   *pb.Manager
	mut                   sync.Mutex
	proposerCtx           context.Context
	cancelProposer        context.CancelFunc
	proposer              *Proposer
	disableLeaderElection bool
	timeout               time.Duration
}

func New(addr string, srvAddrs []string, disableLeaderElection ...bool) *Server {
	disable := false
	if len(disableLeaderElection) > 0 {
		disable = disableLeaderElection[0]
	}
	id := 0
	for i, srvAddr := range srvAddrs {
		if addr == srvAddr {
			id = i
			break
		}
	}
	srv := Server{
		id:                    uint32(id),
		Server:                pb.NewServer(gorums.WithMetrics()),
		data:                  make([]string, 0),
		addr:                  addr,
		peers:                 srvAddrs,
		leader:                "",
		disableLeaderElection: disable,
		timeout:               5 * time.Second,
	}
	srv.configureView()
	srv.acceptor = NewAcceptor(len(srvAddrs), srv.View)
	pb.RegisterPaxosMServer(srv.Server, &srv)
	return &srv
}

func (srv *Server) configureView() {
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

func (srv *Server) Stop() {
	if srv.proposer != nil {
		srv.proposer.stopFunc()
	}
	srv.Server.Stop()
	srv.mgr.Close()
}

func (srv *Server) Start() {
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
	if !srv.disableLeaderElection {
		// start leader election and failure detector
		srv.leaderElection = ld.New(srv.View, id, func(id uint32) *pb.Heartbeat {
			return &pb.Heartbeat{
				Id: id,
			}
		})
		srv.leaderElection.StartLeaderElection()
		srv.proposerCtx, srv.cancelProposer = context.WithCancel(context.Background())
		go srv.listenForLeaderChanges()
	} else {
		srv.leader = srv.peers[0]
		if srv.leader == srv.addr {
			//srv.proposer = NewProposer(srv.id, srv.peers, srv.acceptor.rnd, srv.View, srv.BroadcastAccept)
			go srv.proposer.Start()
		}
	}
	// start gRPC server
	go srv.Serve(lis)
}

func (srv *Server) listenForLeaderChanges() {
	for leader := range srv.leaderElection.Leaders() {
		slog.Warn("new leader", "leader", leader)
		if leader == "" {
			leader = srv.addr
		}
		srv.mut.Lock()
		srv.leader = leader
		if srv.proposer != nil {
			srv.proposer.Stop()
			srv.proposer = nil
		}
		srv.mut.Unlock()
		if srv.isLeader() {
			//srv.proposer = NewProposer(srv.id, srv.peers, srv.acceptor.rnd, srv.View, srv.BroadcastAccept)
			go srv.proposer.Start()
		}
	}
}

func (srv *Server) Write(ctx gorums.ServerCtx, request *pb.PaxosValue) (*pb.PaxosResponse, error) {
	if !srv.isLeader() {
		// alternatives:
		// 1. simply ignore request 			<- ok
		// 2. send it to the leader 			<- ok
		// 3. reply with last committed value	<- ok
		// 4. reply with error					<- not ok
		//broadcast.Forward(request, srv.leader)
		srv.acceptor.respChans[request.MsgID] = make(chan *pb.PaxosResponse)
		defer delete(srv.acceptor.respChans, request.MsgID)
		select {
		case resp := <-srv.acceptor.respChans[request.MsgID]:
			return resp, nil
		case <-time.After(srv.timeout):
			return nil, errors.New("timed out")
		}
	}
	srv.proposer.mut.Lock()
	rnd := srv.proposer.rnd
	adu := srv.proposer.adu
	srv.proposer.adu++
	srv.proposer.mut.Unlock()

	srv.View.Accept(context.Background(), &pb.AcceptMsg{
		Rnd:  rnd,
		Slot: adu,
		Val:  request,
	}, gorums.WithNoSendWaiting())

	//broadcast.Accept(&pb.AcceptMsg{
	//	Rnd:  rnd,
	//	Slot: adu,
	//	Val:  request,
	//})
	srv.acceptor.respChans[request.MsgID] = make(chan *pb.PaxosResponse)
	defer delete(srv.acceptor.respChans, request.MsgID)
	select {
	case resp := <-srv.acceptor.respChans[request.MsgID]:
		return resp, nil
	case <-time.After(srv.timeout):
		return nil, errors.New("timed out")
	}
}

func (srv *Server) isLeader() bool {
	srv.mut.Lock()
	defer srv.mut.Unlock()
	return srv.leader == srv.addr
}

func (srv *Server) Ping(ctx gorums.ServerCtx, request *pb.Heartbeat) {
	srv.leaderElection.Ping(request.GetId())
}

func (srv *Server) Benchmark(ctx gorums.ServerCtx, request *pb.Empty) (*pb.Result, error) {
	metrics := srv.GetStats()
	m := []*pb.Metric{
		{
			TotalNum:              metrics.TotalNum,
			FinishedReqsTotal:     metrics.FinishedReqs.Total,
			FinishedReqsSuccesful: metrics.FinishedReqs.Succesful,
			FinishedReqsFailed:    metrics.FinishedReqs.Failed,
			Dropped:               metrics.Dropped,
			RoundTripLatency: &pb.TimingMetric{
				Avg: uint64(metrics.RoundTripLatency.Avg),
				Min: uint64(metrics.RoundTripLatency.Min),
				Max: uint64(metrics.RoundTripLatency.Max),
			},
		},
	}
	return &pb.Result{
		Metrics: m,
	}, nil
}
