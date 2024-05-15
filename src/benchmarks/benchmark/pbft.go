package bench

import (
	"context"
	"fmt"
	"strconv"

	pbftClient "github.com/aleksander-vedvik/benchmark/pbft/client"
	pbftServer "github.com/aleksander-vedvik/benchmark/pbft/server"
)

type PbftBenchmark struct{}

func (PbftBenchmark) CreateServer(addr string, srvAddrs []string) (*pbftServer.Server, func(), error) {
	srv := pbftServer.New(addr, srvAddrs, true)
	srv.Start()
	return srv, func() {
		srv.Stop()
	}, nil
}

func (PbftBenchmark) CreateClient(addr string, srvAddrs []string, _ int) (*pbftClient.Client, func(), error) {
	qSize := 1 + len(srvAddrs)/2
	c := pbftClient.New(addr, srvAddrs, qSize)
	return c, func() {
		c.Stop()
	}, nil
}

func (PbftBenchmark) Warmup(client *pbftClient.Client) {
	_, err := client.Write("warmup")
	if err != nil {
		panic(err)
	}
}

func (PbftBenchmark) StartBenchmark(config *pbftClient.Client) []Result {
	return nil
}

func (PbftBenchmark) StopBenchmark(config *pbftClient.Client) []Result {
	return nil
}

func (PbftBenchmark) Run(client *pbftClient.Client, ctx context.Context, val int) error {
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
