package server

import (
	"log"
	"net"

	pb "github.com/aleksander-vedvik/benchmark/simple/protos"
	"github.com/golang/protobuf/ptypes/empty"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/relab/gorums"
)

// The storage server should implement the server interface defined in the pbbuf files
type Server struct {
	*pb.Server
	mgr   *pb.Manager
	addr  string
	peers []string
}

// Creates a new StorageServer.
func New(addr string, srvAddresses []string) *Server {
	srv := Server{
		Server: pb.NewServer(gorums.WithMetrics()),
		addr:   addr,
		peers:  srvAddresses,
	}
	srv.configureView()
	pb.RegisterSimpleServer(srv.Server, &srv)
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
}

func (s *Server) BroadcastCall1(ctx gorums.ServerCtx, request *pb.WriteRequest1, broadcast *pb.Broadcast) {
	broadcast.BroadcastIntermediate(&pb.BroadcastRequest{
		Id:      request.Id,
		Message: request.Message,
	})
}

func (s *Server) BroadcastCall2(ctx gorums.ServerCtx, request *pb.WriteRequest2, broadcast *pb.Broadcast) {
	broadcast.BroadcastIntermediate(&pb.BroadcastRequest{
		Id:   request.Id,
		Data: request.Data,
	})
}

func (s *Server) BroadcastIntermediate(ctx gorums.ServerCtx, request *pb.BroadcastRequest, broadcast *pb.Broadcast) {
	broadcast.Broadcast(request)
}

func (s *Server) Broadcast(ctx gorums.ServerCtx, request *pb.BroadcastRequest, broadcast *pb.Broadcast) {
	switch request.Data {
	case nil:
		// if data is empty, then the path is: BroadcastCall1 -> BroadcastIntermediate -> Broadcast
		broadcast.SendToClient(&pb.WriteResponse1{
			Id:      request.Id,
			Message: request.Message,
		}, nil)
	default:
		// if data is non-empty, then the path is: BroadcastCall2 -> BroadcastIntermediate -> Broadcast
		broadcast.SendToClient(&pb.WriteResponse2{
			Id:   request.Id,
			Data: request.Data,
		}, nil)
	}
}

func (srv *Server) Benchmark(ctx gorums.ServerCtx, request *empty.Empty) (*pb.Result, error) {
	metrics := srv.GetStats()
	m := []*pb.Metric{
		{
			Addr:                  srv.addr,
			TotalNum:              metrics.TotalNum,
			GoroutinesStarted:     metrics.GoroutinesStarted,
			GoroutinesStopped:     metrics.GoroutinesStopped,
			FinishedReqsTotal:     metrics.FinishedReqs.Total,
			FinishedReqsSuccesful: metrics.FinishedReqs.Succesful,
			FinishedReqsFailed:    metrics.FinishedReqs.Failed,
			Processed:             metrics.Processed,
			Dropped:               metrics.Dropped,
			Invalid:               metrics.Invalid,
			AlreadyProcessed:      metrics.AlreadyProcessed,
			RoundTripLatency: &pb.TimingMetric{
				Avg: uint64(metrics.RoundTripLatency.Avg.Microseconds()),
				Min: uint64(metrics.RoundTripLatency.Min.Microseconds()),
				Max: uint64(metrics.RoundTripLatency.Max.Microseconds()),
			},
			ReqLatency: &pb.TimingMetric{
				Avg: uint64(metrics.RoundTripLatency.Avg.Microseconds()),
				Min: uint64(metrics.RoundTripLatency.Min.Microseconds()),
				Max: uint64(metrics.RoundTripLatency.Max.Microseconds()),
			},
			ShardDistribution: metrics.ShardDistribution,
		},
	}
	return &pb.Result{
		Metrics: m,
	}, nil
}
