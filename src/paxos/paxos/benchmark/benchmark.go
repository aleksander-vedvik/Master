package main

import (
	"context"
	"paxos/benchmark/metrics"
	"runtime"
)

type benchmark struct {
	numServers  int
	numClients  int
	numRequests int
	async       bool
}

var benchmarks = []benchmark{
	{
		numServers:  3,
		numClients:  1,
		numRequests: 1000,
		async:       false,
	},
	{
		numServers:  3,
		numClients:  1,
		numRequests: 1000,
		async:       true,
	},
	{
		numServers:  5,
		numClients:  1,
		numRequests: 1000,
		async:       false,
	},
	{
		numServers:  5,
		numClients:  1,
		numRequests: 1000,
		async:       true,
	},
	{
		numServers:  7,
		numClients:  1,
		numRequests: 1000,
		async:       false,
	},
	{
		numServers:  7,
		numClients:  1,
		numRequests: 1000,
		async:       true,
	},
}

func runBenchmark[T any](opts benchmark, benchmarkSetup func() *metrics.Metrics, benchmarkFunc func(ctx context.Context, req *T) error, createPayload func(val int) *T) *metrics.Metrics {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var start runtime.MemStats
	var end runtime.MemStats
	errs := make([]error, 0, opts.numRequests)

	// start the recording of metrics
	benchmarkSetup()

	runtime.ReadMemStats(&start)
	for i := 0; i < opts.numRequests; i++ {
		errs[i] = benchmarkFunc(ctx, createPayload(i))
	}
	runtime.ReadMemStats(&end)

	//clientAllocs := (end.Mallocs - start.Mallocs) / resp.TotalOps
	//clientMem := (end.TotalAlloc - start.TotalAlloc) / resp.TotalOps

	//resp.AllocsPerOp = clientAllocs
	//resp.MemPerOp = clientMem
	//return resp, nil

	// stop the recording and return the metrics
	return benchmarkSetup()
}

/*func runServerBenchmark(opts Options, cfg *Configuration, f serverFunc) (*Result, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	payload := make([]byte, opts.Payload)
	var start runtime.MemStats
	var end runtime.MemStats

	benchmarkFunc := func(stopTime time.Time) {
		for !time.Now().After(stopTime) {
			msg := &TimedMsg{SendTime: time.Now().UnixNano(), Payload: payload}
			f(ctx, msg)
		}
	}

	warmupEnd := time.Now().Add(opts.Warmup)
	for n := 0; n < opts.Concurrent; n++ {
		go benchmarkFunc(warmupEnd)
	}
	err := g.Wait()
	if err != nil {
		return nil, err
	}

	_, err = cfg.StartServerBenchmark(ctx, &StartRequest{})
	if err != nil {
		return nil, err
	}

	runtime.ReadMemStats(&start)
	endTime := time.Now().Add(opts.Duration)
	for n := 0; n < opts.Concurrent; n++ {
		benchmarkFunc(endTime)
	}
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	runtime.ReadMemStats(&end)

	resp, err := cfg.StopServerBenchmark(ctx, &StopRequest{})
	if err != nil {
		return nil, err
	}

	clientAllocs := (end.Mallocs - start.Mallocs) / resp.TotalOps
	clientMem := (end.TotalAlloc - start.TotalAlloc) / resp.TotalOps

	resp.AllocsPerOp = clientAllocs
	resp.MemPerOp = clientMem
	return resp, nil
}
*/
