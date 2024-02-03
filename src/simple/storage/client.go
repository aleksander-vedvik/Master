package storage

import (
	"context"
	"time"

	pb "github.com/aleksander-vedvik/Master/protos"
)

type StorageClient struct {
	view   *pb.Configuration
	msgIds int64
}

// Creates a new StorageClient with the provided srvAddresses as the configuration
func NewStorageClient(srvAddresses []string) *StorageClient {
	return &StorageClient{
		view:   getConfig(srvAddresses),
		msgIds: 0,
	}
}

// Writes the provided value to a random server
func (sc *StorageClient) WriteValue(value string) error {
	sc.msgIds++
	_, err := sc.view.Broadcast(context.Background(), &pb.State{
		Id:        sc.msgIds,
		Value:     value,
		Timestamp: time.Now().Unix(),
	})
	return err
}
