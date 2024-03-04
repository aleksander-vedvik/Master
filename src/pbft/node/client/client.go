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
		view: getConfig(addr, srvAddresses),
		id:   0,
		addr: "127.0.0.1:8080",
	}
}

func (sc *StorageClient) handleResponses(resps []*pb.ClientResponse) {
	res := strconv.Itoa(len(resps)) + " CLIENT RETURN:"
	for _, resp := range resps {
		if resp != nil {
			res += " " + resp.Result
		} else {
			res += " nil"
		}
	}
	log.Println(res)
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
	log.Println("\treceived a response at client:", resp.Result)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
