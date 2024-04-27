package bench

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

type ClientResult struct {
	Id         string
	Total      uint64
	Successes  uint64
	Errs       uint64
	LatencyAvg time.Duration
	LatencyMin time.Duration
	LatencyMax time.Duration
}

type Result struct {
	Id                    string
	TotalNum              uint64
	GoroutinesStarted     uint64
	GoroutinesStopped     uint64
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

type Benchmark[S, C any] interface {
	CreateServer(addr string, peers []string) (*S, func(), error)
	CreateClient(addr string, srvAddrs []string, qSize int) (*C, func(), error)
	Warmup(client *C)
	StartBenchmark(config *C)
	Run(client *C, ctx context.Context, payload int) error
	StopBenchmark(config *C) []Result
}

type benchmarkOption struct {
	name           string
	srvAddrs       []string
	numClients     int
	clientBasePort int
	quorumSize     int
	numRequests    int
	timeout        time.Duration
	async          bool
	local          bool
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
		if i >= 3 {
			return nil, nil
		}
		//results[i], errs[i] = benchmark.run(bench)
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
	var start runtime.MemStats
	var end runtime.MemStats
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

	resChan := make(chan struct {
		dur time.Duration
		err error
	}, opts.numClients*opts.numRequests)
	fmt.Println("starting benchmark...")
	// start the recording of metrics
	benchmark.StartBenchmark(config)
	runtime.ReadMemStats(&start)
	for i := 0; i < opts.numRequests; i++ {
		for _, client := range clients {
			go func(client *C, i int) {
				ctx, cancel := context.WithTimeout(context.Background(), opts.timeout)
				defer cancel()
				start := time.Now()
				err := benchmark.Run(client, ctx, i)
				end := time.Since(start)
				resChan <- struct {
					dur time.Duration
					err error
				}{end, err}
			}(client, i)
			//durations[i], errs[i] = timer(benchmark.Run, client, ctx, i)
		}
	}
	maxDur := time.Duration(0)
	minDur := 100 * time.Hour
	avgDur := time.Duration(0)
	for i := 0; i < totalNumReqs; i++ {
		if i%(totalNumReqs/10) == 0 {
			fmt.Printf("%v%s done\n", (100 * float64(i) / float64(totalNumReqs)), "%")
		}
		res := <-resChan
		durations[i], errs[i] = res.dur, res.err
		if res.err != nil {
			clientResult.Errs++
		} else {
			clientResult.Successes++
		}
		if res.dur > maxDur {
			maxDur = res.dur
		}
		if res.dur < minDur {
			minDur = res.dur
		}
		avgDur += res.dur
	}
	fmt.Printf("100%s done\n", "%")
	runtime.ReadMemStats(&end)
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

//func timer[C any](f func(client *C, ctx context.Context, payload int) error, client *C, ctx context.Context, payload int) (time.Duration, error) {
//start := time.Now()
//err := f(client, ctx, payload)
//return time.Since(start), err
//}

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
