package server

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"sync"

	ld "github.com/aleksander-vedvik/benchmark/leaderelection"
	pb "github.com/aleksander-vedvik/benchmark/pbft.o/protos"
	"github.com/golang/protobuf/ptypes/empty"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/relab/gorums"
)

// The storage server should implement the server interface defined in the pbbuf files
type Server struct {
	*pb.Server
	mut            sync.Mutex
	leaderElection *ld.MonLeader[*pb.Node, *pb.Configuration, *pb.Heartbeat]
	leader         string
	data           []string
	addr           string
	peers          []string
	addedMsgs      map[string]bool
	messageLog     *MessageLog
	viewNumber     int32
	state          *pb.ClientResponse
	sequenceNumber int32
	mgr            *pb.Manager
	withoutLeader  bool
}

// Creates a new StorageServer.
func New(addr string, srvAddresses []string, logger *slog.Logger) *Server {
	wL := true
	address, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}
	srv := Server{
		Server:         pb.NewServer(gorums.WithListenAddr(address), gorums.WithOrder(pb.PBFTPrePrepare, pb.PBFTPrepare, pb.PBFTCommit), gorums.WithSLogger(logger)),
		data:           make([]string, 0),
		addr:           addr,
		peers:          srvAddresses,
		addedMsgs:      make(map[string]bool),
		leader:         srvAddresses[0],
		messageLog:     newMessageLog(),
		state:          nil,
		sequenceNumber: 1,
		viewNumber:     1,
		withoutLeader:  wL,
	}
	srv.configureView()
	pb.RegisterPBFTServer(srv.Server, &srv)
	return &srv
}

func (srv *Server) configureView() {
	srv.mgr = pb.NewManager(
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

func (s *Server) Start(local bool) {
	var (
		lis net.Listener
		err error
	)
	env := os.Getenv("PRODUCTION")
	if env == "1" {
		splittedAddr := strings.Split(s.addr, ":")
		//lis, err = net.Listen("tcp4", ":5000")
		lis, err = net.Listen("tcp", ":"+splittedAddr[1])
	} else {
		lis, err = net.Listen("tcp", s.addr)
	}
	if err != nil {
		panic(err)
	}
	//go s.status()
	slog.Info(fmt.Sprintf("Server started. Listening on address: %s\n\t- peers: %v\n", s.addr, s.peers))
	if local {
		go s.Serve(lis)
		return
	}
	s.Serve(lis)
	if s.withoutLeader {
		return
	}
	var id uint32
	for _, node := range s.View.Nodes() {
		if node.Address() == s.addr {
			id = node.ID()
			break
		}
	}
	s.leaderElection = ld.New(s.View, id, func(id uint32) *pb.Heartbeat {
		return &pb.Heartbeat{
			Id: id,
		}
	})
	s.leaderElection.StartLeaderElection()
	go s.listenForLeaderChanges()
}

func (s *Server) listenForLeaderChanges() {
	for leader := range s.leaderElection.Leaders() {
		slog.Warn("leader changed", "leader", leader)
		s.leader = leader
	}
}

func (s *Server) Write(ctx gorums.ServerCtx, request *pb.WriteRequest, broadcast *pb.Broadcast) {
	if !s.isLeader() {
		if val, ok := s.requestIsAlreadyProcessed(request); ok {
			broadcast.SendToClient(val, nil)
		} else {
			broadcast.Forward(request, s.leader)
		}
		return
	}
	s.mut.Lock()
	req := &pb.PrePrepareRequest{
		Id:             request.Id,
		View:           s.viewNumber,
		SequenceNumber: s.sequenceNumber,
		Digest:         "digest",
		Message:        request.Message,
		Timestamp:      request.Timestamp,
	}
	s.sequenceNumber++
	s.mut.Unlock()
	broadcast.PrePrepare(req)
}

func (srv *Server) Benchmark(ctx gorums.ServerCtx, request *empty.Empty) (*pb.Result, error) {
	srv.mut.Lock()
	defer srv.mut.Unlock()
	slog.Info("purging reqs")
	// purge all reqs
	srv.SetView(srv.View)
	srv.messageLog.Clear()
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
