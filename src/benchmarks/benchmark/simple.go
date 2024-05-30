package bench

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	simpleClient "github.com/aleksander-vedvik/benchmark/simple/client"
	simpleServer "github.com/aleksander-vedvik/benchmark/simple/server"
)

type SimpleBenchmark struct {
	clients []*simpleClient.Client
}

func (*SimpleBenchmark) CreateServer(addr string, srvAddrs []string) (*simpleServer.Server, func(), error) {
	srv := simpleServer.New(addr, srvAddrs, nil)
	srv.Start()
	return srv, func() {
		srv.Stop()
	}, nil
}

func (b *SimpleBenchmark) Init(opts RunOptions) {
	b.clients = make([]*simpleClient.Client, len(opts.clients))
	createClients(b, opts)
	warmupFunc(b.clients, b.warmup)
}

func (b *SimpleBenchmark) Clients() []*simpleClient.Client {
	return b.clients
}

func (b *SimpleBenchmark) Stop() {
	for _, client := range b.clients {
		client.Stop()
	}
}

func (b *SimpleBenchmark) AddClient(id int, addr string, srvAddrs []string, logger *slog.Logger) {
	qSize := len(srvAddrs)
	b.clients = append(b.clients, simpleClient.New(id, addr, srvAddrs, qSize, logger))
}

func (*SimpleBenchmark) warmup(client *simpleClient.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client.Write1(ctx, "warmup")
}

func (*SimpleBenchmark) StartBenchmark(config *simpleClient.Client) []Result {
	config.Benchmark()
	return nil
}

func (*SimpleBenchmark) StopBenchmark(config *simpleClient.Client) []Result {
	/*res, err := config.Benchmark()
	if err != nil {
		return nil
	}
	result := make([]Result, len(res.Metrics))
	for i, r := range res.Metrics {
		result[i] = Result{
			Id:                    r.Addr,
			TotalNum:              r.TotalNum,
			FinishedReqsTotal:     r.FinishedReqsTotal,
			FinishedReqsSuccesful: r.FinishedReqsSuccesful,
			FinishedReqsFailed:    r.FinishedReqsFailed,
			Processed:             r.Processed,
			Dropped:               r.Dropped,
			Invalid:               r.Invalid,
			AlreadyProcessed:      r.AlreadyProcessed,
			RoundTripLatencyAvg:   time.Duration(r.RoundTripLatency.Avg),
			RoundTripLatencyMin:   time.Duration(r.RoundTripLatency.Min),
			RoundTripLatencyMax:   time.Duration(r.RoundTripLatency.Max),
			ReqLatencyAvg:         time.Duration(r.ReqLatency.Avg),
			ReqLatencyMin:         time.Duration(r.ReqLatency.Min),
			ReqLatencyMax:         time.Duration(r.ReqLatency.Max),
			ShardDistribution:     r.ShardDistribution,
		}
	}
	return result*/
	return nil
}

func (*SimpleBenchmark) Run(client *simpleClient.Client, ctx context.Context, val int) error {
	resp, err := client.Write1(ctx, strconv.Itoa(val))
	if err != nil {
		return err
	}
	if resp.GetMessage() != strconv.Itoa(val) {
		return fmt.Errorf("wrong result, got: %s, want:%s", resp.GetMessage(), strconv.Itoa(val))
	}
	return nil
}
