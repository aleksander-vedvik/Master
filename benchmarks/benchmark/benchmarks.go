package bench

import "fmt"

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

type initializable interface {
	Init(RunOptions)
}

type benchStruct struct {
	run  func(benchmarkOption, any) (ClientResult, []Result, error)
	init func() initializable
}

var benchTypes = map[string]benchStruct{
	PaxosBroadcastCall: {
		run: func(opts benchmarkOption, bench any) (ClientResult, []Result, error) {
			return runBenchmark(opts, bench.(*PaxosBenchmark))
		},
		init: func() initializable {
			return &PaxosBenchmark{}
		},
	},
	PaxosQuorumCall: {
		run: func(opts benchmarkOption, bench any) (ClientResult, []Result, error) {
			return runBenchmark(opts, bench.(*PaxosQCBenchmark))
		},
		init: func() initializable {
			return &PaxosQCBenchmark{}
		},
	},
	PaxosQuorumCallBroadcastOption: {
		run: func(opts benchmarkOption, bench any) (ClientResult, []Result, error) {
			return runBenchmark(opts, bench.(*PaxosQCBBenchmark))
		},
		init: func() initializable {
			return &PaxosQCBBenchmark{}
		},
	},
	PBFTWithGorums: {
		run: func(opts benchmarkOption, bench any) (ClientResult, []Result, error) {
			return runBenchmark(opts, bench.(*PbftBenchmark))
		},
		init: func() initializable {
			return &PbftBenchmark{}
		},
	},
	PBFTWithoutGorums: {
		run: func(opts benchmarkOption, bench any) (ClientResult, []Result, error) {
			return runBenchmark(opts, bench.(*PbftSBenchmark))
		},
		init: func() initializable {
			return &PbftSBenchmark{}
		},
	},
	PBFTNoOrder: {
		run: func(opts benchmarkOption, bench any) (ClientResult, []Result, error) {
			return runBenchmark(opts, bench.(*PbftBenchmark))
		},
		init: func() initializable {
			return &PbftBenchmark{}
		},
	},
	PBFTOrder: {
		run: func(opts benchmarkOption, bench any) (ClientResult, []Result, error) {
			return runBenchmark(opts, bench.(*PbftOBenchmark))
		},
		init: func() initializable {
			return &PbftOBenchmark{}
		},
	},
	Simple: {
		run: func(opts benchmarkOption, bench any) (ClientResult, []Result, error) {
			return runBenchmark(opts, bench.(*SimpleBenchmark))
		},
		init: func() initializable {
			return &SimpleBenchmark{}
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

func createClients[S, C any](bench Benchmark[S, C], opts RunOptions) {
	fmt.Print("creating clients")
	for i := 0; i < opts.numClients; i++ {
		addr := fmt.Sprintf("127.0.0.1:%v", opts.clientBasePort+i)
		if opts.clients != nil {
			addr = opts.clients[i]
		}
		bench.AddClient(i, addr, opts.srvAddrs, opts.logger)
	}
}

func warmupFunc[C any](clients []*C, warmup func(*C)) {
	fmt.Print(": warming up")
	warmupChan := make(chan struct{}, len(clients))
	for _, client := range clients {
		go func(client *C) {
			warmup(client)
			warmupChan <- struct{}{}
		}(client)
	}
	for i := 0; i < len(clients); i++ {
		if i%2 == 0 {
			fmt.Print(".")
		}
		<-warmupChan
	}
	fmt.Println()
	fmt.Println()
}
