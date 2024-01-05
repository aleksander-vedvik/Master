package gorums

import (
	"context"
	"log"
	"math/rand"
	"time"

	pb "github.com/aleksander-vedvik/Master/protos"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type StorageClient struct {
	allNodesConfig *pb.Configuration
}

// Creates a new StorageClient with the provided srvAddresses as the configuration
func NewStorageClient(srvAddresses []string) *StorageClient {
	mgr := pb.NewManager(
		gorums.WithDialTimeout(500*time.Millisecond),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	allNodesConfig, err := mgr.NewConfiguration(
		&QSpec{len(srvAddresses)},
		gorums.WithNodeList(srvAddresses),
	)
	if err != nil {
		log.Fatal("error creating read config:", err)
		return nil
	}
	return &StorageClient{
		allNodesConfig: allNodesConfig,
	}
}

// Writes the provided value to a random server
func (sc *StorageClient) WriteValue(value string) error {
	nodes := sc.allNodesConfig.Nodes()
	randomIndex := rand.Intn(len(nodes))
	node := nodes[randomIndex]

	_, err := node.Write(context.Background(), &pb.WriteRequest{
		Value: value,
	})
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// Returns a slice of values stored on all servers
func (sc *StorageClient) ReadValues() ([]string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	reply, err := sc.allNodesConfig.Read(ctx, &emptypb.Empty{})
	cancel()
	if err != nil {
		log.Fatalln("read rpc returned error:", err)
		return []string{}, nil
	}
	return reply.Values, nil
}

type QSpec struct {
	quorumSize int
}

func (qs *QSpec) ReadQF(in *emptypb.Empty, replies map[uint32]*pb.ReadResponse) (*pb.ReadResponse, bool) {
	if len(replies) < qs.quorumSize {
		return nil, false
	}
	return combineResponses(replies), true
}

func combineResponses(replies map[uint32]*pb.ReadResponse) *pb.ReadResponse {
	res := make([]string, 0, len(replies))
	for _, value := range replies {
		res = append(res, value.GetValues()...)
	}
	return &pb.ReadResponse{
		Values: res,
	}
}
