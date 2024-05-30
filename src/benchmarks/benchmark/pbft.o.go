package bench

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	pbftClient "github.com/aleksander-vedvik/benchmark/pbft.o/client"
	pbftServer "github.com/aleksander-vedvik/benchmark/pbft.o/server"
)

type PbftOBenchmark struct {
	clients []*pbftClient.Client
}

func (*PbftOBenchmark) CreateServer(addr string, srvAddrs []string) (*pbftServer.Server, func(), error) {
	srv := pbftServer.New(addr, srvAddrs, nil)
	srv.Start()
	return srv, func() {
		srv.Stop()
	}, nil
}

func (b *PbftOBenchmark) Init(opts RunOptions) {
	b.clients = make([]*pbftClient.Client, len(opts.clients))
	createClients(b, opts)
	warmupFunc(b.clients, b.warmup)
}

func (b *PbftOBenchmark) Clients() []*pbftClient.Client {
	return b.clients
}

func (b *PbftOBenchmark) Stop() {
	for _, client := range b.clients {
		client.Stop()
	}
}

func (b *PbftOBenchmark) AddClient(id int, addr string, srvAddrs []string, logger *slog.Logger) {
	qSize := 2 * len(srvAddrs) / 3
	b.clients = append(b.clients, pbftClient.New(id, addr, srvAddrs[0:1], qSize, logger))
}

func (*PbftOBenchmark) warmup(client *pbftClient.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client.Write(ctx, "warmup")
}

func (*PbftOBenchmark) StartBenchmark(config *pbftClient.Client) []Result {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	config.Benchmark(ctx)
	return nil
}

func (*PbftOBenchmark) StopBenchmark(config *pbftClient.Client) []Result {
	return nil
}

func (*PbftOBenchmark) Run(client *pbftClient.Client, ctx context.Context, val int) error {
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
