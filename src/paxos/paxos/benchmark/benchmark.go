package main

import (
	"context"
	"runtime"
	"time"
)

type Result interface {
	GetThroughput() uint64
	GetAvgLatency() time.Duration
}

type Benchmark[S, C, P any] interface {
	CreateServer(addr string, peers []string) (*S, error)
	CreateClient(addr string, srvAddrs []string, qSize int) (*C, error)
	Warmup(config *C)
	StartBenchmark(config *C)
	Run(client *C, ctx context.Context, payload *P) error
	CreatePayload(val int) *P
	StopBenchmark(config *C) Result
}

type PaxosPayload struct{}
type PaxosClient struct{}
type PaxosServer struct{}
type PaxosBenchmark struct{}

func (p PaxosBenchmark) CreateServer(addr string, srvAddrs []string) (*PaxosServer, error) {
	return nil, nil
}

func (p PaxosBenchmark) CreateClient(addr string, srvAddrs []string, qSize int) (*PaxosClient, error) {
	return nil, nil
}

func (p PaxosBenchmark) Warmup(config *PaxosClient) {
}

func (p PaxosBenchmark) StartBenchmark(config *PaxosClient) {
}

func (p PaxosBenchmark) StopBenchmark(config *PaxosClient) Result {
	return nil
}

func (p PaxosBenchmark) Run(client *PaxosClient, ctx context.Context, payload *PaxosPayload) error {
	return nil
}

func (p PaxosBenchmark) CreatePayload(val int) *PaxosPayload {
	return nil
}

type benchmarkOption struct {
	numServers  int
	srvAddrs    []string
	numClients  int
	clientAddrs []string
	numRequests int
	async       bool
	local       bool
}

var benchmarks = []benchmarkOption{
	{
		numServers:  3,
		numClients:  1,
		numRequests: 1000,
		async:       false,
		local:       true,
	},
	{
		numServers:  3,
		numClients:  1,
		numRequests: 1000,
		async:       true,
		local:       true,
	},
	{
		numServers:  5,
		numClients:  1,
		numRequests: 1000,
		async:       false,
		local:       true,
	},
	{
		numServers:  5,
		numClients:  1,
		numRequests: 1000,
		async:       true,
		local:       true,
	},
	{
		numServers:  7,
		numClients:  1,
		numRequests: 1000,
		async:       false,
		local:       true,
	},
	{
		numServers:  7,
		numClients:  1,
		numRequests: 1000,
		async:       true,
		local:       true,
	},
}

var benchTypes = map[string]struct {
	run func(benchmarkOption)
}{
	"Paxos": {
		run: func(bench benchmarkOption) {
			runBenchmark[PaxosServer, PaxosClient, PaxosPayload](bench, PaxosBenchmark{})
		},
	},
}

func runBenchmarks(name string) {
	benchmark, ok := benchTypes[name]
	if !ok {
		return
	}
	for _, bench := range benchmarks {
		benchmark.run(bench)
	}
}

func runBenchmark[S, C, P any](opts benchmarkOption, benchmark Benchmark[S, C, P]) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var start runtime.MemStats
	var end runtime.MemStats
	errs := make([]error, 0, opts.numRequests)

	clients := make([]*C, opts.numClients)
	for i := 0; i < opts.numClients; i++ {
		var err error
		clients[i], err = benchmark.CreateClient("", opts.srvAddrs, len(opts.srvAddrs))
		if err != nil {
			panic(err)
		}
	}

	var servers []*S
	if opts.local {
		servers = make([]*S, len(opts.srvAddrs))
		for i, addr := range opts.srvAddrs {
			var err error
			servers[i], err = benchmark.CreateServer(addr, opts.srvAddrs)
			if err != nil {
				panic(err)
			}
		}
	}

	// start the recording of metrics
	runtime.ReadMemStats(&start)
	for i := 0; i < opts.numRequests; i++ {
		for _, client := range clients {
			errs[i] = benchmark.Run(client, ctx, benchmark.CreatePayload(i))
		}
	}
	for i := 0; i < opts.numRequests; i++ {
		for _, client := range clients {
			errs[i] = benchmark.Run(client, ctx, benchmark.CreatePayload(i))
		}
	}
	runtime.ReadMemStats(&end)

	//clientAllocs := (end.Mallocs - start.Mallocs) / resp.TotalOps
	//clientMem := (end.TotalAlloc - start.TotalAlloc) / resp.TotalOps

	//resp.AllocsPerOp = clientAllocs
	//resp.MemPerOp = clientMem
	//return resp, nil

	// stop the recording and return the metrics
	return
}

func timer() {
	defer time.Since(time.Now())

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
