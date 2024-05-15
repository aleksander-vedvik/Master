package bench

import (
	"context"
	"strconv"

	paxosClient "github.com/aleksander-vedvik/benchmark/paxosqc/client"
	paxosServer "github.com/aleksander-vedvik/benchmark/paxosqc/server"
)

type PaxosQCBenchmark struct{}

func (PaxosQCBenchmark) CreateServer(addr string, srvAddrs []string) (*paxosServer.PaxosReplica, func(), error) {
	srv := paxosServer.New(addr, srvAddrs)
	srv.Start()
	return srv, func() {
		srv.Stop()
	}, nil
}

func (PaxosQCBenchmark) CreateClient(addr string, srvAddrs []string, _ int) (*paxosClient.Client, func(), error) {
	qSize := 1 + len(srvAddrs)/2
	c := paxosClient.New(addr, srvAddrs, qSize)
	return c, func() {
		c.Stop()
	}, nil
}

func (PaxosQCBenchmark) Warmup(client *paxosClient.Client) {
	client.Write("warmup")
}

func (PaxosQCBenchmark) StartBenchmark(config *paxosClient.Client) []Result {
	//config.Benchmark()
	return nil
}

func (PaxosQCBenchmark) StopBenchmark(config *paxosClient.Client) []Result {
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

func (PaxosQCBenchmark) Run(client *paxosClient.Client, ctx context.Context, val int) error {
	_, err := client.Write(strconv.Itoa(val))
	if err != nil {
		return err
	}
	return nil
}
