package node

import (
	"context"
	"log"
	pb "pbft/protos"
)

type StorageClient struct {
	view *pb.Configuration
}

// Creates a new StorageClient with the provided srvAddresses as the configuration
func NewStorageClient(srvAddresses []string, addr string) *StorageClient {
	return &StorageClient{
		view: getConfig(addr, srvAddresses),
	}
}

func (sc *StorageClient) WriteValue(value string) error {
	resp, err := sc.view.PrePrepare(context.Background(), &pb.PrePrepareRequest{
		Value: value,
	})
	log.Println(resp)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
