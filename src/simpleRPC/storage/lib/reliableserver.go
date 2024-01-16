package lib

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/aleksander-vedvik/Master/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type RServer interface {
	Write(context.Context, *pb.State) (*pb.WriteResponse, error)
	Read(context.Context, *pb.ReadRequest) (*pb.State, error)
}

// The storage server should implement the server interface defined in the pbbuf files
type ReliableServer struct {
	sync.RWMutex
	pb.UnimplementedStorageServer
	addr           string
	recievedFrom   map[int64]map[string]bool
	server         RServer
	multipartyChan chan any
	view           *View

	fireAndForgetChan chan FF
}

// Creates a new StorageServer.
func NewReliableServer(srvAddresses []string, addr string) *ReliableServer {
	otherServers := make([]string, 0, len(srvAddresses)-1)
	for _, srvAddr := range srvAddresses {
		if srvAddr == addr {
			continue
		}
		otherServers = append(otherServers, srvAddr)
	}
	srv := ReliableServer{
		addr:              addr,
		recievedFrom:      make(map[int64]map[string]bool),
		multipartyChan:    make(chan any),
		view:              NewView(otherServers),
		fireAndForgetChan: make(chan FF, 1000),
	}
	go srv.startServer(addr)
	go srv.multiparty()
	go srv.fireAndForgetQueue()
	return &srv
}

func (s *ReliableServer) RegisterServer(srv RServer) {
	s.server = srv
}

// Start the server listening on the provided address string
// The function should be non-blocking
// Returns the full listening address of the server as string
// Hint: Use go routine to start the server.
func (s *ReliableServer) startServer(addr string) {
	lis, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterStorageServer(grpcServer, s)
	grpcServer.Serve(lis)
}

func (s *ReliableServer) Read(ctx context.Context, request *pb.ReadRequest) (response *pb.State, err error) {
	s.Lock()
	defer s.Unlock()
	return s.server.Read(ctx, request)
}

func (s *ReliableServer) Write(ctx context.Context, request *pb.State) (response *pb.WriteResponse, err error) {
	s.Lock()
	defer s.Unlock()
	if s.shouldBroadcast(request.GetId()) {
		go s.broadcast(request)
	}
	if !s.alreadyReceivedFromPeer(ctx, request.GetId()) {
		response, err = s.server.Write(ctx, request)
	}
	return response, err
}

func (s *ReliableServer) shouldBroadcast(msgId int64) bool {
	_, ok := s.recievedFrom[msgId]
	if !ok {
		s.recievedFrom[msgId] = make(map[string]bool)
		return true
	}
	return false
}

func (s *ReliableServer) alreadyReceivedFromPeer(ctx context.Context, msgId int64) bool {
	p, _ := peer.FromContext(ctx)
	addr := p.Addr.String()
	receivedMsgFromNode, ok := s.recievedFrom[msgId][addr]
	if !ok {
		if !receivedMsgFromNode {
			s.recievedFrom[msgId][addr] = true
			return false
		}
	}
	return true

}

func (s *ReliableServer) broadcast(request *pb.State) {
	time.Sleep(5 * time.Second)
	log.Println("broadcasting", request.Value, "from", s.addr[len(s.addr)-4:])
	s.multipartyChan <- request
}

// running in a go routine
func (s *ReliableServer) multiparty() {
	for msg := range s.multipartyChan {
		state := msg.(*pb.State)
		success := s.view.Write(state.GetValue(), state.GetId())
		if !success {
			fmt.Println("write failed")
		}
	}
}

type FF struct {
	ctx     context.Context
	request *pb.State
}

func (s *ReliableServer) WriteFireAndForget(ctx context.Context, request *pb.State) (response *pb.WriteResponse, err error) {
	s.fireAndForgetChan <- FF{ctx, request}
	return
}

func (s *ReliableServer) fireAndForgetQueue() {
	for ff := range s.fireAndForgetChan {
		if s.shouldBroadcast(ff.request.GetId()) {
			go s.broadcast(ff.request)
		}
		if !s.alreadyReceivedFromPeer(ff.ctx, ff.request.GetId()) {
			_, _ = s.server.Write(ff.ctx, ff.request)
		}
	}
}
