package node

import (
	"context"
	pb "pbft/protos"
)

func (s *StorageServer) ConvertPrePrepareToPrepareRequest(ctx context.Context, request *pb.PrePrepareRequest) (response *pb.PrepareRequest) {
	return &pb.PrepareRequest{
		Value: request.GetValue(),
	}
}

func (s *StorageServer) ConvertPrepareToCommitRequest(ctx context.Context, request *pb.PrepareRequest) (response *pb.CommitRequest) {
	return &pb.CommitRequest{
		Value: request.GetValue(),
	}
}
