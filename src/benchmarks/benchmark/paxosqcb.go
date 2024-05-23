package bench

import (
	"context"
	"strconv"
	"time"

	paxosClient "github.com/aleksander-vedvik/benchmark/paxosqcb/client"
	paxosServer "github.com/aleksander-vedvik/benchmark/paxosqcb/server"
)

type PaxosQCBBenchmark struct{}

func (PaxosQCBBenchmark) CreateServer(addr string, srvAddrs []string) (*paxosServer.PaxosReplica, func(), error) {
	srv := paxosServer.New(addr, srvAddrs)
	srv.Start()
	return srv, func() {
		srv.Stop()
	}, nil
}

func (PaxosQCBBenchmark) CreateClient(id int, addr string, srvAddrs []string, _ int) (*paxosClient.Client, func(), error) {
	qSize := 1 + len(srvAddrs)/2
	c := paxosClient.New(id, addr, srvAddrs, qSize)
	return c, func() {
		c.Stop()
	}, nil
}

func (PaxosQCBBenchmark) Warmup(client *paxosClient.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client.Write(ctx, "warmup")
}

func (PaxosQCBBenchmark) StartBenchmark(config *paxosClient.Client) []Result {
	//config.Benchmark()
	return nil
}

func (PaxosQCBBenchmark) StopBenchmark(config *paxosClient.Client) []Result {
	//res, err := config.Benchmark()
	//if err != nil {
	//	return nil
	//}
	//result := make([]Result, len(res.Metrics))
	//for i, r := range res.Metrics {
	//	result[i] = Result{
	//		TotalNum:              r.TotalNum,
	//		FinishedReqsTotal:     r.FinishedReqsTotal,
	//		FinishedReqsSuccesful: r.FinishedReqsSuccesful,
	//		FinishedReqsFailed:    r.FinishedReqsFailed,
	//		Processed:             r.Processed,
	//		Dropped:               r.Dropped,
	//		Invalid:               r.Invalid,
	//		AlreadyProcessed:      r.AlreadyProcessed,
	//		RoundTripLatencyAvg:   time.Duration(r.RoundTripLatency.Avg),
	//		RoundTripLatencyMin:   time.Duration(r.RoundTripLatency.Min),
	//		RoundTripLatencyMax:   time.Duration(r.RoundTripLatency.Max),
	//		ReqLatencyAvg:         time.Duration(r.ReqLatency.Avg),
	//		ReqLatencyMin:         time.Duration(r.ReqLatency.Min),
	//		ReqLatencyMax:         time.Duration(r.ReqLatency.Max),
	//		ShardDistribution:     r.ShardDistribution,
	//	}
	//}
	//return result
	return nil
}

func (PaxosQCBBenchmark) Run(client *paxosClient.Client, ctx context.Context, val int) error {
	_, err := client.Write(ctx, strconv.Itoa(val))
	if err != nil {
		return err
	}
	return nil
}
