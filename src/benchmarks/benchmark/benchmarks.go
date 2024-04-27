package bench

var benchTypes = map[string]struct {
	run func(benchmarkOption) (ClientResult, []Result, error)
}{
	"Paxos": {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, PaxosBenchmark{})
		},
	},
	"PBFT": {
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
		name:           "S3.C1.R1000",
		srvAddrs:       threeServers,
		numClients:     1,
		clientBasePort: 8080,
		numRequests:    1000,
		local:          true,
	},
	{
		name:           "S3.C10.R100",
		srvAddrs:       threeServers,
		numClients:     10,
		clientBasePort: 8080,
		numRequests:    100,
		local:          true,
	},
	{
		name:           "S3.C100.R100",
		srvAddrs:       threeServers,
		numClients:     100,
		clientBasePort: 8080,
		numRequests:    100,
		local:          true,
	},
	{
		name:           "S5.C1.R1000",
		srvAddrs:       fiveServers,
		numClients:     1,
		clientBasePort: 8080,
		numRequests:    1000,
		local:          true,
	},
	{
		name:           "S5.C10.R100",
		srvAddrs:       fiveServers,
		numClients:     10,
		clientBasePort: 8080,
		numRequests:    100,
		async:          true,
		local:          true,
	},
	{
		name:           "S5.C100.R100",
		srvAddrs:       fiveServers,
		numClients:     100,
		clientBasePort: 8080,
		numRequests:    100,
		local:          true,
	},
	{
		name:           "S7.C1.R1000",
		srvAddrs:       sevenServers,
		numClients:     1,
		clientBasePort: 8080,
		numRequests:    1000,
		local:          true,
	},
	{
		name:           "S7.C10.R100",
		srvAddrs:       sevenServers,
		numClients:     10,
		clientBasePort: 8080,
		numRequests:    100,
		local:          true,
	},
	{
		name:           "S7.C100.R100",
		srvAddrs:       sevenServers,
		numClients:     100,
		clientBasePort: 8080,
		numRequests:    100,
		local:          true,
	},
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
