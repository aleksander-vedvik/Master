package storage

import (
	"context"
	"log"
	"time"

	pb "github.com/aleksander-vedvik/Master/protos"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StorageClient struct {
	quorum *pb.Configuration
}

// Creates a new StorageClient with the provided srvAddresses as the configuration
func NewStorageClient(srvAddresses []string) *StorageClient {
	mgr := pb.NewManager(
		gorums.WithDialTimeout(50*time.Millisecond),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	/*quorum, err := mgr.NewConfiguration(
		NewQSpec(len(srvAddresses)),
		gorums.WithNodeList(srvAddresses),
	)*/
	quorum, err := mgr.NewConfiguration2(func(opts *pb.GorumsOptions) {
		opts.AddQuorumSpec(NewQSpec(len(srvAddresses)))
		opts.AddServers(srvAddresses)
		opts.AddServer("localhost:5000")
	})
	if err != nil {
		log.Println("error creating config:", err)
		//return nil
	}
	return &StorageClient{
		quorum: quorum,
	}
}

// Writes the provided value to a random server
func (sc *StorageClient) WriteValue(value string) error {
	_, err := sc.quorum.Write(context.Background(), &pb.State{
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
	reply, err := sc.quorum.Read(ctx, &pb.ReadRequest{})
	defer cancel()
	if err != nil {
		log.Fatalln("read rpc returned error:", err)
		return "", nil
	}
	return reply.Value, nil
}
