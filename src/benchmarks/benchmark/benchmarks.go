package bench

const (
	PaxosBroadcastCall             string = "Paxos.BroadcastCall"
	PaxosQuorumCall                string = "Paxos.QuorumCall"
	PaxosQuorumCallBroadcastOption string = "Paxos.QuorumCallBroadcastOption"
	PBFTWithGorums                 string = "PBFT.With.Gorums"
	PBFTWithoutGorums              string = "PBFT.Without.Gorums"
	PBFTNoOrder                    string = "PBFT.NoOrder"
	PBFTOrder                      string = "PBFT.Order"
	Simple                         string = "Simple"
)

var benchTypes = map[string]struct {
	run func(benchmarkOption) (ClientResult, []Result, error)
}{
	PaxosBroadcastCall: {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, PaxosBenchmark{})
		},
	},
	PaxosQuorumCall: {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, PaxosQCBenchmark{})
		},
	},
	PaxosQuorumCallBroadcastOption: {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, PaxosQCBBenchmark{})
		},
	},
	PBFTWithGorums: {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, PbftBenchmark{})
		},
	},
	PBFTWithoutGorums: {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, PbftSBenchmark{})
		},
	},
	PBFTNoOrder: {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, PbftBenchmark{})
		},
	},
	PBFTOrder: {
		run: func(bench benchmarkOption) (ClientResult, []Result, error) {
			return runBenchmark(bench, PbftOBenchmark{})
		},
	},
	Simple: {
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
	/*{
		srvAddrs:       threeServers,
		numClients:     10,
		clientBasePort: 8080,
		numRequests:    500,
		local:          true,
		runType:        Async,
	},*/
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
