package bench

import (
	"context"
	"fmt"
	"strconv"
	"time"

	simpleClient "github.com/aleksander-vedvik/benchmark/simple/client"
	simpleServer "github.com/aleksander-vedvik/benchmark/simple/server"
)

type SimpleBenchmark struct{}

func (SimpleBenchmark) CreateServer(addr string, srvAddrs []string) (*simpleServer.Server, func(), error) {
	srv := simpleServer.New(addr, srvAddrs)
	srv.Start()
	return srv, func() {
		srv.Stop()
	}, nil
}

func (SimpleBenchmark) CreateClient(addr string, srvAddrs []string, qSize int) (*simpleClient.Client, func(), error) {
	c := simpleClient.New(addr, srvAddrs, qSize)
	return c, func() {
		c.Stop()
	}, nil
}

func (SimpleBenchmark) Warmup(client *simpleClient.Client) {
	client.Write1(context.Background(), "warmup")
}

func (SimpleBenchmark) StartBenchmark(config *simpleClient.Client) []Result {
	config.Benchmark()
	return nil
}

func (SimpleBenchmark) StopBenchmark(config *simpleClient.Client) []Result {
	res, err := config.Benchmark()
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
	return result
}

func (SimpleBenchmark) Run(client *simpleClient.Client, ctx context.Context, val int) error {
	resp, err := client.Write1(ctx, strconv.Itoa(val))
	if err != nil {
		return err
	}
	if resp.GetMessage() != strconv.Itoa(val) {
		return fmt.Errorf("wrong result, got: %s, want:%s", resp.GetMessage(), strconv.Itoa(val))
	}
	return nil
}
