package storage

import (
	"context"
	"log"
	"time"

	pb "github.com/aleksander-vedvik/Master/protos"
)

type StorageClient struct {
	view *pb.Configuration
}

// Creates a new StorageClient with the provided srvAddresses as the configuration
func NewStorageClient(srvAddresses []string) *StorageClient {
	return &StorageClient{
		view: getConfig(srvAddresses),
	}
}

// Writes the provided value to a random server
func (sc *StorageClient) WriteValue(value string) error {
	_, err := sc.view.Write(context.Background(), &pb.State{
		Value:     value,
		Timestamp: time.Now().Unix(),
	})
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// Returns a slice of values stored on all servers
func (sc *StorageClient) ReadValue() (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	reply, err := sc.view.Read(ctx, &pb.ReadRequest{})
	defer cancel()
	if err != nil {
		log.Fatalln("read rpc returned error:", err)
		return "", nil
	}
	return reply.Value, nil
}
