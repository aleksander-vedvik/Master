package gorums

import (
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/aleksander-vedvik/Master/protos"
	"github.com/relab/gorums"
	"google.golang.org/protobuf/types/known/emptypb"
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
	pb.RegisterStorageServiceServer(gorumsSrv, &srv)
	return &srv
}

// Start the server listening on the provided address string
// The function should be non-blocking
// Returns the full listening address of the server as string
// Hint: Use go routine to start the server.
func (s *StorageServer) StartServer(addr string) string {
	addrChan := make(chan string, 0)
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

	return <-addrChan
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

func (s *StorageServer) Write(ctx gorums.ServerCtx, request *pb.WriteRequest) (response *emptypb.Empty, err error) {
	s.Lock()
	defer s.Unlock()
	s.data = append(s.data, request.Value)
	return &emptypb.Empty{}, nil
}

func (s *StorageServer) Read(ctx gorums.ServerCtx, request *emptypb.Empty) (response *pb.ReadResponse, err error) {
	s.Lock()
	defer s.Unlock()
	response = &pb.ReadResponse{
		Values: s.data,
	}
	return response, nil
}
