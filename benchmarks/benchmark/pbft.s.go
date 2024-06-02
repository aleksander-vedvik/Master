package bench

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	pbftClient "github.com/aleksander-vedvik/benchmark/pbft.s/client"
	pbftServer "github.com/aleksander-vedvik/benchmark/pbft.s/server"
)

type PbftSBenchmark struct {
	clients []*pbftClient.Client
}

func (*PbftSBenchmark) CreateServer(addr string, srvAddrs []string) (*pbftServer.Server, func(), error) {
	srv := pbftServer.New(addr, srvAddrs, true)
	srv.Start(true)
	return srv, func() {
		srv.Stop()
	}, nil
}

func (b *PbftSBenchmark) Init(opts RunOptions) {
	b.clients = make([]*pbftClient.Client, 0, len(opts.clients))
	createClients(b, opts)
	warmupFunc(b.clients, b.warmup)
}

func (b *PbftSBenchmark) Clients() []*pbftClient.Client {
	return b.clients
}

func (b *PbftSBenchmark) Config() *pbftClient.Client {
	return b.clients[0]
}

func (b *PbftSBenchmark) Stop() {
	for _, client := range b.clients {
		client.Stop()
	}
}

func (b *PbftSBenchmark) AddClient(id int, addr string, srvAddrs []string, logger *slog.Logger) {
	qSize := 2 * len(srvAddrs) / 3
	c := pbftClient.New(addr, srvAddrs, qSize)
	c.Start()
	b.clients = append(b.clients, c)
}

func (*PbftSBenchmark) warmup(client *pbftClient.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client.WriteVal(ctx, "warmup")
}

func (*PbftSBenchmark) StartBenchmark(config *pbftClient.Client) []Result {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	config.SendBenchmark(ctx)
	return nil
}

func (*PbftSBenchmark) StopBenchmark(config *pbftClient.Client) []Result {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	config.SendBenchmark(ctx)
	return nil
}

func (*PbftSBenchmark) Run(client *pbftClient.Client, ctx context.Context, val int) error {
	value := strconv.Itoa(val)
	resp, err := client.WriteVal(ctx, value)
	if err != nil {
		return err
	}
	if resp.GetResult() != value {
		return fmt.Errorf("wrong result. want: %s, got: %s", value, resp.GetResult())
	}
	return nil
}
