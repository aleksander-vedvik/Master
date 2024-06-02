package bench

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	pbftClient "github.com/aleksander-vedvik/benchmark/pbft/client"
	pbftServer "github.com/aleksander-vedvik/benchmark/pbft/server"
)

type PbftBenchmark struct {
	clients []*pbftClient.Client
}

func (*PbftBenchmark) CreateServer(addr string, srvAddrs []string) (*pbftServer.Server, func(), error) {
	srv := pbftServer.New(addr, srvAddrs, nil)
	srv.Start(true)
	return srv, func() {
		srv.Stop()
	}, nil
}

func (b *PbftBenchmark) Init(opts RunOptions) {
	b.clients = make([]*pbftClient.Client, 0, len(opts.clients))
	createClients(b, opts)
	warmupFunc(b.clients, b.warmup)
}

func (b *PbftBenchmark) Clients() []*pbftClient.Client {
	return b.clients
}

func (b *PbftBenchmark) Config() *pbftClient.Client {
	return b.clients[0]
}

func (b *PbftBenchmark) Stop() {
	for _, client := range b.clients {
		client.Stop()
	}
}

func (b *PbftBenchmark) AddClient(id int, addr string, srvAddrs []string, logger *slog.Logger) {
	srv := 1
	if id == 0 {
		srv = len(srvAddrs)
	}
	qSize := 2 * len(srvAddrs) / 3
	b.clients = append(b.clients, pbftClient.New(id, addr, srvAddrs[0:srv], qSize, logger))
}

func (*PbftBenchmark) warmup(client *pbftClient.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client.Write(ctx, "warmup")
}

func (*PbftBenchmark) StartBenchmark(config *pbftClient.Client) []Result {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	config.Benchmark(ctx)
	return nil
}

func (*PbftBenchmark) StopBenchmark(config *pbftClient.Client) []Result {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	config.Benchmark(ctx)
	return nil
}

func (*PbftBenchmark) Run(client *pbftClient.Client, ctx context.Context, val int) error {
	value := strconv.Itoa(val)
	resp, err := client.Write(ctx, value)
	if err != nil {
		return err
	}
	if resp.GetResult() != value {
		return fmt.Errorf("wrong result. want: %s, got: %s", value, resp.GetResult())
	}
	return nil
}
