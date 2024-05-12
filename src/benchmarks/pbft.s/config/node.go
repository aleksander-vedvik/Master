package config

import (
	"context"

	pb "github.com/aleksander-vedvik/benchmark/pbft.s/protos"
	"google.golang.org/grpc"
)

type Config struct {
	nodes []pb.PBFTNodeClient
}

func NewConfig(srvAddrs []string) *Config {
	config := &Config{
		nodes: make([]pb.PBFTNodeClient, len(srvAddrs)),
	}
	for i, addr := range srvAddrs {
		cc, err := grpc.DialContext(context.Background(), addr)
		if err != nil {
			panic(err)
		}
		config.nodes[i] = pb.NewPBFTNodeClient(cc)
	}
	return config
}

func (c *Config) Write(ctx context.Context, req *pb.WriteRequest) (*pb.ClientResponse, error) {
	respChan := make(chan struct{}, len(c.nodes))
	for _, node := range c.nodes {
		go func(node pb.PBFTNodeClient) {
			node.Write(ctx, req)
			respChan <- struct{}{}
		}(node)
	}
	for sentMsgs := len(c.nodes); sentMsgs > 0; sentMsgs-- {
		<-respChan
	}
	return nil, nil
}

func (c *Config) PrePrepare(req *pb.PrePrepareRequest) {
	respChan := make(chan struct{}, len(c.nodes))
	for _, node := range c.nodes {
		go func(node pb.PBFTNodeClient) {
			node.PrePrepare(context.Background(), req)
			respChan <- struct{}{}
		}(node)
	}
	for sentMsgs := len(c.nodes); sentMsgs > 0; sentMsgs-- {
		<-respChan
	}
}

func (c *Config) Prepare(req *pb.PrepareRequest) {
	respChan := make(chan struct{}, len(c.nodes))
	for _, node := range c.nodes {
		go func(node pb.PBFTNodeClient) {
			node.Prepare(context.Background(), req)
			respChan <- struct{}{}
		}(node)
	}
	for sentMsgs := len(c.nodes); sentMsgs > 0; sentMsgs-- {
		<-respChan
	}
}

func (c *Config) Commit(req *pb.CommitRequest) {
	respChan := make(chan struct{}, len(c.nodes))
	for _, node := range c.nodes {
		go func(node pb.PBFTNodeClient) {
			node.Commit(context.Background(), req)
			respChan <- struct{}{}
		}(node)
	}
	for sentMsgs := len(c.nodes); sentMsgs > 0; sentMsgs-- {
		<-respChan
	}
}
