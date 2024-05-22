package bench

var benchTypes = map[string]struct {
	run func(benchmarkOption) (ClientResult, []Result, error)
}{
	"Paxos.BroadcastCall": {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, PaxosBenchmark{})
		},
	},
	"Paxos.QuorumCall": {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, PaxosQCBenchmark{})
		},
	},
	"PBFT.With.Gorums": {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, PbftBenchmark{})
		},
	},
	"PBFT.Without.Gorums": {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, PbftSBenchmark{})
		},
	},
	"PBFT.NoOrder": {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, PbftBenchmark{})
		},
	},
	"PBFT.Order": {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, PbftOBenchmark{})
		},
	},
	"Simple": {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, SimpleBenchmark{})
		},
	},
}

var threeServers = []string{
	"127.0.0.1:5000",
	"127.0.0.1:5001",
	"127.0.0.1:5002",
}

var benchmarks = []benchmarkOption{
	{
		srvAddrs:       threeServers,
		numClients:     1,
		clientBasePort: 8080,
		numRequests:    10000,
		local:          true,
		runType:        Async,
	},
	{
		srvAddrs:       threeServers,
		numClients:     1,
		clientBasePort: 8080,
		numRequests:    10000,
		local:          true,
		runType:        Sync,
	},
	{
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
		srvAddrs:       threeServers,
		numClients:     10,
		clientBasePort: 8080,
		numRequests:    1000,
		local:          true,
		runType:        Async,
	},
	{
		srvAddrs:       threeServers,
		numClients:     10,
		clientBasePort: 8080,
		numRequests:    1000,
		local:          true,
		runType:        Sync,
	},
	{
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
		srvAddrs:       threeServers,
		numClients:     100,
		clientBasePort: 8080,
		numRequests:    1000,
		local:          true,
		runType:        Async,
	},
	{
		srvAddrs:       threeServers,
		numClients:     100,
		clientBasePort: 8080,
		numRequests:    1000,
		local:          true,
		runType:        Sync,
	},
	{
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

var throughputBenchmarks = []benchmarkOption{
	{
		srvAddrs:       threeServers,
		numClients:     10,
		clientBasePort: 8080,
		numRequests:    10,
		local:          true,
		runType:        Async,
	},
	{
		srvAddrs:       threeServers,
		numClients:     10,
		clientBasePort: 8080,
		numRequests:    20,
		local:          true,
		runType:        Async,
	},
	{
		srvAddrs:       threeServers,
		numClients:     10,
		clientBasePort: 8080,
		numRequests:    50,
		local:          true,
		runType:        Async,
	},
	{
		srvAddrs:       threeServers,
		numClients:     10,
		clientBasePort: 8080,
		numRequests:    100,
		local:          true,
		runType:        Async,
	},
	{
		srvAddrs:       threeServers,
		numClients:     10,
		clientBasePort: 8080,
		numRequests:    200,
		local:          true,
		runType:        Async,
	},
	{
		srvAddrs:       threeServers,
		numClients:     10,
		clientBasePort: 8080,
		numRequests:    500,
		local:          true,
		runType:        Async,
	},
	{
		srvAddrs:       threeServers,
		numClients:     10,
		clientBasePort: 8080,
		numRequests:    1000,
		local:          true,
		runType:        Async,
	},
	{
		srvAddrs:       threeServers,
		numClients:     10,
		clientBasePort: 8080,
		numRequests:    5000,
		local:          true,
		runType:        Async,
	},
	{
		srvAddrs:       threeServers,
		numClients:     10,
		clientBasePort: 8080,
		numRequests:    10000,
		local:          true,
		runType:        Async,
	},
}
