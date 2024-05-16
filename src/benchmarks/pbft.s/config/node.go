package config

import (
	"context"
	"sync"

	pb "github.com/aleksander-vedvik/benchmark/pbft.s/protos"
	"google.golang.org/grpc"
)

type Config struct {
	mut     sync.Mutex
	nodes   []pb.PBFTNodeClient
	methods map[string][]string
}

func NewConfig(srvAddrs []string) *Config {
	config := &Config{
		nodes:   make([]pb.PBFTNodeClient, len(srvAddrs)),
		methods: make(map[string][]string),
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

func (c *Config) NumNodes() int {
	return len(c.nodes)
}

func (c *Config) Write(ctx context.Context, req *pb.WriteRequest) (*pb.ClientResponse, error) {
	c.mut.Lock()
	defer c.mut.Unlock()
	methodName := "Write"
	var methods []string
	if m, ok := c.methods[req.Id]; ok {
		for _, method := range m {
			if method == methodName {
				return nil, nil
			}
		}
		methods = append(m, methodName)
	} else {
		methods = []string{methodName}
	}
	c.methods[req.Id] = methods
	respChan := make(chan struct{}, len(c.nodes))
	sentMsgs := 0
	for i, node := range c.nodes {
		if i > 0 {
			break
		}
		go func(node pb.PBFTNodeClient) {
			node.Write(ctx, req)
			respChan <- struct{}{}
		}(node)
		sentMsgs++
	}
	for ; sentMsgs > 0; sentMsgs-- {
		<-respChan
	}
	return nil, nil
}

func (c *Config) PrePrepare(req *pb.PrePrepareRequest) {
	c.mut.Lock()
	defer c.mut.Unlock()
	methodName := "PrePrepare"
	var methods []string
	if m, ok := c.methods[req.Id]; ok {
		for _, method := range m {
			if method == methodName {
				return
			}
		}
		methods = append(m, methodName)
	} else {
		methods = []string{methodName}
	}
	c.methods[req.Id] = methods
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
	c.mut.Lock()
	defer c.mut.Unlock()
	methodName := "Prepare"
	var methods []string
	if m, ok := c.methods[req.Id]; ok {
		for _, method := range m {
			if method == methodName {
				return
			}
		}
		methods = append(m, methodName)
	} else {
		methods = []string{methodName}
	}
	c.methods[req.Id] = methods
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
	c.mut.Lock()
	defer c.mut.Unlock()
	methodName := "Commit"
	var methods []string
	if m, ok := c.methods[req.Id]; ok {
		for _, method := range m {
			if method == methodName {
				return
			}
		}
		methods = append(m, methodName)
	} else {
		methods = []string{methodName}
	}
	c.methods[req.Id] = methods
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

func (c *Config) ClientHandler(req *pb.ClientResponse) {
	c.mut.Lock()
	defer c.mut.Unlock()
	methodName := "ClientHandler"
	var methods []string
	if m, ok := c.methods[req.Id]; ok {
		for _, method := range m {
			if method == methodName {
				return
			}
		}
		methods = append(m, methodName)
	} else {
		methods = []string{methodName}
	}
	c.methods[req.Id] = methods
	cc, err := grpc.DialContext(context.Background(), req.From)
	if err != nil {
		panic(err)
	}
	client := pb.NewPBFTNodeClient(cc)
	client.ClientHandler(context.Background(), req)
}
