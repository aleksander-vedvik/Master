package server

import (
	"net"
	"sync"

	pb "reliablebroadcast/proto"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	*pb.Server
	mut       sync.Mutex
	delivered []*pb.Message
	mgr       *pb.Manager
	addr      string
	srvAddrs  []string
}

func New(addr string, srvAddrs []string) *Server {
	address, _ := net.ResolveTCPAddr("tcp", addr)
	srv := Server{
		Server:    pb.NewServer(gorums.WithListenAddr(address)),
		addr:      addr,
		srvAddrs:  srvAddrs,
		delivered: make([]*pb.Message, 0),
	}
	srv.configureView()
	pb.RegisterReliableBroadcastServer(srv.Server, &srv)
	return &srv
}

func (srv *Server) configureView() {
	srv.mgr = pb.NewManager(
		gorums.WithGrpcDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	view, err := srv.mgr.NewConfiguration(gorums.WithNodeList(srv.srvAddrs))
	if err != nil {
		panic(err)
	}
	srv.SetView(view)
}

func (s *Server) Start() {
	lis, _ := net.Listen("tcp", s.addr)
	go s.Serve(lis)
}

func (s *Server) Broadcast(ctx gorums.ServerCtx, request *pb.Message, broadcast *pb.Broadcast) {
	broadcast.Deliver(request)
}

func (s *Server) Deliver(ctx gorums.ServerCtx, request *pb.Message, broadcast *pb.Broadcast) {
	s.mut.Lock()
	defer s.mut.Unlock()
	if !s.isDelivered(request) {
		s.addDelivered(request)
		broadcast.Deliver(request)
		broadcast.SendToClient(request, nil)
	}
}

func (s *Server) isDelivered(message *pb.Message) bool {
	for _, msg := range s.delivered {
		if msg.Data == message.Data {
			return true
		}
	}
	return false
}

func (s *Server) addDelivered(message *pb.Message) {
	s.delivered = append(s.delivered, message)
}
