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

type PbftBenchmark struct{}

func (PbftBenchmark) CreateServer(addr string, srvAddrs []string) (*pbftServer.Server, func(), error) {
	srv := pbftServer.New(addr, srvAddrs, nil)
	srv.Start()
	return srv, func() {
		srv.Stop()
	}, nil
}

func (PbftBenchmark) CreateClient(id int, addr string, srvAddrs []string, _ int, logger *slog.Logger) (*pbftClient.Client, func(), error) {
	qSize := 2 * len(srvAddrs) / 3
	c := pbftClient.New(id, addr, srvAddrs[0:1], qSize, logger)
	return c, func() {
		c.Stop()
	}, nil
}

func (PbftBenchmark) Warmup(client *pbftClient.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, _ = client.Write(ctx, "warmup")
	//if err != nil {
	//panic(err)
	//}
}

func (PbftBenchmark) StartBenchmark(config *pbftClient.Client) []Result {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	config.Benchmark(ctx)
	return nil
}

func (PbftBenchmark) StopBenchmark(config *pbftClient.Client) []Result {
	return nil
}

func (PbftBenchmark) Run(client *pbftClient.Client, ctx context.Context, val int) error {
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
