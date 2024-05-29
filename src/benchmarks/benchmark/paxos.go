package bench

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	paxosClient "github.com/aleksander-vedvik/benchmark/paxos.b/client"
	paxosServer "github.com/aleksander-vedvik/benchmark/paxos.b/server"
)

type PaxosBenchmark struct{}

func (PaxosBenchmark) CreateServer(addr string, srvAddrs []string) (*paxosServer.Server, func(), error) {
	srv := paxosServer.New(addr, srvAddrs, nil)
	srv.Start()
	return srv, func() {
		srv.Stop()
	}, nil
}

func (PaxosBenchmark) CreateClient(id int, addr string, srvAddrs []string, _ int, logger *slog.Logger) (*paxosClient.Client, func(), error) {
	qSize := 1 + len(srvAddrs)/2
	c := paxosClient.New(id, addr, srvAddrs, qSize, logger)
	return c, func() {
		c.Stop()
	}, nil
}

func (PaxosBenchmark) Warmup(client *paxosClient.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = client.Write(ctx, "warmup")
}

func (PaxosBenchmark) StartBenchmark(config *paxosClient.Client) []Result {
	res, err := config.Benchmark()
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

func (PaxosBenchmark) StopBenchmark(config *paxosClient.Client) []Result {
	return nil
	/*res, err := config.Benchmark()
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
	return result*/
}

func (PaxosBenchmark) Run(client *paxosClient.Client, ctx context.Context, val int) error {
	resp, err := client.Write(ctx, strconv.Itoa(val))
	if err != nil {
		return err
	}
	if resp.GetError() {
		return errors.New("not successful")
	}
	return nil
}
