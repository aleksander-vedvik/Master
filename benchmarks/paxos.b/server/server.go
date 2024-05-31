package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"runtime"
	"sync"

	ld "github.com/aleksander-vedvik/benchmark/leaderelection"
	pb "github.com/aleksander-vedvik/benchmark/paxos.b/proto"

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
}

func New(addr string, srvAddrs []string, logger *slog.Logger) *Server {
	//if logger == nil {
	//loggerOpts := &slog.HandlerOptions{
	//AddSource: true,
	//Level:     slog.LevelDebug,
	//}
	//handler := slog.NewTextHandler(os.Stdout, loggerOpts)
	//logger = slog.New(handler)
	//}
	disable := true
	//if len(disableLeaderElection) > 0 {
	//disable = disableLeaderElection[0]
	//}
	id := 0
	for i, srvAddr := range srvAddrs {
		if addr == srvAddr {
			id = i
			break
		}
	}
	srv := Server{
		id:     uint32(id),
		Server: pb.NewServer(gorums.WithSLogger(logger), gorums.WithListenAddr(addr)),
		//Server:                pb.NewServer(),
		acceptor:              NewAcceptor(addr, len(srvAddrs)),
		data:                  make([]string, 0),
		addr:                  addr,
		peers:                 srvAddrs,
		leader:                srvAddrs[0],
		disableLeaderElection: disable,
	}
	srv.configureView(logger)
	pb.RegisterMultiPaxosServer(srv.Server, &srv)
	return &srv
}

func (srv *Server) configureView(logger *slog.Logger) {
	srv.mgr = pb.NewManager(
		gorums.WithGrpcDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
		gorums.WithLogger(logger),
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

func (srv *Server) Start(local bool) {
	// create listener
	var (
		lis net.Listener
		err error
	)
	env := os.Getenv("PRODUCTION")
	if env == "1" {
		lis, err = net.Listen("tcp4", ":5000")
	} else {
		lis, err = net.Listen("tcp4", srv.addr)
	}
	if err != nil {
		panic(err)
	}
	// add the address
	//srv.addr = lis.Addr().String()
	slog.Info(fmt.Sprintf("Server started. Listening on address: %s, lis=%s\n", srv.addr, lis.Addr().String()))
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
			srv.proposer = NewProposer(srv.id, srv.peers, srv.acceptor.rnd+1, srv.View, srv.BroadcastAccept)
			//go srv.proposer.Start()
		}
	}
	// don't block the caller in case the servers are started locally
	if local {
		go srv.Serve(lis)
		return
	}
	// start gRPC server
	srv.Serve(lis)
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

func (srv *Server) Write(ctx gorums.ServerCtx, request *pb.PaxosValue, broadcast *pb.Broadcast) {
	if !srv.isLeader() {
		// alternatives:
		// 1. simply ignore request 			<- ok
		// 2. send it to the leader 			<- ok
		// 3. reply with last committed value	<- ok
		// 4. reply with error					<- not ok
		//broadcast.Forward(request, srv.leader)
		return
	}
	srv.proposer.mut.Lock()
	rnd := srv.proposer.rnd
	adu := srv.proposer.adu
	srv.proposer.adu++
	srv.proposer.mut.Unlock()

	broadcast.Accept(&pb.AcceptMsg{
		Rnd:  rnd,
		Slot: adu,
		Val:  request,
	})
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
	srv.mut.Lock()
	defer srv.mut.Unlock()
	slog.Info("purging reqs")
	// purge all reqs
	srv.SetView(srv.View)
	srv.acceptor = NewAcceptor(srv.addr, len(srv.peers))
	srv.data = make([]string, 0)
	if srv.leader == srv.addr {
		srv.proposer = NewProposer(srv.id, srv.peers, srv.acceptor.rnd+1, srv.View, srv.BroadcastAccept)
		//go srv.proposer.Start()
	}
	runtime.GC()
	metrics := srv.GetStats()
	m := []*pb.Metric{
		{
			TotalNum: metrics.TotalNum,
			Dropped:  metrics.Dropped,
		},
	}
	return &pb.Result{
		Metrics: m,
	}, nil
}
