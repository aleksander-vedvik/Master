package server

import (
	"context"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/aleksander-vedvik/benchmark/pbft.s/config"
	pb "github.com/aleksander-vedvik/benchmark/pbft.s/protos"
	"github.com/golang/protobuf/ptypes/empty"

	"google.golang.org/grpc"
)

// The storage server should implement the server interface defined in the pbbuf files
type Server struct {
	pb.PBFTNodeServer
	mut            sync.Mutex
	leader         string
	data           []string
	addr           string
	peers          []string
	addedMsgs      map[string]bool
	messageLog     *MessageLog
	viewNumber     int32
	state          *pb.ClientResponse
	sequenceNumber int32
	withoutLeader  bool
	view           *config.Config
	srv            *grpc.Server
}

// Creates a new StorageServer.
func New(addr string, srvAddresses []string, withoutLeader ...bool) *Server {
	wL := false
	if len(withoutLeader) > 0 {
		wL = withoutLeader[0]
	}
	srv := Server{
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
		view:           config.NewConfig(addr, srvAddresses),
		srv:            grpc.NewServer(),
	}
	pb.RegisterPBFTNodeServer(srv.srv, &srv)
	return &srv
}

func (s *Server) Start(local bool) {
	slog.Info("server: started", "addr", s.addr)
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		panic(err)
	}
	if local {
		go s.srv.Serve(lis)
		time.Sleep(1 * time.Second)
		return
	}
	s.srv.Serve(lis)
}

func (s *Server) Stop() {
	s.srv.Stop()
}

func (s *Server) Write(ctx context.Context, request *pb.WriteRequest) (*empty.Empty, error) {
	//slog.Info("0 server: received msg")
	if !s.isLeader() {
		s.mut.Lock()
		if val, ok := s.requestIsAlreadyProcessed(request); ok {
			go s.view.ClientHandler(val)
			//} else {
			//broadcast.Forward(request, s.leader)
		}
		s.mut.Unlock()
		return nil, nil
	}
	s.mut.Lock()
	req := &pb.PrePrepareRequest{
		Id:             request.Id,
		View:           s.viewNumber,
		SequenceNumber: s.sequenceNumber,
		Digest:         "digest",
		Message:        request.Message,
		Timestamp:      request.Timestamp,
		From:           request.From,
	}
	s.sequenceNumber++
	s.mut.Unlock()
	s.view.PrePrepare(req)
	return nil, nil
}

// only used by the client
func (s *Server) ClientHandler(ctx context.Context, request *pb.ClientResponse) (*empty.Empty, error) {
	return nil, nil
}

func (srv *Server) Benchmark(ctx context.Context, request *empty.Empty) (*pb.Result, error) {
	srv.messageLog.Clear()
	return &pb.Result{}, nil
}
