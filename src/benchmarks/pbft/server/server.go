package server

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	ld "github.com/aleksander-vedvik/benchmark/leaderelection"
	pb "github.com/aleksander-vedvik/benchmark/pbft/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/relab/gorums"
)

// The storage server should implement the server interface defined in the pbbuf files
type Server struct {
	*pb.Server
	leaderElection *ld.MonLeader[*pb.Node, *pb.Configuration, *pb.Heartbeat]
	leader         string
	data           []string
	addr           string
	peers          []string
	messages       int
	addedMsgs      map[string]bool
	messageLog     *MessageLog
	viewNumber     int32
	state          *pb.ClientResponse
	//requestQueue   []*pb.PrePrepareRequest
	//maxLimitOfReqs int
	sequenceNumber int32
	mgr            *pb.Manager
}

// Creates a new StorageServer.
func New(addr string, srvAddresses []string) *Server {
	srv := Server{
		Server:         pb.NewServer(gorums.WithMetrics()),
		data:           make([]string, 0),
		addr:           addr,
		peers:          srvAddresses,
		addedMsgs:      make(map[string]bool),
		leader:         "",
		messageLog:     newMessageLog(),
		state:          nil,
		sequenceNumber: 1,
		viewNumber:     1,
	}
	srv.configureView()
	pb.RegisterPBFTNodeServer(srv.Server, &srv)
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

func (s *Server) Start() {
	lis, err := net.Listen("tcp4", s.addr)
	if err != nil {
		log.Fatal(err)
	}
	//go s.status()
	go s.Serve(lis)
	s.addr = lis.Addr().String()
	var id uint32
	for _, node := range s.View.Nodes() {
		if node.Address() == s.addr {
			id = node.ID()
			break
		}
	}
	slog.Info(fmt.Sprintf("Server started. Listening on address: %s\n\t- peers: %v\n", s.addr, s.peers))
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

func (s *Server) status() {
	for {
		time.Sleep(5 * time.Second)
		state := ""
		if s.state != nil {
			state = s.state.Result
		}
		str := fmt.Sprintf("Server %s running with:\n\t- number of messages: %v\n\t- commited value: %v\n\t- peers: %v\n", s.addr[len(s.addr)-4:], s.messages, state, s.peers)
		log.Println(str)
	}
}

func (s *Server) Write(ctx gorums.ServerCtx, request *pb.WriteRequest, broadcast *pb.Broadcast) {
	if !s.isLeader() {
		if val, ok := s.requestIsAlreadyProcessed(request); ok {
			broadcast.SendToClient(val, nil)
		} else {
			slog.Info("not the leader")
			broadcast.Forward(request, s.leader)
		}
		return
	}
	slog.Info("got client request. initiating a pBFT round")
	req := &pb.PrePrepareRequest{
		Id:             request.Id,
		View:           s.viewNumber,
		SequenceNumber: s.sequenceNumber,
		Digest:         "digest",
		Message:        request.Message,
		Timestamp:      request.Timestamp,
	}
	broadcast.PrePrepare(req)
	s.sequenceNumber++
}

//func (srv *Server) Benchmark(ctx gorums.ServerCtx, request *empty.Empty) (*pb.Result, error) {
////srv.PrintStats()
//metrics := srv.GetStats()
//m := []*pb.Metric{
//{
//TotalNum:              metrics.TotalNum,
//GoroutinesStarted:     metrics.GoroutinesStarted,
//GoroutinesStopped:     metrics.GoroutinesStopped,
//FinishedReqsTotal:     metrics.FinishedReqs.Total,
//FinishedReqsSuccesful: metrics.FinishedReqs.Succesful,
//FinishedReqsFailed:    metrics.FinishedReqs.Failed,
//Processed:             metrics.Processed,
//Dropped:               metrics.Dropped,
//Invalid:               metrics.Invalid,
//AlreadyProcessed:      metrics.AlreadyProcessed,
//RoundTripLatency: &pb.TimingMetric{
//Avg: uint64(metrics.RoundTripLatency.Avg),
//Min: uint64(metrics.RoundTripLatency.Min),
//Max: uint64(metrics.RoundTripLatency.Max),
//},
//ReqLatency: &pb.TimingMetric{
//Avg: uint64(metrics.RoundTripLatency.Avg),
//Min: uint64(metrics.RoundTripLatency.Min),
//Max: uint64(metrics.RoundTripLatency.Max),
//},
//ShardDistribution: metrics.ShardDistribution,
//},
//}
//return &pb.Result{
//Metrics: m,
//}, nil
//}
