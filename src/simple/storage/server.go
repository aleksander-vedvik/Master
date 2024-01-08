package storage

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/aleksander-vedvik/Master/protos"
	"github.com/relab/gorums"
)

// The storage server should implement the server interface defined in the pbbuf files
type StorageServer struct {
	sync.RWMutex
	data      []string
	gorumsSrv *gorums.Server
	addr      string
}

// Creates a new StorageServer.
func NewStorageServer() *StorageServer {
	gorumsSrv := gorums.NewServer()
	srv := StorageServer{
		data:      make([]string, 0),
		gorumsSrv: gorumsSrv,
		addr:      "",
	}
	pb.RegisterQCStorageServer(gorumsSrv, &srv)
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
		s.gorumsSrv.Serve(lis)
	}()
	go s.status()

	return <-addrChan
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
	s.addr = fmt.Sprintf("%s", lis.Addr())
	go s.status()
	log.Printf("Server started. Listening on address: %s\n", s.addr)
	s.gorumsSrv.Serve(lis)
}

func (s *StorageServer) status() {
	for {
		time.Sleep(5 * time.Second)
		val := ""
		if len(s.data) > 0 {
			val = s.data[len(s.data)-1]
		}
		str := fmt.Sprintf("Server %s running with last value: \"%s\"", s.addr[len(s.addr)-4:], val)
		log.Println(str)
	}
}

// Returns the data slice on this server
func (s *StorageServer) GetData() []string {
	s.RLock()
	defer s.RUnlock()
	return s.data
}

// Sets the data slice to a value
func (s *StorageServer) SetData(data []string) {
	s.Lock()
	defer s.Unlock()
	s.data = data
}

func (s *StorageServer) Write(ctx gorums.ServerCtx, request *pb.State) (response *pb.WriteResponse, err error) {
	s.Lock()
	defer s.Unlock()
	s.data = append(s.data, request.Value)
	return &pb.WriteResponse{New: true}, nil
}

func (s *StorageServer) Read(ctx gorums.ServerCtx, request *pb.ReadRequest) (response *pb.State, err error) {
	s.Lock()
	defer s.Unlock()
	if len(s.data) <= 0 {
		return &pb.State{}, nil
	}
	response = &pb.State{
		Value:     s.data[len(s.data)-1],
		Timestamp: time.Now().Unix(),
	}
	return response, nil
}