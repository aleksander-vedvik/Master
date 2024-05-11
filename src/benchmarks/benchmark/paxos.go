package bench

import (
	"context"
	"errors"
	"strconv"
	"time"

	paxosClient "github.com/aleksander-vedvik/benchmark/paxos.b/client"
	paxosServer "github.com/aleksander-vedvik/benchmark/paxos.b/server"
)

type PaxosBenchmark struct{}

func (PaxosBenchmark) CreateServer(addr string, srvAddrs []string) (*paxosServer.Server, func(), error) {
	srv := paxosServer.New(addr, srvAddrs, true)
	srv.Start()
	return srv, func() {
		srv.Stop()
	}, nil
}

func (PaxosBenchmark) CreateClient(addr string, srvAddrs []string, _ int) (*paxosClient.Client, func(), error) {
	qSize := 1 + len(srvAddrs)/2
	c := paxosClient.New(addr, srvAddrs, qSize)
	return c, func() {
		c.Stop()
	}, nil
}

func (PaxosBenchmark) Warmup(client *paxosClient.Client) {
	client.Write("warmup")
}

func (PaxosBenchmark) StartBenchmark(config *paxosClient.Client) {
	config.Benchmark()
}

func (PaxosBenchmark) StopBenchmark(config *paxosClient.Client) []Result {
	res, err := config.Benchmark()
	if err != nil {
		return nil
	}
	result := make([]Result, len(res.Metrics))
	for i, r := range res.Metrics {
		result[i] = Result{
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

func (PaxosBenchmark) Run(client *paxosClient.Client, ctx context.Context, val int) error {
	resp, err := client.Write(strconv.Itoa(val))
	if err != nil {
		return err
	}
	if resp.GetError() {
		return errors.New("not successful")
	}
	return nil
}
