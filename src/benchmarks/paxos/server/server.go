package server

import (
	"context"
	"log"
	"log/slog"
	"net"
	"sync"

	ld "github.com/aleksander-vedvik/benchmark/leaderelection"
	pb "github.com/aleksander-vedvik/benchmark/paxos/proto"

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
	senders               map[uint64]int
	//metrics               *metrics.Metrics
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
		acceptor:              NewAcceptor(len(srvAddrs)),
		data:                  make([]string, 0),
		addr:                  addr,
		peers:                 srvAddrs,
		leader:                "",
		disableLeaderElection: disable,
		senders:               make(map[uint64]int),
	}
	srv.configureView()
	pb.RegisterMultiPaxosServer(srv.Server, &srv)
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
	//slog.Info(fmt.Sprintf("Server started. Listening on address: %s\n\t- peers: %v\n", srv.addr, srv.peers))
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
			srv.proposer = NewProposer(srv.id, srv.peers, srv.acceptor.rnd, srv.View, srv.BroadcastAccept)
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
			srv.proposer = NewProposer(srv.id, srv.peers, srv.acceptor.rnd, srv.View, srv.BroadcastAccept)
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

func (srv *Server) Write(ctx gorums.ServerCtx, request *pb.PaxosValue, broadcast *pb.Broadcast) {
	if !srv.isLeader() {
		// alternatives:
		// 1. simply ignore request 			<- ok
		// 2. send it to the leader 			<- ok
		// 3. reply with last committed value	<- ok
		// 4. reply with error					<- not ok
		broadcast.Forward(request, srv.leader)
		return
	}
	srv.proposer.mut.Lock()
	broadcast.Accept(&pb.AcceptMsg{
		Rnd:  srv.proposer.rnd,
		Slot: srv.proposer.adu,
		Val:  request,
	})
	srv.proposer.adu++
	srv.proposer.mut.Unlock()
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
	//srv.PrintStats()
	metrics := srv.GetStats()
	m := []*pb.Metric{
		{
			TotalNum:              metrics.TotalNum,
			FinishedReqsTotal:     metrics.FinishedReqs.Total,
			FinishedReqsSuccesful: metrics.FinishedReqs.Succesful,
			FinishedReqsFailed:    metrics.FinishedReqs.Failed,
			Processed:             metrics.Processed,
			Dropped:               metrics.Dropped,
			Invalid:               metrics.Invalid,
			AlreadyProcessed:      metrics.AlreadyProcessed,
			RoundTripLatency: &pb.TimingMetric{
				Avg: uint64(metrics.RoundTripLatency.Avg),
				Min: uint64(metrics.RoundTripLatency.Min),
				Max: uint64(metrics.RoundTripLatency.Max),
			},
			ReqLatency: &pb.TimingMetric{
				Avg: uint64(metrics.RoundTripLatency.Avg),
				Min: uint64(metrics.RoundTripLatency.Min),
				Max: uint64(metrics.RoundTripLatency.Max),
			},
			ShardDistribution: metrics.ShardDistribution,
		},
	}
	return &pb.Result{
		Metrics: m,
	}, nil
}
