package client

import (
	"context"
	"time"

	pb "github.com/aleksander-vedvik/Master/protos"
)

type StorageClient struct {
	view   *pb.Configuration
	msgIds int64
}

const (
	BroadcastCall string = "broadcast"
	QuorumCall    string = "quorum"
)

// Creates a new StorageClient with the provided srvAddresses as the configuration
func NewStorageClient(srvAddresses []string, callType string, qSize ...int) *StorageClient {
	var config *pb.Configuration
	if callType == BroadcastCall {
		config = GetBConfig(srvAddresses, qSize[0])
	} else if callType == QuorumCall {
		config = GetQConfig(srvAddresses)
	} else {
		return nil
	}
	return &StorageClient{
		view:   config,
		msgIds: 0,
	}
}

// Writes the provided value to a random server
func (sc *StorageClient) WriteValue(value string) error {
	sc.msgIds++
	//ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//defer cancel()
	ctx := context.Background()
	//ctx := context.WithValue(bg, "test", "test")
	_, err := sc.view.Broadcast(ctx, &pb.State{
		Id:        sc.msgIds,
		Value:     value,
		Timestamp: time.Now().Unix(),
	})
	return err
}

func (sc *StorageClient) CreateStudent(value string) error {
	sc.msgIds++
	ctx := context.Background()
	_, err := sc.view.SaveStudent(ctx, &pb.State{
		Id:        sc.msgIds,
		Value:     value,
		Timestamp: time.Now().Unix(),
	})
	return err
}
