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
	addr string
}

// Creates a new StorageClient with the provided srvAddresses as the configuration
func NewStorageClient(srvAddresses []string, addr string) *StorageClient {
	return &StorageClient{
		view: getConfig("127.0.0.1:8080", srvAddresses),
		id:   0,
	}
}

func (sc *StorageClient) WriteValue(value string) error {
	sc.id++
	ctx := context.Background()
	req := &pb.WriteRequest{
		Id:        strconv.Itoa(sc.id),
		From:      "client",
		Message:   value,
		Timestamp: time.Now().Unix(),
	}
	resp, err := sc.view.Write(ctx, req)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("\treceived a response at client:", resp.Result)
	return nil
}
