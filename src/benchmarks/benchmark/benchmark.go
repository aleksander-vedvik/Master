package bench

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

type Result interface {
	GetThroughput() uint64
	GetAvgLatency() uint64
}

type Benchmark[S, C any] interface {
	CreateServer(addr string, peers []string) (*S, func(), error)
	CreateClient(addr string, srvAddrs []string, qSize int) (*C, func(), error)
	Warmup(client *C)
	StartBenchmark(config *C)
	Run(client *C, ctx context.Context, payload int) error
	StopBenchmark(config *C) Result
}

type benchmarkOption struct {
	name           string
	srvAddrs       []string
	numClients     int
	clientBasePort int
	quorumSize     int
	numRequests    int
	async          bool
	local          bool
}

var threeServers = []string{
	"127.0.0.1:5000",
	"127.0.0.1:5001",
	"127.0.0.1:5002",
}

var fiveServers = []string{
	"127.0.0.1:5000",
	"127.0.0.1:5001",
	"127.0.0.1:5002",
	"127.0.0.1:5003",
	"127.0.0.1:5004",
}

var sevenServers = []string{
	"127.0.0.1:5000",
	"127.0.0.1:5001",
	"127.0.0.1:5002",
	"127.0.0.1:5003",
	"127.0.0.1:5004",
	"127.0.0.1:5005",
	"127.0.0.1:5006",
}

var benchmarks = []benchmarkOption{
	{
		name:           "first",
		srvAddrs:       threeServers,
		numClients:     1,
		clientBasePort: 8080,
		numRequests:    1,
		async:          false,
		local:          true,
	},
	{
		name:           "second",
		srvAddrs:       threeServers,
		numClients:     1,
		clientBasePort: 8080,
		numRequests:    1000,
		async:          true,
		local:          true,
	},
	{
		name:           "third",
		srvAddrs:       threeServers,
		numClients:     100,
		clientBasePort: 8080,
		numRequests:    1000,
		async:          true,
		local:          true,
	},
	{
		srvAddrs:       fiveServers,
		numClients:     1,
		clientBasePort: 8080,
		numRequests:    1000,
		async:          false,
		local:          true,
	},
	{
		srvAddrs:       fiveServers,
		numClients:     1,
		clientBasePort: 8080,
		numRequests:    1000,
		async:          true,
		local:          true,
	},
	{
		srvAddrs:       fiveServers,
		numClients:     100,
		clientBasePort: 8080,
		numRequests:    1000,
		async:          true,
		local:          true,
	},
	{
		srvAddrs:       sevenServers,
		numClients:     1,
		clientBasePort: 8080,
		numRequests:    1000,
		async:          false,
		local:          true,
	},
	{
		srvAddrs:       sevenServers,
		numClients:     1,
		clientBasePort: 8080,
		numRequests:    1000,
		async:          true,
		local:          true,
	},
	{
		srvAddrs:       sevenServers,
		numClients:     100,
		clientBasePort: 8080,
		numRequests:    1000,
		async:          true,
		local:          true,
	},
}

func RunAll() {
	for bench := range benchTypes {
		RunSingleBenchmark(bench)
	}
}

func RunSingleBenchmark(name string) ([]Result, []error) {
	benchmark, ok := benchTypes[name]
	if !ok {
		return nil, nil
	}
	fmt.Println("running benchmark:", name)
	results := make([]Result, len(benchmarks))
	errs := make([]error, len(benchmarks))
	for i, bench := range benchmarks {
		fmt.Println("running:", bench.name)
		results[i], errs[i] = benchmark.run(bench)
		break
	}
	return results, errs
}

func runBenchmark[S, C any](opts benchmarkOption, benchmark Benchmark[S, C]) (Result, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var config *C
	var start runtime.MemStats
	var end runtime.MemStats
	errs := make([]error, opts.numRequests)
	durations := make([]time.Duration, opts.numRequests)

	if opts.quorumSize <= 0 {
		opts.quorumSize = len(opts.srvAddrs)
	}

	fmt.Println("creating clients...")
	clients := make([]*C, opts.numClients)
	for i := 0; i < opts.numClients; i++ {
		var (
			err     error
			cleanup func()
		)
		clients[i], cleanup, err = benchmark.CreateClient(fmt.Sprintf("127.0.0.1:%v", opts.clientBasePort+i), opts.srvAddrs, opts.quorumSize)
		defer cleanup()
		if err != nil {
			return nil, err
		}
		if i == 0 {
			config = clients[0]
		}
	}

	fmt.Println("creating servers...")
	var servers []*S
	if opts.local {
		servers = make([]*S, len(opts.srvAddrs))
		for i, addr := range opts.srvAddrs {
			var (
				err     error
				cleanup func()
			)
			servers[i], cleanup, err = benchmark.CreateServer(addr, opts.srvAddrs)
			defer cleanup()
			if err != nil {
				return nil, err
			}
		}
	}

	fmt.Println("warming up...")
	for _, client := range clients {
		benchmark.Warmup(client)
	}

	fmt.Println("starting benchmark...")
	// start the recording of metrics
	benchmark.StartBenchmark(config)
	runtime.ReadMemStats(&start)
	for i := 0; i < opts.numRequests; i++ {
		for _, client := range clients {
			durations[i], errs[i] = timer(benchmark.Run, client, ctx, i)
		}
	}
	// stop the recording and return the metrics
	runtime.ReadMemStats(&end)
	result := benchmark.StopBenchmark(config)
	fmt.Println("stopped benchmark...")

	//clientAllocs := (end.Mallocs - start.Mallocs) / resp.TotalOps
	//clientMem := (end.TotalAlloc - start.TotalAlloc) / resp.TotalOps

	//resp.AllocsPerOp = clientAllocs
	//resp.MemPerOp = clientMem
	//return resp, nil
	return result, nil
}

func timer[C any](f func(client *C, ctx context.Context, payload int) error, client *C, ctx context.Context, payload int) (time.Duration, error) {
	start := time.Now()
	err := f(client, ctx, payload)
	return time.Since(start), err
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
