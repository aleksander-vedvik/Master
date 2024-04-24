package bench

import (
	"context"
	"errors"
	"strconv"

	paxosClient "github.com/aleksander-vedvik/benchmark/paxos/client"
	paxosServer "github.com/aleksander-vedvik/benchmark/paxos/server"
)

type PaxosBenchmark struct{}

func (p PaxosBenchmark) CreateServer(addr string, srvAddrs []string) (*paxosServer.Server, func(), error) {
	srv := paxosServer.New(addr, srvAddrs, true)
	srv.Start()
	return srv, func() {
		srv.Stop()
	}, nil
}

func (p PaxosBenchmark) CreateClient(addr string, srvAddrs []string, _ int) (*paxosClient.Client, func(), error) {
	qSize := 1 + len(srvAddrs)/2
	c := paxosClient.New(addr, srvAddrs, qSize)
	return c, func() {
		c.Stop()
	}, nil
}

func (p PaxosBenchmark) Warmup(client *paxosClient.Client) {
	client.Write("warmup")
}

func (p PaxosBenchmark) StartBenchmark(config *paxosClient.Client) {
	config.Benchmark()
}

func (p PaxosBenchmark) StopBenchmark(config *paxosClient.Client) []Result {
	res, err := config.Benchmark()
	if err != nil {
		return nil
	}
	return res.Metrics
}

func (p PaxosBenchmark) Run(client *paxosClient.Client, ctx context.Context, val int) error {
	resp, err := client.Write(strconv.Itoa(val))
	if err != nil {
		return err
	}
	if resp.GetError() {
		return errors.New("not successful")
	}
	return nil
}
