package bench

import (
	"context"
	"fmt"
	"strconv"

	pbftClient "github.com/aleksander-vedvik/benchmark/pbft/client"
	pbftServer "github.com/aleksander-vedvik/benchmark/pbft/server"
)

type PbftBenchmark struct{}

func (p PbftBenchmark) CreateServer(addr string, srvAddrs []string) (*pbftServer.Server, func(), error) {
	srv := pbftServer.New(addr, srvAddrs)
	srv.Start()
	return srv, func() {
		srv.Stop()
	}, nil
}

func (p PbftBenchmark) CreateClient(addr string, srvAddrs []string, _ int) (*pbftClient.Client, func(), error) {
	qSize := 1 + len(srvAddrs)/2
	c := pbftClient.New(addr, srvAddrs, qSize)
	return c, func() {
		c.Stop()
	}, nil
}

func (p PbftBenchmark) Warmup(client *pbftClient.Client) {
	client.Write("warmup")
}

func (p PbftBenchmark) StartBenchmark(config *pbftClient.Client) {
}

func (p PbftBenchmark) StopBenchmark(config *pbftClient.Client) Result {
	return nil
}

func (p PbftBenchmark) Run(client *pbftClient.Client, ctx context.Context, val int) error {
	value := strconv.Itoa(val)
	resp, err := client.Write(value)
	if err != nil {
		return err
	}
	if resp.GetResult() != value {
		return fmt.Errorf("wrong result. want: %s, got: %s", value, resp.GetResult())
	}
	return nil
}
