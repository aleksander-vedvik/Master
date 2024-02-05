package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "pbft/protos"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// The storage server should implement the server interface defined in the pbbuf files
type PBFTServer struct {
	*pb.Server
	sync.RWMutex
	data            []string
	addr            string
	peers           []string
	messages        int
	handledMessages map[string]map[string]int
	addedMsgs       map[string]bool
	leader          bool
}

// Creates a new StorageServer.
func NewStorageServer(addr string, srvAddresses []string) *PBFTServer {
	handledMessages := make(map[string]map[string]int)
	handledMessages["PrePrepare"] = make(map[string]int)
	handledMessages["Prepare"] = make(map[string]int)
	handledMessages["Commit"] = make(map[string]int)
	srv := PBFTServer{
		Server:          pb.NewServer(),
		data:            make([]string, 0),
		addr:            "",
		peers:           make([]string, 0),
		handledMessages: handledMessages,
		addedMsgs:       make(map[string]bool),
		leader:          true,
	}
	//srv.gorumsSrv.AddTmp(addr)
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

func (s *PBFTServer) authenticate(ctx gorums.BroadcastCtx) error {
	//log.Println("CTX:", ctx.GetBroadcastValues())
	//log.Println(s.addr, "CTX:", ctx.GetBroadcastValues())
	return nil
}

func (s *PBFTServer) countMsgs(ctx gorums.BroadcastCtx) error {
	s.messages++
	bv := ctx.GetBroadcastValues()
	s.addToHandledMessages(bv.Method, bv.BroadcastID)
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
		s.addr = fmt.Sprintf("%s", lis.Addr())
		addrChan <- s.addr
		s.Serve(lis)
	}()
	go s.status()
	return <-addrChan
}

func (s *PBFTServer) Run() {
	time.Sleep(1 * time.Second)
	s.RegisterConfiguration(s.addr, s.peers,
		gorums.WithDialTimeout(50*time.Millisecond),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
}

//func (s *StorageServer) AddConfig(srvAddresses []string) {
//time.Sleep(2 * time.Second)
//otherServers := make([]string, 0, len(srvAddresses)-1)
//for _, srvAddr := range srvAddresses {
//if srvAddr == s.addr {
//continue
//}
//otherServers = append(otherServers, srvAddr)
//}
//s.peers = otherServers
//config := getConfig(s.addr, otherServers)
////config.AddSender(s.addr)
//s.RegisterConfiguration(config)
////s.gorumsSrv.CreateMapping(pb.Map(pb.PrePrepare, pb.Prepare))
//}

func (s *PBFTServer) Start(addr string) {
	lis, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatal(err)
	}
	s.addr = fmt.Sprintf("%s", lis.Addr())
	go s.status()
	log.Printf("Server started. Listening on address: %s\n", s.addr)
	s.Serve(lis)
}

func (s *PBFTServer) status() {
	for {
		time.Sleep(5 * time.Second)
		str := fmt.Sprintf("Server %s running with:\n\t- number of messages: %v\n\t- commited values: %v\n\t- peers: %v\n", s.addr[len(s.addr)-4:], s.messages, s.data, s.peers)
		log.Println(str)
	}
}

func (s *PBFTServer) PrePrepare(ctx gorums.BroadcastCtx, request *pb.PrePrepareRequest, broadcast *pb.Broadcast) (err error) {
	if !s.leader {
		if val, ok := s.requestIsAlreadyProcessed(request); ok {
			broadcast.ReturnToClient(val, nil)
		}
		return nil
	}
	broadcast.Prepare(&pb.PrepareRequest{
		Value: request.GetValue(),
	})
	return nil
}

func (s *PBFTServer) Prepare(ctx gorums.BroadcastCtx, request *pb.PrepareRequest, broadcast *pb.Broadcast) (err error) {
	broadcast.Commit(&pb.CommitRequest{
		Value: request.GetValue(),
	})
	return nil
}

func (s *PBFTServer) Commit(ctx gorums.BroadcastCtx, request *pb.CommitRequest, broadcast *pb.Broadcast) (err error) {
	if s.quorum(ctx) {
		s.addMessage(request.GetValue())
		broadcast.ReturnToClient(&pb.ClientResponse{
			Value: request.GetValue(),
		}, nil)
	}
	return nil
}

func (s *PBFTServer) quorum(ctx gorums.BroadcastCtx) bool {
	bv := ctx.GetBroadcastValues()
	return s.handledMessages[bv.Method][bv.BroadcastID] >= len(s.peers)-1 // does not include itself
}

func (s *PBFTServer) alreadyAdded(val string) bool {
	added, ok := s.addedMsgs[val]
	return ok && added
}

func (s *PBFTServer) addToHandledMessages(method, broadcastID string) {
	if _, ok := s.handledMessages[method][broadcastID]; !ok {
		s.handledMessages[method][broadcastID] = 0
	}
	s.handledMessages[method][broadcastID]++
}

func (s *PBFTServer) addMessage(val string) {
	if s.alreadyAdded(val) {
		return
	}
	s.data = append(s.data, val)
	s.addedMsgs[val] = true
}

func (s *PBFTServer) requestIsAlreadyProcessed(req *pb.PrePrepareRequest) (*pb.ClientResponse, bool) {
	return nil, false
}

/*func (s *StorageServer) prePrepare(ctx gorums.ServerCtx, request *pb.PrePrepareRequest) (response *pb.Empty, err error, broadcast bool) {
	s.messages++
	response = &pb.Empty{}
	err = nil
	broadcast = true
	return
}

func (s *StorageServer) prepare(ctx gorums.ServerCtx, request *pb.PrepareRequest) (response *pb.Empty, err error, broadcast bool) {
	s.messages++
	response = &pb.Empty{}
	err = nil
	broadcast = false
	if _, ok := s.handledMessages["Prepare"][request.Value]; !ok {
		s.handledMessages["Prepare"][request.Value] = 0
	}
	s.handledMessages["Prepare"][request.Value]++
	if s.quorum(request.Value, "Prepare") {
		broadcast = true
	}
	//fmt.Println(s.addr, "received Prepare quorum", broadcast)
	return
}

func (s *StorageServer) commit(ctx gorums.ServerCtx, request *pb.CommitRequest) (response *pb.Empty, err error) {
	fmt.Println(s.addr, "received Commit")
	s.messages++
	response = &pb.Empty{}
	err = nil
	if _, ok := s.handledMessages["Commit"]; !ok {
		s.handledMessages["Commit"][request.Value] = 0
	}
	s.handledMessages["Commit"][request.Value]++
	if s.quorum(request.Value, "Commit") && !s.alreadyAdded(request.GetValue()) {
		s.data = append(s.data, request.GetValue())
		s.addedMsgs[request.GetValue()] = true
	}
	return
}

func (s *StorageServer) PrePrepare2(ctx gorums.ServerCtx, request *pb.PrePrepareRequest, broadcast func(*pb.PrepareRequest)) (response *pb.Empty, err error) {
	s.messages++
	response = &pb.Empty{}
	err = nil
	broadcast(&pb.PrepareRequest{
		Value: request.GetValue(),
	})
	return
}

func (s *StorageServer) Prepare2(ctx gorums.ServerCtx, request *pb.PrepareRequest, broadcast func(*pb.CommitRequest)) (response *pb.Empty, err error) {
	s.messages++
	response = &pb.Empty{}
	err = nil
	s.addToHandledMessages("Prepare", request.GetValue())
	broadcast(&pb.CommitRequest{
		Value: request.GetValue(),
	})
	return
}

func (s *StorageServer) Commit2(ctx gorums.ServerCtx, request *pb.CommitRequest, returnToClient func(*pb.Empty)) (response *pb.Empty, err error) {
	fmt.Println(s.addr, "received Commit")
	s.messages++
	response = &pb.Empty{}
	err = nil
	s.addToHandledMessages("Commit", request.GetValue())
	if s.quorum(request.Value, "Commit") && !s.alreadyAdded(request.GetValue()) {
		s.addMessage(request.GetValue())
		returnToClient(&pb.Empty{})
		return response, fmt.Errorf("successfully returned to client")
	}
	return
}

func (s *StorageServer) PrePrepare(ctx gorums.ServerCtx, request *pb.PrePrepareRequest, broadcast *pb.Broadcast) (response *pb.ClientResponse, err error) {
	s.messages++
	response = &pb.ClientResponse{}
	err = nil
	broadcast.Prepare(&pb.PrepareRequest{
		Value: request.GetValue(),
	})
	return
}

func (s *StorageServer) Prepare(ctx gorums.ServerCtx, request *pb.PrepareRequest, broadcast *pb.Broadcast) (response *pb.Empty, err error) {
	s.messages++
	response = &pb.Empty{}
	err = nil
	s.addToHandledMessages("Prepare", request.GetValue())
	broadcast.Commit(&pb.CommitRequest{
		Value: request.GetValue(),
	})
	return
}

func (s *StorageServer) Commit(ctx gorums.ServerCtx, request *pb.CommitRequest, broadcast *pb.Broadcast) (response *pb.Empty, err error) {
	//fmt.Println(s.addr, "received Commit")
	s.messages++
	response = &pb.Empty{}
	err = fmt.Errorf("test")
	s.addToHandledMessages("Commit", request.GetValue())
	if s.quorum(request.Value, "Commit") && !s.alreadyAdded(request.GetValue()) {
		err = fmt.Errorf("successfully returned to client")
		s.addMessage(request.GetValue())
		broadcast.ReturnToClient(&pb.ClientResponse{
			Value: request.GetValue(),
		}, err)
	}
	broadcast.ReturnToClient(&pb.ClientResponse{
		Value: request.GetValue(),
	}, err)
	return
}*/
