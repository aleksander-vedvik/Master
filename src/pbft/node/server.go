package node

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
type StorageServer struct {
	*pb.Server
	sync.RWMutex
	data            []string
	addr            string
	peers           []string
	messages        int
	handledMessages map[string]map[string]int
	addedMsgs       map[string]bool
}

// Creates a new StorageServer.
func NewStorageServer(addr string) *StorageServer {
	handledMessages := make(map[string]map[string]int)
	handledMessages["PrePrepare"] = make(map[string]int)
	handledMessages["Prepare"] = make(map[string]int)
	handledMessages["Commit"] = make(map[string]int)
	srv := StorageServer{
		Server:          pb.NewServer(),
		data:            make([]string, 0),
		addr:            "",
		peers:           make([]string, 0),
		handledMessages: handledMessages,
		addedMsgs:       make(map[string]bool),
	}
	//srv.gorumsSrv.AddTmp(addr)
	pb.RegisterPBFTNodeServer(srv.Server, &srv)
	return &srv
}

// Start the server listening on the provided address string
// The function should be non-blocking
// Returns the full listening address of the server as string
// Hint: Use go routine to start the server.
func (s *StorageServer) StartServer(addr string) string {
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

func (s *StorageServer) AddConfig(srvAddresses []string) {
	time.Sleep(2 * time.Second)
	otherServers := make([]string, 0, len(srvAddresses)-1)
	for _, srvAddr := range srvAddresses {
		if srvAddr == s.addr {
			continue
		}
		otherServers = append(otherServers, srvAddr)
	}
	s.peers = otherServers
	config := getConfig(s.addr, otherServers)
	//config.AddSender(s.addr)
	s.RegisterConfiguration(config)
	//s.gorumsSrv.CreateMapping(pb.Map(pb.PrePrepare, pb.Prepare))
}

func (s *StorageServer) Start(addr string) {
	lis, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatal(err)
	}
	s.addr = fmt.Sprintf("%s", lis.Addr())
	go s.status()
	log.Printf("Server started. Listening on address: %s\n", s.addr)
	s.Serve(lis)
}

func (s *StorageServer) status() {
	for {
		time.Sleep(5 * time.Second)
		str := fmt.Sprintf("Server %s running with:\n\t- number of messages: %v\n\t- commited values: %v\n\t- peers: %v", s.addr[len(s.addr)-4:], s.messages, s.data, s.peers)
		log.Println(str)
	}
}

func (s *StorageServer) prePrepare(ctx gorums.ServerCtx, request *pb.PrePrepareRequest) (response *pb.Empty, err error, broadcast bool) {
	s.messages++
	response = &pb.Empty{}
	err = nil
	broadcast = true
	return
}

func (s *StorageServer) PrePrepare(ctx gorums.ServerCtx, request *pb.PrePrepareRequest, broadcast func(*pb.PrepareRequest)) (response *pb.Empty, err error) {
	s.messages++
	response = &pb.Empty{}
	err = nil
	broadcast(&pb.PrepareRequest{
		Value: request.GetValue(),
	})
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

func (s *StorageServer) Prepare(ctx gorums.ServerCtx, request *pb.PrepareRequest, broadcast func(*pb.CommitRequest)) (response *pb.Empty, err error) {
	s.messages++
	response = &pb.Empty{}
	err = nil
	if _, ok := s.handledMessages["Prepare"][request.Value]; !ok {
		s.handledMessages["Prepare"][request.Value] = 0
	}
	s.handledMessages["Prepare"][request.Value]++
	if s.quorum(request.Value, "Prepare") {
		broadcast(&pb.CommitRequest{
			Value: request.GetValue(),
		})
	}
	//fmt.Println(s.addr, "received Prepare quorum", broadcast)
	return
}

func (s *StorageServer) Commit(ctx gorums.ServerCtx, request *pb.CommitRequest) (response *pb.Empty, err error) {
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

func (s *StorageServer) quorum(id, step string) bool {
	return s.handledMessages[step][id] >= len(s.peers)-1 // does not include itself
}

func (s *StorageServer) alreadyAdded(val string) bool {
	added, ok := s.addedMsgs[val]
	return ok && added
}
