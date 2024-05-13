package bench

var benchTypes = map[string]struct {
	run func(benchmarkOption) (ClientResult, []Result, error)
}{
	"Paxos": {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, PaxosBenchmark{})
		},
	},
	"PaxosQC": {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, PaxosQCBenchmark{})
		},
	},
	"PBFT": {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, PbftBenchmark{})
		},
	},
	"PBFT.S": {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, PbftBenchmark{})
		},
	},
	"Simple": {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, SimpleBenchmark{})
		},
	},
}

var benchmarks = []benchmarkOption{
	{
		name:           "S3.C1.R10000.Async",
		srvAddrs:       threeServers,
		numClients:     1,
		clientBasePort: 8080,
		numRequests:    10000,
		local:          true,
		runType:        Async,
	},
	{
		name:           "S3.C1.R10000.Sync",
		srvAddrs:       threeServers,
		numClients:     1,
		clientBasePort: 8080,
		numRequests:    10000,
		local:          true,
		runType:        Sync,
	},
	{
		name:           "S3.C1.R10000.Random",
		srvAddrs:       threeServers,
		numClients:     1,
		clientBasePort: 8080,
		numRequests:    10000,
		local:          true,
		runType:        Random,
		reqInterval: struct {
			start int
			end   int
		}{50, 400},
	},
	{
		name:           "S3.C10.R1000.Async",
		srvAddrs:       threeServers,
		numClients:     10,
		clientBasePort: 8080,
		numRequests:    1000,
		local:          true,
		runType:        Async,
	},
	{
		name:           "S3.C10.R1000.Sync",
		srvAddrs:       threeServers,
		numClients:     10,
		clientBasePort: 8080,
		numRequests:    1000,
		local:          true,
		runType:        Sync,
	},
	{
		name:           "S3.C10.R1000.Random",
		srvAddrs:       threeServers,
		numClients:     10,
		clientBasePort: 8080,
		numRequests:    1000,
		local:          true,
		runType:        Random,
		reqInterval: struct {
			start int
			end   int
		}{50, 400},
	},
	{
		name:           "S3.C100.R1000.Async",
		srvAddrs:       threeServers,
		numClients:     100,
		clientBasePort: 8080,
		numRequests:    1000,
		local:          true,
		runType:        Async,
	},
	{
		name:           "S3.C100.R1000.Sync",
		srvAddrs:       threeServers,
		numClients:     100,
		clientBasePort: 8080,
		numRequests:    1000,
		local:          true,
		runType:        Sync,
	},
	{
		name:           "S3.C100.R1000.Random",
		srvAddrs:       threeServers,
		numClients:     100,
		clientBasePort: 8080,
		numRequests:    1000,
		local:          true,
		runType:        Random,
		reqInterval: struct {
			start int
			end   int
		}{50, 400},
	},
}

var threeServers = []string{
	"127.0.0.1:5000",
	"127.0.0.1:5001",
	"127.0.0.1:5002",
}
