package nodeServer

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	"pbft/leaderelection.go"
	pb "pbft/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/relab/gorums"
)

// The storage server should implement the server interface defined in the pbbuf files
type PBFTServer struct {
	*pb.Server
	leaderElection *leaderelection.MonLeader
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
func NewPBFTServer(addr string, srvAddresses []string) *PBFTServer {
	srv := PBFTServer{
		Server:         pb.NewServer(),
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

func (srv *PBFTServer) configureView() {
	srv.mgr = pb.NewManager(
		gorums.WithDialTimeout(50*time.Millisecond),
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

func (s *PBFTServer) Start() {
	lis, err := net.Listen("tcp4", s.addr)
	if err != nil {
		log.Fatal(err)
	}
	//go s.status()
	go s.Serve(lis)
	s.addr = lis.Addr().String()
	slog.Info(fmt.Sprintf("Server started. Listening on address: %s\n\t- peers: %v\n", s.addr, s.peers))
	s.leaderElection = leaderelection.New(s.View)
	s.leaderElection.StartLeaderElection()
	go s.listenForLeaderChanges()
}

func (s *PBFTServer) listenForLeaderChanges() {
	for leader := range s.leaderElection.Leaders() {
		slog.Warn("leader changed", "leader", leader)
		s.leader = leader
	}
}

func (s *PBFTServer) status() {
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

func (s *PBFTServer) Write(ctx gorums.ServerCtx, request *pb.WriteRequest, broadcast *pb.Broadcast) {
	if !s.isLeader() {
		if val, ok := s.requestIsAlreadyProcessed(request); ok {
			broadcast.SendToClient(val, nil)
		} else {
			slog.Info("not the leader")
			broadcast.Write(request, gorums.WithSubset(s.leader))
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
	go broadcast.PrePrepare(req)
	s.sequenceNumber++
}
