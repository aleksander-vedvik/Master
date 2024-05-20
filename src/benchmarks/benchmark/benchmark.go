package bench

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"runtime"
	"strconv"
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
	Throughput      float64
	ReqDistribution []uint64 // shows how many request in transit per unit of time
	BucketSize      int      // number of microseconds
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
	StartBenchmark(config *C) []Result
	Run(client *C, ctx context.Context, payload int) error
	StopBenchmark(config *C) []Result
}

type runType int

const (
	Sync runType = iota
	Async
	Random
	Throughput
)

func (r runType) String() string {
	res := ""
	switch r {
	case Sync:
		res = "Sync"
	case Async:
		res = "Async"
	case Random:
		res = "Random"
	}
	return res
}

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
	runType        runType
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
	//throughputVsLatency := make([][]string, 0)
	for _, bench := range benchmarks {
		if bench.numRequests > 1000 {
			continue
		}
		bench.name = fmt.Sprintf("S%v.C%v.R%v.%s", len(bench.srvAddrs), bench.numClients, bench.numRequests, bench.runType)
		start := time.Now()
		clientResult, ress, err := benchmark.run(bench)
		fmt.Println("took:", time.Since(start))
		if err != nil {
			panic(err)
		}
		//throughputVsLatency = append(throughputVsLatency, []string{fmt.Sprintf("%.2f", clientResult.Throughput), strconv.Itoa(int(clientResult.LatencyAvg.Microseconds()))})
		fmt.Println("writing to csv...")
		err = WriteToCsv(name, bench.name, ress, clientResult)
		if err != nil {
			panic(err)
		}
		fmt.Println("done")
	}
	//err := WriteThroughputVsLatency(name, throughputVsLatency)
	//if err != nil {
	//panic(err)
	//}
	return results, errs
}

func RunThroughputVsLatencyBenchmark(name string) ([]Result, []error) {
	benchmark, ok := benchTypes[name]
	if !ok {
		return nil, nil
	}
	fmt.Println("running benchmark:", name)
	//throughputs := []int{
	//1000,
	//2000,
	//3000,
	//4000,
	//5000,
	//10000,
	//}
	maxTarget := 15000
	targetIncrement := 1000
	results := make([]Result, len(benchmarks))
	errs := make([]error, len(benchmarks))
	throughputVsLatency := make([][]string, 0)
	for target := targetIncrement; target <= maxTarget; target += targetIncrement {
		//for _, target := range throughputs {
		bench := benchmarkOption{
			name:           fmt.Sprintf("%s.S3.C10.R%v.Throughput", name, target),
			srvAddrs:       threeServers,
			numClients:     10,
			clientBasePort: 8080,
			numRequests:    target,
			local:          true,
			runType:        Throughput,
		}
		start := time.Now()
		clientResult, _, err := benchmark.run(bench)
		throughputVsLatency = append(throughputVsLatency, []string{strconv.Itoa(int(clientResult.Throughput)), strconv.Itoa(int(clientResult.LatencyAvg.Milliseconds()))})
		if err != nil {
			panic(err)
		}
		fmt.Println("took:", time.Since(start))
		fmt.Println("done")
	}
	err := WriteThroughputVsLatency(name, throughputVsLatency)
	if err != nil {
		panic(err)
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
	//errs := make([]error, opts.numRequests*opts.numClients)
	//durations := make([]time.Duration, opts.numRequests*opts.numClients)

	if opts.quorumSize <= 0 {
		opts.quorumSize = len(opts.srvAddrs)
	}
	if opts.timeout <= 0 {
		opts.timeout = 15 * time.Second
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
	resultBefore := benchmark.StartBenchmark(config)
	//runtime.ReadMemStats(&start)

	//runStart := time.Now()
	switch opts.runType {
	case Sync:
		go runSync(opts, benchmark, resChan, clients)
	case Async:
		go runAsync(opts, benchmark, resChan, clients)
	case Random:
		go runRandom(opts, benchmark, resChan, clients)
	case Throughput:
		go runThroughput(opts, benchmark, resChan, clients)
	}

	switch opts.runType {
	case Throughput:
		fmt.Println("collecting responses...")
		avgDur := time.Duration(0)
		numFailed := 0
		for i := 0; i < 10*opts.numRequests; i++ {
			if i%(opts.numRequests/2) == 0 {
				fmt.Printf("%v%s done\n", (100 * float64(i) / float64(10*opts.numRequests)), "%")
			}
			var res RequestResult
			// prevent deadlock
			select {
			case res = <-resChan:
			case <-time.After(2 * time.Minute):
				slog.Info("benchmark:", "replies", i, "total", 10*opts.numRequests, "successes", i-numFailed, "failures", numFailed)
				return clientResult, nil, errors.New("could not collect all responses")
			}
			if res.err != nil {
				numFailed++
				continue
			}
			avgDur += res.end.Sub(res.start)
		}
		fmt.Printf("100%s done, numFailed: %v out of %v\n", "%", numFailed, 10*opts.numRequests)
		benchmark.StopBenchmark(config)
		fmt.Println("stopped benchmark...")

		avgDur /= time.Duration(10 * opts.numRequests)
		clientResult.LatencyAvg = avgDur
		clientResult.Throughput = float64(opts.numRequests)
		fmt.Println("done")
		return clientResult, nil, nil
	default:
	}
	numBuckets := 100
	durDistribution := make([]uint64, numBuckets)
	durations := make([]time.Duration, 0, opts.numRequests*opts.numClients)
	maxDur := time.Duration(0)
	minDur := 100 * time.Hour
	avgDur := time.Duration(0)
	firstReqStart := time.Now().Add(24 * time.Hour)
	lastReqDone := time.Now()
	for i := 0; i < totalNumReqs; i++ {
		if i%(totalNumReqs/10) == 0 {
			fmt.Printf("%v%s done\n", (100 * float64(i) / float64(totalNumReqs)), "%")
		}
		var res RequestResult
		// prevent deadlock
		select {
		case res = <-resChan:
		case <-time.After(30 * time.Second):
			slog.Info("benchmark:", "replies", i, "total", totalNumReqs)
			return clientResult, nil, errors.New("could not collect all responses")
		}

		//durations[i], errs[i] = dur, res.err
		if res.err != nil {
			clientResult.Errs++
			//slog.Info("failed req", "err", res.err)
			continue
		} else {
			clientResult.Successes++
		}
		if firstReqStart.Sub(res.start) > 0 {
			firstReqStart = res.start
		}
		if res.end.Sub(lastReqDone) > 0 {
			lastReqDone = res.end
		}
		dur := res.end.Sub(res.start)
		durations = append(durations, dur)
		if dur > maxDur {
			maxDur = dur
		}
		if dur < minDur {
			minDur = dur
		}
		avgDur += dur
	}
	clientResult.TotalDur = lastReqDone.Sub(firstReqStart)
	clientResult.ReqDistribution = durDistribution
	fmt.Printf("100%s done\n", "%")
	//runtime.ReadMemStats(&end)
	// stop the recording and return the metrics
	result := benchmark.StopBenchmark(config)
	fmt.Println("stopped benchmark...")

	fmt.Println("calculating histogram...")
	bucketSize := (maxDur.Microseconds() - minDur.Microseconds()) / int64(numBuckets)
	clientResult.BucketSize = int(bucketSize)
	for _, dur := range durations {
		bucket := int(math.Floor(float64(dur.Microseconds()-minDur.Microseconds()) / float64(bucketSize)))
		if bucket >= numBuckets {
			bucket = numBuckets - 1
		}
		if bucket <= 0 || bucket >= len(durDistribution) {
			continue
		}
		durDistribution[bucket]++
	}

	//clientAllocs := (end.Mallocs - start.Mallocs) / resp.TotalOps
	//clientMem := (end.TotalAlloc - start.TotalAlloc) / resp.TotalOps

	//resp.AllocsPerOp = clientAllocs
	//resp.MemPerOp = clientMem
	//return resp, nil

	avgDur /= time.Duration(opts.numRequests * opts.numClients)
	clientResult.LatencyAvg = avgDur
	clientResult.LatencyMax = maxDur
	clientResult.LatencyMin = minDur
	if result != nil && resultBefore != nil {
		for i := range result {
			result[i].TotalNum -= resultBefore[i].TotalNum
		}
	}
	totalDur := clientResult.TotalDur.Seconds()
	if totalDur <= 0 {
		totalDur = 1
	}
	clientResult.Throughput = float64(clientResult.Total) / totalDur
	fmt.Println("done")
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

// each client runs asynchronously. I.e. sends all requests at once.
// runs for 10 seconds and sends reqs at target throughput.
func runThroughput[S, C any](opts benchmarkOption, benchmark Benchmark[S, C], resChan chan RequestResult, clients []*C) {
	// opts.numRequests denotes the target throughput in reqs/s. We thus
	// need to divide the number of reqs by the number of clients.
	numReqsPerClient := opts.numRequests / len(clients)
	// it will run for 10 seconds
	for t := 0; t < 10; t++ {
		fmt.Printf("\t%v: sending reqs...\n", t)
		for _, client := range clients {
			for i := 0; i < numReqsPerClient; i++ {
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
		time.Sleep(1 * time.Second)
	}
}
