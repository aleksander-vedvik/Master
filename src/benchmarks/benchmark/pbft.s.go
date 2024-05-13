package bench

import (
	"context"
	"fmt"
	"strconv"

	pbftClient "github.com/aleksander-vedvik/benchmark/pbft.s/client"
	pbftServer "github.com/aleksander-vedvik/benchmark/pbft.s/server"
)

type PbftSBenchmark struct{}

func (PbftSBenchmark) CreateServer(addr string, srvAddrs []string) (*pbftServer.Server, func(), error) {
	srv := pbftServer.New(addr, srvAddrs, true)
	srv.Start()
	return srv, func() {
		srv.Stop()
	}, nil
}

func (PbftSBenchmark) CreateClient(addr string, srvAddrs []string, _ int) (*pbftClient.Client, func(), error) {
	qSize := 1 + len(srvAddrs)/2
	c := pbftClient.New(addr, srvAddrs, qSize)
	return c, func() {
		c.Stop()
	}, nil
}

func (PbftSBenchmark) Warmup(client *pbftClient.Client) {
	_, err := client.WriteVal(context.Background(), "warmup")
	if err != nil {
		panic(err)
	}
}

func (PbftSBenchmark) StartBenchmark(config *pbftClient.Client) {
}

func (PbftSBenchmark) StopBenchmark(config *pbftClient.Client) []Result {
	return nil
}

func (PbftSBenchmark) Run(client *pbftClient.Client, ctx context.Context, val int) error {
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
