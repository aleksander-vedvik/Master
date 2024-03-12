package nodeServer

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"sync"
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
	mu             sync.RWMutex
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
	requestQueue   []*pb.PrePrepareRequest
	maxLimitOfReqs int
	sequenceNumber int32
	mgr            *pb.Manager
}

// Creates a new StorageServer.
func NewStorageServer(addr string, srvAddresses []string) *PBFTServer {
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
	srv.mgr = pb.NewManager(
		gorums.WithDialTimeout(51*time.Millisecond),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	view, err := srv.mgr.NewConfiguration(gorums.WithNodeListBroadcast(srvAddresses))
	if err != nil {
		panic(err)
	}
	srv.SetView(view)
	pb.RegisterPBFTNodeServer(srv.Server, &srv)
	return &srv
}

func (s *PBFTServer) Start() {
	lis, err := net.Listen("tcp4", s.addr)
	if err != nil {
		log.Fatal(err)
	}
	go s.status()
	go s.Serve(lis)
	s.addr = lis.Addr().String()
	log.Printf("Server started. Listening on address: %s\n", s.addr)
	s.leaderElection = leaderelection.New(s.View)
	s.leaderElection.StartLeaderElection()
	go s.listenForLeaderChanges()
}

func (s *PBFTServer) listenForLeaderChanges() {
	for leader := range s.leaderElection.Leaders() {
		slog.Info("leader changed:", leader)
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
			broadcast.Write(request, gorums.WithSubset(s.leader))
		}
		return
	}
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
