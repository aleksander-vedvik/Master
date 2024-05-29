package bench

import (
	"context"
	"fmt"
	"strconv"

	pbftClient "github.com/aleksander-vedvik/benchmark/pbft.o/client"
	pbftServer "github.com/aleksander-vedvik/benchmark/pbft.o/server"
)

type PbftOBenchmark struct{}

func (PbftOBenchmark) CreateServer(addr string, srvAddrs []string) (*pbftServer.Server, func(), error) {
	srv := pbftServer.New(addr, srvAddrs, nil)
	srv.Start()
	return srv, func() {
		srv.Stop()
	}, nil
}

func (PbftOBenchmark) CreateClient(id int, addr string, srvAddrs []string, _ int) (*pbftClient.Client, func(), error) {
	qSize := 2 * len(srvAddrs) / 3
	c := pbftClient.New(id, addr, srvAddrs[0:1], qSize)
	return c, func() {
		c.Stop()
	}, nil
}

func (PbftOBenchmark) Warmup(client *pbftClient.Client) {
	_, err := client.Write(context.Background(), "warmup")
	if err != nil {
		panic(err)
	}
}

func (PbftOBenchmark) StartBenchmark(config *pbftClient.Client) []Result {
	return nil
}

func (PbftOBenchmark) StopBenchmark(config *pbftClient.Client) []Result {
	return nil
}

func (PbftOBenchmark) Run(client *pbftClient.Client, ctx context.Context, val int) error {
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
