package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	pb "github.com/aleksander-vedvik/Master/protos"
	"github.com/relab/gorums"
)

type Data struct {
	Id          int64
	Value       string
	Timestamp   int64
	BroadcastID string
}

func newData(req *pb.State, broadcastID string) *Data {
	return &Data{
		Id:          req.GetId(),
		Value:       req.GetValue(),
		Timestamp:   req.GetTimestamp(),
		BroadcastID: broadcastID,
	}
}

// The storage server should implement the server interface defined in the pbbuf files
type StorageServer struct {
	*pb.Server
	sync.RWMutex
	data     []*Data
	pending  []*Data
	acks     map[int64]int
	messages int
	addr     string
	peers    []string
}

// Creates a new StorageServer.
func NewStorageServer(addr string, srvAddresses []string) *StorageServer {
	srv := StorageServer{
		Server:   pb.NewServer(),
		data:     make([]*Data, 0),
		pending:  make([]*Data, 0),
		acks:     make(map[int64]int),
		messages: 0,
		addr:     "",
		peers:    make([]string, 0),
	}
	pb.RegisterUniformBroadcastServer(srv.Server, &srv)
	otherServers := make([]string, 0, len(srvAddresses)-1)
	for _, srvAddr := range srvAddresses {
		if srvAddr == srv.addr {
			continue
		}
		otherServers = append(otherServers, srvAddr)
	}
	srv.peers = otherServers
	srv.RegisterMiddlewares(srv.authenticate, srv.countMsgs)
	return &srv
}

func (s *StorageServer) authenticate(metadata gorums.BroadcastMetadata) error {
	//log.Println("CTX:", ctx.GetBroadcastValues())
	//log.Println(s.addr, "CTX:", ctx.GetBroadcastValues())
	//bd := ctx.GetBroadcastValues()
	//log.Println(s.addr, bd.Method, bd.SenderAddr)
	return nil
}

func (s *StorageServer) countMsgs(gorums.BroadcastMetadata) error {
	s.Lock()
	defer s.Unlock()
	s.messages++
	return nil
}

func printable(addrs []string) string {
	ret := "[ "
	for _, addr := range addrs {
		num, _ := strconv.Atoi(addr[len(addr)-1:])
		ret += "node " + strconv.Itoa(num+1) + " "
	}
	ret += "]"
	return ret
}

func printablevals(vals []*Data) string {
	ret := "[ "
	for _, val := range vals {
		ret += "\"" + val.Value + "\"" + " "
	}
	ret += "]"
	return ret
}

// Start the server listening on the provided address string
// The function should be non-blocking
// Returns the full listening address of the server as string
// Hint: Use go routine to start the server.
func (s *StorageServer) Start(addr string) {
	lis, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatal(err)
	}
	s.addr = fmt.Sprintf("%v", lis.Addr())
	go s.status()
	go s.deliver()
	log.Printf("Server started. Listening on address: %s\n", s.addr)
	s.Serve(lis)
}

func (s *StorageServer) Run() {
	time.Sleep(10 * time.Millisecond)
	s.SetView(s.addr, s.peers)
}

func (s *StorageServer) status() {
	for {
		time.Sleep(5 * time.Second)
		s.Lock()
		num, _ := strconv.Atoi(s.addr[len(s.addr)-1:])
		str := fmt.Sprintf("\nnode %v running with peers = %v, msgs = %v\n\t- pending:\t%v\n\t- data:\t\t%v\n\t- acks:\t\t%v", num+1, printable(s.peers), s.messages, printablevals(s.pending), printablevals(s.data), s.acks)
		fmt.Println(str)
		s.Unlock()
	}
}

func (s *StorageServer) Broadcast(ctx gorums.ServerCtx, request *pb.State, broadcast *pb.Broadcast) {
	s.Lock()
	defer s.Unlock()
	md := broadcast.GetMetadata()
	s.pending = append(s.pending, newData(request, md.BroadcastID))
	go broadcast.Deliver(request, gorums.WithSubset(s.peers...), gorums.WithGossip(2.0))
}

func (s *StorageServer) Deliver(ctx gorums.ServerCtx, request *pb.State, broadcast *pb.Broadcast) {
	s.Lock()
	defer s.Unlock()
	s.addAck(request.GetId())
	if !s.inPending(request) {
		md := broadcast.GetMetadata()
		s.pending = append(s.pending, newData(request, md.BroadcastID))
		go broadcast.Deliver(request)
	}
}

func (s *StorageServer) addAck(reqId int64) {
	if _, ok := s.acks[reqId]; !ok {
		s.acks[reqId] = 0
	}
	s.acks[reqId]++
}

func (s *StorageServer) inPending(req *pb.State) bool {
	for _, msg := range s.pending {
		if msg.Id == req.GetId() {
			return true
		}
	}
	return false
}

func (srv *StorageServer) deliver() {
	for {
		srv.Lock()
		newPending := make([]*Data, 0, len(srv.pending))
		for _, msg := range srv.pending {
			if srv.canDeliver(msg.Id) {
				srv.data = append(srv.data, msg)
				go srv.ReplyToClient(&pb.ClientResponse{
					Success: true,
					Value:   msg.Value,
				}, nil, msg.BroadcastID)
			} else {
				newPending = append(newPending, msg)
			}
		}
		srv.pending = newPending
		srv.Unlock()
		time.Sleep(6 * time.Second)
	}
}

func (s *StorageServer) canDeliver(reqId int64) bool {
	acks, ok := s.acks[reqId]
	if !ok {
		return false
	}
	return acks >= 1+len(s.peers)/2
}
