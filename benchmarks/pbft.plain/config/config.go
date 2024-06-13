package config

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	pb "github.com/aleksander-vedvik/benchmark/pbft.plain/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Node struct {
	addr   string
	client pb.PBFTNodeClient
	conn   *grpc.ClientConn
}

func (n *Node) Close() {
	n.conn.Close()
}

type Config struct {
	mut     sync.Mutex
	nodes   []*Node
	methods map[string][]string
	who     string
	clients map[string]*Node
}

func NewConfig(addr string, srvAddrs []string) *Config {
	config := &Config{
		who:     addr,
		nodes:   make([]*Node, 0, len(srvAddrs)),
		methods: make(map[string][]string),
		clients: make(map[string]*Node),
	}
	go func() {
		// delayed connection to prevent connection refused.
		// this can happen if some of the servers are not
		// online yet.
		time.Sleep(5 * time.Second)
		for _, srvAddr := range srvAddrs {
			if srvAddr == addr {
				// do not add yourself in the configuration.
				// PBFT only sends messages to other nodes,
				// not itself.
				continue
			}
			cc, err := grpc.DialContext(context.Background(), srvAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				panic(err)
			}
			node := &Node{
				addr:   srvAddr,
				client: pb.NewPBFTNodeClient(cc),
				conn:   cc,
			}
			config.nodes = append(config.nodes, node)
		}
	}()
	return config
}

func (c *Config) Close() {
	for _, n := range c.nodes {
		n.Close()
	}
	for _, c := range c.clients {
		c.Close()
	}
}

func (c *Config) NumNodes() int {
	return len(c.nodes)
}

func (c *Config) Write(ctx context.Context, req *pb.WriteRequest) (*pb.ClientResponse, error) {
	c.mut.Lock()
	defer c.mut.Unlock()
	respChan := make(chan struct{}, len(c.nodes))
	sentMsgs := 0
	for i, node := range c.nodes {
		if i > 0 {
			break
		}
		go func(node pb.PBFTNodeClient, j int) {
			_, err := node.Write(context.Background(), req)
			if err != nil {
				panic(fmt.Sprintf("%s: %s", c.who, err))
			}
			respChan <- struct{}{}
		}(node.client, i)
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
	for _, node := range c.nodes {
		go func(node pb.PBFTNodeClient) {
			sendFn := func() error {
				_, err := node.PrePrepare(context.Background(), req)
				return err
			}
			err := sendWithRetries(sendFn, 10)
			if err != nil {
				panic(err)
			}
		}(node.client)
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
	for _, node := range c.nodes {
		go func(node *Node) {
			sendFn := func() error {
				_, err := node.client.Prepare(context.Background(), req)
				return err
			}
			err := sendWithRetries(sendFn, 10)
			if err != nil {
				panic(err)
			}
		}(node)
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
	for _, node := range c.nodes {
		go func(node pb.PBFTNodeClient) {
			sendFn := func() error {
				_, err := node.Commit(context.Background(), req)
				return err
			}
			err := sendWithRetries(sendFn, 10)
			if err != nil {
				panic(err)
			}
		}(node.client)
	}
}

func (c *Config) ClientHandler(req *pb.ClientResponse) {
	c.mut.Lock()
	methodName := "ClientHandler"
	var methods []string
	if m, ok := c.methods[req.Id]; ok {
		for _, method := range m {
			if method == methodName {
				c.mut.Unlock()
				return
			}
		}
		methods = append(m, methodName)
	} else {
		methods = []string{methodName}
	}
	c.methods[req.Id] = methods
	var (
		cc   *grpc.ClientConn
		ok   bool
		err  error
		node *Node
	)

	if node, ok = c.clients[req.From]; !ok {
		cc, err = grpc.DialContext(context.Background(), req.From, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		client := &Node{
			addr:   req.From,
			client: pb.NewPBFTNodeClient(cc),
			conn:   cc,
		}
		c.clients[req.From] = client
		node = client
	}
	c.mut.Unlock()
	sendFn := func() error {
		_, err := node.client.ClientHandler(context.Background(), req)
		return err
	}
	sendWithRetries(sendFn, 10)
	if err != nil {
		panic(err)
	}
}

func (c *Config) Benchmark(ctx context.Context) (*pb.ClientResponse, error) {
	slog.Info("client: sending benchmark")
	c.mut.Lock()
	defer c.mut.Unlock()
	respChan := make(chan struct{}, len(c.nodes))
	sentMsgs := 0
	for i, node := range c.nodes {
		go func(node pb.PBFTNodeClient, j int) {
			sendFn := func() error {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				_, err := node.Benchmark(ctx, &emptypb.Empty{})
				return err
			}
			err := sendWithRetries(sendFn, 10)
			if err != nil {
				panic(err)
			}
			respChan <- struct{}{}
		}(node.client, i)
		sentMsgs++
	}
	for ; sentMsgs > 0; sentMsgs-- {
		<-respChan
	}
	return nil, nil
}

func sendWithRetries(sendFn func() error, retries int) error {
	var err error
	for r := 0; r < retries; r++ {
		err = sendFn()
		if err == nil {
			return nil
		}
		time.Sleep((1 + time.Duration(r)) * 10 * time.Millisecond)
	}
	return err
}
