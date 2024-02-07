package node

import (
	"context"
	"log"
	pb "pbft/protos"
	"strconv"
	"time"
)

type StorageClient struct {
	view *pb.Configuration
	id   int
}

// Creates a new StorageClient with the provided srvAddresses as the configuration
func NewStorageClient(srvAddresses []string, addr string) *StorageClient {
	return &StorageClient{
		view: getConfig(addr, srvAddresses),
		id:   0,
	}
}

func (sc *StorageClient) WriteValue(value string) error {
	sc.id++
	resp, err := sc.view.Write(context.Background(), &pb.WriteRequest{
		Id:        strconv.Itoa(sc.id),
		From:      "client",
		Message:   value,
		Timestamp: time.Now().Unix(),
	})
	log.Println("\treceived a response at client:", resp)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
