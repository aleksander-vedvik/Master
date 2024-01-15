package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/aleksander-vedvik/Master/protos"
)

// The storage server should implement the server interface defined in the pbbuf files
type StorageServer struct {
	data            []string
	addr            string
	messages        int
	handledMessages map[int64]bool
}

// Creates a new StorageServer.
func NewStorageServer(srvAddresses []string, addr string) *StorageServer {
	srv := StorageServer{
		data:            make([]string, 0),
		addr:            addr,
		handledMessages: make(map[int64]bool),
	}
	reliableServer := newReliableServer(srvAddresses, addr)
	reliableServer.RegisterServer(&srv)
	return &srv
}

func (s *StorageServer) StartServer() {
	go s.status()
}

func (s *StorageServer) status() {
	for {
		time.Sleep(5 * time.Second)
		/*val := ""
		if len(s.data) > 0 {
			val = s.data[len(s.data)-1]
		}*/
		str := fmt.Sprintf("Server %s running values: \"%s\"", s.addr[len(s.addr)-4:], s.data)
		log.Println(str)
	}
}

func (s *StorageServer) Write(ctx context.Context, request *pb.State) (response *pb.WriteResponse, err error) {
	s.messages++
	if !s.alreadyAdded(request) {
		s.data = append(s.data, request.Value)
		s.handledMessages[request.Id] = true
	}
	return &pb.WriteResponse{}, nil
}

func (s *StorageServer) Read(ctx context.Context, request *pb.ReadRequest) (response *pb.State, err error) {
	s.messages++
	if len(s.data) <= 0 {
		return &pb.State{}, nil
	}
	response = &pb.State{
		Value:     s.data[len(s.data)-1],
		Timestamp: time.Now().Unix(),
	}
	return response, nil
}

func (s *StorageServer) alreadyAdded(request *pb.State) bool {
	added, ok := s.handledMessages[request.Id]
	return ok && added
}
