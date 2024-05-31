package bench

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	paxosClient "github.com/aleksander-vedvik/benchmark/paxosqc/client"
	paxosServer "github.com/aleksander-vedvik/benchmark/paxosqc/server"
)

type PaxosQCBenchmark struct {
	clients []*paxosClient.Client
}

func (*PaxosQCBenchmark) CreateServer(addr string, srvAddrs []string) (*paxosServer.PaxosReplica, func(), error) {
	srv := paxosServer.New(addr, srvAddrs, nil)
	srv.Start(true)
	return srv, func() {
		srv.Stop()
	}, nil
}

func (b *PaxosQCBenchmark) Init(opts RunOptions) {
	b.clients = make([]*paxosClient.Client, 0, len(opts.clients))
	createClients(b, opts)
	warmupFunc(b.clients, b.warmup)
}

func (b *PaxosQCBenchmark) Clients() []*paxosClient.Client {
	return b.clients
}

func (b *PaxosQCBenchmark) Stop() {
	for _, client := range b.clients {
		client.Stop()
	}
}

func (b *PaxosQCBenchmark) AddClient(id int, addr string, srvAddrs []string, logger *slog.Logger) {
	qSize := 1 + len(srvAddrs)/2
	b.clients = append(b.clients, paxosClient.New(id, addr, srvAddrs, qSize, logger))
}

func (*PaxosQCBenchmark) warmup(client *paxosClient.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client.Write(ctx, "warmup")
}

func (*PaxosQCBenchmark) StartBenchmark(config *paxosClient.Client) []Result {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	config.Benchmark(ctx)
	return nil
}

func (*PaxosQCBenchmark) StopBenchmark(config *paxosClient.Client) []Result {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	config.Benchmark(ctx)
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

func (*PaxosQCBenchmark) Run(client *paxosClient.Client, ctx context.Context, val int) error {
	_, err := client.Write(ctx, strconv.Itoa(val))
	if err != nil {
		return err
	}
	return nil
}
