package nodeServer

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "pbft/protos"

	"github.com/relab/gorums"
)

// The storage server should implement the server interface defined in the pbbuf files
type PBFTServer struct {
	*pb.Server
	sync.RWMutex
	data           []string
	addr           string
	peers          []string
	messages       int
	addedMsgs      map[string]bool
	leader         string
	messageLog     *MessageLog
	viewNumber     int32
	state          *pb.ClientResponse
	requestQueue   []*pb.PrePrepareRequest
	maxLimitOfReqs int
	sequenceNumber int32
}

// Creates a new StorageServer.
func NewStorageServer(addr string, srvAddresses []string) *PBFTServer {
	srv := PBFTServer{
		Server:         pb.NewServer(),
		data:           make([]string, 0),
		addr:           "",
		peers:          make([]string, 0),
		addedMsgs:      make(map[string]bool),
		leader:         "127.0.0.1:5000",
		messageLog:     newMessageLog(),
		state:          nil,
		sequenceNumber: 1,
		viewNumber:     1,
	}
	otherServers := make([]string, 0, len(srvAddresses)-1)
	for _, srvAddr := range srvAddresses {
		if srvAddr == srv.addr {
			continue
		}
		otherServers = append(otherServers, srvAddr)
	}
	srv.peers = otherServers
	srv.RegisterMiddlewares(srv.authenticate, srv.countMsgs)
	pb.RegisterPBFTNodeServer(srv.Server, &srv)
	return &srv
}

func (s *PBFTServer) authenticate(ctx gorums.BroadcastMetadata) error {
	//log.Println("CTX:", ctx.GetBroadcastValues())
	//log.Println(s.addr, "CTX:", ctx.GetBroadcastValues())
	return nil
}

func (s *PBFTServer) countMsgs(ctx gorums.BroadcastMetadata) error {
	s.messages++
	//bv := ctx.GetBroadcastValues()
	//s.addToHandledMessages(bv.Method, bv.BroadcastID)
	//bv := ctx.GetBroadcastValues()
	//log.Println(s.addr, bv.Method)
	return nil
}

// Start the server listening on the provided address string
// The function should be non-blocking
// Returns the full listening address of the server as string
// Hint: Use go routine to start the server.
func (s *PBFTServer) StartServer(addr string) string {
	addrChan := make(chan string)
	go func() {
		lis, err := net.Listen("tcp4", addr)
		if err != nil {
			log.Fatal(err)
			addrChan <- ""
			return
		}
		s.addr = fmt.Sprintf("%v", lis.Addr())
		addrChan <- s.addr
		s.Serve(lis)
	}()
	go s.status()
	s.SetView(<-addrChan, s.peers)
	//s.RegisterConfiguration(<-addrChan, s.peers,
	//	gorums.WithDialTimeout(50*time.Millisecond),
	//	gorums.WithGrpcDialOptions(
	//		grpc.WithBlock(),
	//		grpc.WithTransportCredentials(insecure.NewCredentials()),
	//	),
	//)
	return s.addr
}

//func (s *PBFTServer) Run() {
//	time.Sleep(10 * time.Millisecond)
//	s.RegisterConfiguration(s.addr, s.peers,
//		gorums.WithDialTimeout(50*time.Millisecond),
//		gorums.WithGrpcDialOptions(
//			grpc.WithBlock(),
//			grpc.WithTransportCredentials(insecure.NewCredentials()),
//		),
//	)
//}

func (s *PBFTServer) Start(addr string) {
	lis, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatal(err)
	}
	s.addr = fmt.Sprintf("%v", lis.Addr())
	go s.status()
	log.Printf("Server started. Listening on address: %s\n", s.addr)
	s.Serve(lis)
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

func (s *PBFTServer) PrePrepare(ctx gorums.ServerCtx, request *pb.PrePrepareRequest, broadcast *pb.Broadcast) {
	if !s.isInView(request.View) {
		return
	}
	if !s.sequenceNumberIsValid(request.SequenceNumber) {
		return
	}
	if s.hasAlreadyAcceptedSequenceNumber(request.SequenceNumber) {
		return
	}
	s.messageLog.add(request, s.viewNumber, request.SequenceNumber)
	broadcast.Prepare(&pb.PrepareRequest{
		Id:             request.Id,
		View:           request.View,
		SequenceNumber: request.SequenceNumber,
		Digest:         request.Digest,
		Message:        request.Message,
		Timestamp:      request.Timestamp,
	})
}

func (s *PBFTServer) Prepare(ctx gorums.ServerCtx, request *pb.PrepareRequest, broadcast *pb.Broadcast) {
	if !s.isInView(request.View) {
		return
	}
	if !s.sequenceNumberIsValid(request.SequenceNumber) {
		return
	}
	s.messageLog.add(request, s.viewNumber, request.SequenceNumber)
	if s.prepared(request.SequenceNumber) {
		broadcast.Commit(&pb.CommitRequest{
			Id:             request.Id,
			Timestamp:      request.Timestamp,
			View:           request.View,
			SequenceNumber: request.SequenceNumber,
			Message:        request.Message,
		})
	}
}

func (s *PBFTServer) Commit(ctx gorums.ServerCtx, request *pb.CommitRequest, broadcast *pb.Broadcast) {
	if !s.isInView(request.View) {
		return
	}
	if !s.sequenceNumberIsValid(request.SequenceNumber) {
		return
	}
	s.messageLog.add(request, s.viewNumber, request.SequenceNumber)
	if s.committed(request.SequenceNumber) {
		s.state = &pb.ClientResponse{
			Id:        request.Id,
			Result:    request.Message,
			Timestamp: request.Timestamp,
			View:      request.View,
		}
		broadcast.SendToClient(s.state, nil)
	}
}

func (s *PBFTServer) requestIsAlreadyProcessed(req *pb.WriteRequest) (*pb.ClientResponse, bool) {
	return s.state, s.state != nil && s.state.Id == req.Id
}

func (s *PBFTServer) getLeaderAddr() string {
	return s.leader
}

func (s *PBFTServer) isLeader() bool {
	return s.leader == s.addr
}

func (s *PBFTServer) isInView(view int32) bool {
	// request is in the current view of this node
	return s.viewNumber == view
}

func (s *PBFTServer) sequenceNumberIsValid(n int32) bool {
	// the sequence number is between h and H
	return true
}

func (s *PBFTServer) hasAlreadyAcceptedSequenceNumber(n int32) bool {
	// checks if the node has already accepted a pre-prepare request
	// for this view and sequence number
	_, found := s.messageLog.getPrePrepareReq(n, s.viewNumber)
	return found
}

func (s *PBFTServer) prepared(n int32) bool {
	// the request m, a pre-prepare for m in view v with sequence number n,
	// and 2f prepares from different backups that match the pre-prepare.
	// The replicas verify whether the prepares match the pre-prepare by
	// checking that they have the same view, sequence number, and digest.
	req, found := s.messageLog.getPrePrepareReq(n, s.viewNumber)
	if !found {
		return false
	}
	reqs, found := s.messageLog.getPrepareReqs(req.Digest, req.SequenceNumber, req.View)
	return found && len(reqs) > 2*len(s.peers)/3
}

func (s *PBFTServer) committed(n int32) bool {
	/*
		We define the committed and committed-local predi-
		cates as follows: committed(m, v, n) is true if and only
		if prepared(m, v, n, i) is true for all i in some set of
		f+1 non-faulty replicas; and committed-local(m, v, n, i)
		is true if and only if prepared(m, v, n, i) is true and i has
		accepted 2f + 1 commits (possibly including its own)
		from different replicas that match the pre-prepare for m;
		a commit matches a pre-prepare if they have the same
		view, sequence number, and digest.
	*/
	req, found := s.messageLog.getPrePrepareReq(n, s.viewNumber)
	if !found {
		return false
	}
	reqs, found := s.messageLog.getPrepareReqs(req.Digest, req.SequenceNumber, req.View)
	if !found || len(reqs) <= 2*len(s.peers)/3 {
		return false
	}
	commits, found := s.messageLog.getCommitReqs(req.Digest, req.SequenceNumber, req.View)
	return found && len(commits) > 2*len(s.peers)/3
}
