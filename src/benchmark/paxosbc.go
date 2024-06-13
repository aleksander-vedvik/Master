package bench

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	paxosClient "github.com/aleksander-vedvik/benchmark/paxos.bc/client"
	paxosServer "github.com/aleksander-vedvik/benchmark/paxos.bc/server"
)

type PaxosBenchmark struct {
	clients []*paxosClient.Client
}

func (b *PaxosBenchmark) Init(opts RunOptions) {
	b.clients = make([]*paxosClient.Client, 0, len(opts.clients))
	createClients(b, opts)
	warmupFunc(b.clients, b.warmup)
}

func (*PaxosBenchmark) CreateServer(addr string, srvAddrs []string) (*paxosServer.Server, func(), error) {
	srv := paxosServer.New(addr, srvAddrs, nil)
	srv.Start(true)
	return srv, func() {
		srv.Stop()
	}, nil
}

func (b *PaxosBenchmark) Clients() []*paxosClient.Client {
	return b.clients
}

func (b *PaxosBenchmark) Config() *paxosClient.Client {
	return b.clients[0]
}

func (b *PaxosBenchmark) Stop() {
	for _, client := range b.clients {
		client.Stop()
	}
}

func (b *PaxosBenchmark) AddClient(id int, addr string, srvAddrs []string, logger *slog.Logger) {
	qSize := 1 + len(srvAddrs)/2
	b.clients = append(b.clients, paxosClient.New(id, addr, srvAddrs, qSize, logger))
}

func (*PaxosBenchmark) warmup(client *paxosClient.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = client.Write(ctx, "warmup")
}

func (*PaxosBenchmark) StartBenchmark(config *paxosClient.Client) []Result {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	res, err := config.Benchmark(ctx)
	if err != nil {
		return nil
	}
	result := make([]Result, len(res.Metrics))
	for i, r := range res.Metrics {
		result[i] = Result{
			TotalNum: r.TotalNum,
			Dropped:  r.Dropped,
		}
	}
	return result
}

func (*PaxosBenchmark) StopBenchmark(config *paxosClient.Client) []Result {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	config.Benchmark(ctx)
	return nil
}

func (*PaxosBenchmark) Run(client *paxosClient.Client, ctx context.Context, val int) error {
	resp, err := client.Write(ctx, strconv.Itoa(val))
	if err != nil {
		return err
	}
	if resp.GetError() {
		return errors.New("not successful")
	}
	return nil
}
