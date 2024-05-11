package bench

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

type ClientResult struct {
	Id              string
	Total           uint64
	Successes       uint64
	Errs            uint64
	LatencyAvg      time.Duration
	LatencyMin      time.Duration
	LatencyMax      time.Duration
	TotalDur        time.Duration
	ReqDistribution []uint64 // shows how many request in transit per unit of time
}

type Result struct {
	Id                    string
	TotalNum              uint64
	FinishedReqsTotal     uint64
	FinishedReqsSuccesful uint64
	FinishedReqsFailed    uint64
	Processed             uint64
	Dropped               uint64
	Invalid               uint64
	AlreadyProcessed      uint64
	RoundTripLatencyAvg   time.Duration
	RoundTripLatencyMin   time.Duration
	RoundTripLatencyMax   time.Duration
	ReqLatencyAvg         time.Duration
	ReqLatencyMin         time.Duration
	ReqLatencyMax         time.Duration
	ShardDistribution     map[uint32]uint64
}

type RequestResult struct {
	err   error
	start time.Time
	end   time.Time
}

type Benchmark[S, C any] interface {
	CreateServer(addr string, peers []string) (*S, func(), error)
	CreateClient(addr string, srvAddrs []string, qSize int) (*C, func(), error)
	Warmup(client *C)
	StartBenchmark(config *C)
	Run(client *C, ctx context.Context, payload int) error
	StopBenchmark(config *C) []Result
}

const (
	Sync = iota
	Async
	Random
)

type benchmarkOption struct {
	name           string
	srvAddrs       []string
	numClients     int
	clientBasePort int
	quorumSize     int
	numRequests    int
	timeout        time.Duration
	reqInterval    struct{ start, end int } // reqs will be sent in the interval: [start, end] µs
	local          bool
	runType        int
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
		if i < 3 || i >= 6 {
			continue
		}
		start := time.Now()
		clientResult, ress, err := benchmark.run(bench)
		if err != nil {
			panic(err)
		}
		fmt.Println("took:", time.Since(start))
		err = WriteToCsv(fmt.Sprintf("./csv/%s.%s.csv", name, bench.name), ress, clientResult)
		if err != nil {
			panic(err)
		}
	}
	return results, errs
}

func runBenchmark[S, C any](opts benchmarkOption, benchmark Benchmark[S, C]) (ClientResult, []Result, error) {
	fmt.Printf("\nRunning benchmark: %s\n\n", opts.name)
	runtime.GC()
	var config *C
	//var start runtime.MemStats
	//var end runtime.MemStats
	totalNumReqs := opts.numClients * opts.numRequests
	clientResult := ClientResult{
		Id:    "clients",
		Total: uint64(totalNumReqs),
	}
	errs := make([]error, opts.numRequests*opts.numClients)
	durations := make([]time.Duration, opts.numRequests*opts.numClients)

	if opts.quorumSize <= 0 {
		opts.quorumSize = len(opts.srvAddrs)
	}
	if opts.timeout <= 0 {
		opts.timeout = 5 * time.Second
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
			return clientResult, nil, err
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
				return clientResult, nil, err
			}
		}
	}

	fmt.Println("warming up...")
	for _, client := range clients {
		benchmark.Warmup(client)
	}

	resChan := make(chan RequestResult, opts.numClients*opts.numRequests)
	fmt.Println("starting benchmark...")

	// start the recording of metrics
	benchmark.StartBenchmark(config)
	//runtime.ReadMemStats(&start)
	//for i := 0; i < opts.numRequests; i++ {
	//for _, client := range clients {
	//go func(client *C, i int) {
	//ctx, cancel := context.WithTimeout(context.Background(), opts.timeout)
	//defer cancel()
	//start := time.Now()
	//err := benchmark.Run(client, ctx, i)
	//end := time.Now()
	//resChan <- RequestResult{
	//err:   err,
	//start: start,
	//end:   end,
	//}
	//}(client, i)
	//}
	//}
	runStart := time.Now()
	switch opts.runType {
	case Sync:
		runSync(opts, benchmark, resChan, clients)
	case Async:
		runAsync(opts, benchmark, resChan, clients)
	case Random:
		runRandom(opts, benchmark, resChan, clients)
	}

	maxDur := time.Duration(0)
	minDur := 100 * time.Hour
	avgDur := time.Duration(0)
	for i := 0; i < totalNumReqs; i++ {
		if i%(totalNumReqs/10) == 0 {
			fmt.Printf("%v%s done\n", (100 * float64(i) / float64(totalNumReqs)), "%")
		}
		var res RequestResult
		// prevent deadlock
		select {
		case res = <-resChan:
		case <-time.After(30 * time.Second):
			return clientResult, nil, errors.New("could not collect all responses")
		}

		dur := res.end.Sub(res.start)
		durations[i], errs[i] = dur, res.err
		if res.err != nil {
			clientResult.Errs++
		} else {
			clientResult.Successes++
		}
		if dur > maxDur {
			maxDur = dur
		}
		if dur < minDur {
			minDur = dur
		}
		avgDur += dur
	}
	clientResult.TotalDur = time.Since(runStart)
	fmt.Printf("100%s done\n", "%")
	//runtime.ReadMemStats(&end)
	// stop the recording and return the metrics
	result := benchmark.StopBenchmark(config)
	fmt.Println("stopped benchmark...")

	//clientAllocs := (end.Mallocs - start.Mallocs) / resp.TotalOps
	//clientMem := (end.TotalAlloc - start.TotalAlloc) / resp.TotalOps

	//resp.AllocsPerOp = clientAllocs
	//resp.MemPerOp = clientMem
	//return resp, nil

	avgDur /= time.Duration(len(durations))
	clientResult.LatencyAvg = avgDur
	clientResult.LatencyMax = maxDur
	clientResult.LatencyMin = minDur
	return clientResult, result, nil
}

// each client runs synchronously. I.e. waits for response before sending next msg.
func runSync[S, C any](opts benchmarkOption, benchmark Benchmark[S, C], resChan chan RequestResult, clients []*C) {
	for _, client := range clients {
		go func(client *C) {
			for i := 0; i < opts.numRequests; i++ {
				ctx, cancel := context.WithTimeout(context.Background(), opts.timeout)
				defer cancel()
				start := time.Now()
				err := benchmark.Run(client, ctx, i)
				end := time.Now()
				resChan <- RequestResult{
					err:   err,
					start: start,
					end:   end,
				}
			}
		}(client)
	}
}

// each client runs asynchronously. I.e. sends all requests at once.
func runAsync[S, C any](opts benchmarkOption, benchmark Benchmark[S, C], resChan chan RequestResult, clients []*C) {
	for i := 0; i < opts.numRequests; i++ {
		for _, client := range clients {
			go func(client *C, i int) {
				ctx, cancel := context.WithTimeout(context.Background(), opts.timeout)
				defer cancel()
				start := time.Now()
				err := benchmark.Run(client, ctx, i)
				end := time.Now()
				resChan <- RequestResult{
					err:   err,
					start: start,
					end:   end,
				}
			}(client, i)
		}
	}
}

// each client runs synchronously but with an added delay between requests.
func runRandom[S, C any](opts benchmarkOption, benchmark Benchmark[S, C], resChan chan RequestResult, clients []*C) {
	for _, client := range clients {
		go func(client *C) {
			for i := 0; i < opts.numRequests; i++ {
				ctx, cancel := context.WithTimeout(context.Background(), opts.timeout)
				defer cancel()
				start := time.Now()
				err := benchmark.Run(client, ctx, i)
				end := time.Now()
				resChan <- RequestResult{
					err:   err,
					start: start,
					end:   end,
				}
				t := rand.Intn(opts.reqInterval.end-opts.reqInterval.start) + opts.reqInterval.start
				time.Sleep(time.Duration(t) * time.Microsecond)
			}
		}(client)
	}
}
