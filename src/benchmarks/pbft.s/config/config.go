package config

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/aleksander-vedvik/benchmark/pbft.s/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	mut         sync.Mutex
	nodes       []pb.PBFTNodeClient
	nodeAddrs   map[int]string
	connections []*grpc.ClientConn
	methods     map[string][]string
	who         string
}

func NewConfig(addr string, srvAddrs []string) *Config {
	config := &Config{
		who:         addr,
		nodes:       make([]pb.PBFTNodeClient, len(srvAddrs)),
		connections: make([]*grpc.ClientConn, 0, len(srvAddrs)),
		nodeAddrs:   make(map[int]string, len(srvAddrs)),
		methods:     make(map[string][]string),
	}
	for i, addr := range srvAddrs {
		cc, err := grpc.DialContext(context.Background(), addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(fmt.Sprintf("%s: %s", config.who, err))
		}
		config.connections = append(config.connections, cc)
		config.nodes[i] = pb.NewPBFTNodeClient(cc)
		config.nodeAddrs[i] = addr
	}
	return config
}

func (c *Config) Close() {
	for _, cc := range c.connections {
		cc.Close()
	}
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
		go func(node pb.PBFTNodeClient, j int) {
			_, err := node.Write(context.Background(), req)
			if err != nil {
				//panic(fmt.Sprintf("%s: %s", c.who, err))
			}
			respChan <- struct{}{}
			//slog.Info("sent msg", "to", c.nodeAddrs[i])
		}(node, i)
		sentMsgs++
	}
	for ; sentMsgs > 0; sentMsgs-- {
		<-respChan
	}
	//slog.Info("returning")
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
	//respChan := make(chan string, len(c.nodes))
	//sentMsgs := 0
	for _, node := range c.nodes {
		//slog.Info("server: sending preprepare", "to", c.nodeAddrs[i], "from", c.who)
		go func(node pb.PBFTNodeClient) {
			//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			//defer cancel()
			_, err := node.PrePrepare(context.Background(), req)
			if err != nil {
				//panic(fmt.Sprintf("%s: %s", c.who, err))
			}
			//respChan <- fmt.Sprintf("server: sent preprepare to: %s from: %s", c.nodeAddrs[i], c.who)
		}(node)
		//sentMsgs++
	}
	//for ; sentMsgs > 0; sentMsgs-- {
	//fmt.Println(<-respChan)
	//<-respChan
	//}
	//slog.Warn("server: returning preprepare", "who", c.who)
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
	//respChan := make(chan struct{}, len(c.nodes))
	//sentMsgs := 0
	for i, node := range c.nodes {
		//slog.Info("server: sending prepare", "to", c.nodeAddrs[i], "from", c.who)
		go func(node pb.PBFTNodeClient, to string) {
			//success := false
			for r := 0; r < 10; r++ {
				//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				//defer cancel()
				_, err := node.Prepare(context.Background(), req)
				if err != nil {
					//slog.Error("config:", "err", err, "who", c.who)
				} else {
					//success = true
					break
				}
				time.Sleep(10 * time.Duration(r) * time.Millisecond)
			}
			//if !success {
			//slog.Error("config: not a success", "to", to, "from", c.who)
			//panic(fmt.Sprintf("erh: %s", c.who))
			//}
			//respChan <- struct{}{}
		}(node, c.nodeAddrs[i])
		//sentMsgs++
	}
	//for ; sentMsgs > 0; sentMsgs-- {
	//<-respChan
	//}
	//slog.Warn("server: returning prepare", "who", c.who)
}

func (c *Config) Commit(req *pb.CommitRequest) {
	//slog.Info("config: sending commit")
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
	//respChan := make(chan struct{}, len(c.nodes))
	//sentMsgs := 0
	for _, node := range c.nodes {
		go func(node pb.PBFTNodeClient) {
			//success := false
			for r := 0; r < 10; r++ {
				//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				//defer cancel()
				_, err := node.Commit(context.Background(), req)
				if err != nil {
					//slog.Error("config:", "err", err, "who", c.who)
				} else {
					//success = true
					break
				}
				time.Sleep(10 * time.Millisecond)
			}
			//if !success {
			////panic(fmt.Sprintf("erh: %s", c.who))
			//}
			//respChan <- struct{}{}
		}(node)
		//sentMsgs++
	}
	//for ; sentMsgs > 0; sentMsgs-- {
	//<-respChan
	//}
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
	//slog.Info("from", "addr", req.From)
	cc, err := grpc.DialContext(context.Background(), req.From, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer cc.Close()
	client := pb.NewPBFTNodeClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = client.ClientHandler(ctx, req)
	//if err != nil {
	//panic(err)
	//}
}
