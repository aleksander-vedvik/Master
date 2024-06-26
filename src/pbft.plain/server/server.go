package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aleksander-vedvik/benchmark/pbft.plain/config"
	pb "github.com/aleksander-vedvik/benchmark/pbft.plain/protos"
	"github.com/golang/protobuf/ptypes/empty"

	"google.golang.org/grpc"
)

type state int

const (
	PrePrepare state = iota
	Prepare
	Commit
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
	order          map[string]state
	tmpPrepares    map[string][]*pb.PrepareRequest
	tmpCommits     map[string][]*pb.CommitRequest
	log            bool
}

// Creates a new StorageServer.
func New(addr string, srvAddresses []string, withoutLeader ...bool) *Server {
	if len(srvAddresses) < 4 {
		panic("should run with at least 4 servers")
	}
	wL := false
	if len(withoutLeader) > 0 {
		wL = withoutLeader[0]
	}
	useLog := os.Getenv("LOG")
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
		order:          make(map[string]state),
		tmpPrepares:    make(map[string][]*pb.PrepareRequest),
		tmpCommits:     make(map[string][]*pb.CommitRequest),
		log:            useLog == "1",
	}
	pb.RegisterPBFTNodeServer(srv.srv, &srv)
	return &srv
}

func (s *Server) Start(local bool) {
	slog.Info("server: started", "addr", s.addr)
	var (
		lis net.Listener
		err error
	)
	env := os.Getenv("PRODUCTION")
	if env == "1" {
		splittedAddr := strings.Split(s.addr, ":")
		lis, err = net.Listen("tcp", ":"+splittedAddr[1])
	} else {
		lis, err = net.Listen("tcp", s.addr)
	}
	slog.Info(fmt.Sprintf("Server started. Listening on address: %s, %s\n\t- peers: %v\n", s.addr, lis.Addr().String(), s.peers))
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
	if !s.isLeader() {
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
	s.messageLog.add(req, s.viewNumber, req.SequenceNumber)
	s.progressState(request.Id)
	go s.view.PrePrepare(req)
	return nil, nil
}

// only used by the client
func (s *Server) ClientHandler(ctx context.Context, request *pb.ClientResponse) (*empty.Empty, error) {
	return nil, nil
}

func (srv *Server) Benchmark(ctx context.Context, request *empty.Empty) (*pb.Result, error) {
	srv.mut.Lock()
	srv.order = make(map[string]state)
	srv.tmpPrepares = make(map[string][]*pb.PrepareRequest)
	srv.tmpCommits = make(map[string][]*pb.CommitRequest)
	srv.mut.Unlock()
	srv.messageLog.Clear()
	return &pb.Result{}, nil
}
