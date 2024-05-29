package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	bench "github.com/aleksander-vedvik/benchmark/benchmark"
	paxosBroadcastCall "github.com/aleksander-vedvik/benchmark/paxos.b/server"
	paxosQuorumCall "github.com/aleksander-vedvik/benchmark/paxosqc/server"
	paxosQuorumCallBroadcastOption "github.com/aleksander-vedvik/benchmark/paxosqcb/server"
	pbftOrder "github.com/aleksander-vedvik/benchmark/pbft.o/server"
	pbftWithoutGorums "github.com/aleksander-vedvik/benchmark/pbft.s/server"
	pbftWithGorums "github.com/aleksander-vedvik/benchmark/pbft/server"
	simple "github.com/aleksander-vedvik/benchmark/simple/server"
	"gopkg.in/yaml.v3"
)

//import (
//"flag"
//"fmt"
//"paxos/client"
//"paxos/server"
//"time"
//)

//var srvAddrs = map[int]string{
//0: "127.0.0.1:5000",
//1: "127.0.0.1:5001",
//2: "127.0.0.1:5002",
//3: "127.0.0.1:5003",
//4: "127.0.0.1:5004",
//}

/*func main() {
	//bench.RunSingleBenchmark("Simple")
	//bench.RunSingleBenchmark("Paxos.BroadcastCall")
	//bench.RunSingleBenchmark("Paxos.QuorumCall")
	//bench.RunSingleBenchmark("Paxos.QuorumCallBroadcastOption")
	//bench.RunSingleBenchmark("PBFT.With.Gorums")
	//bench.RunSingleBenchmark("PBFT.Without.Gorums")
	//bench.RunThroughputVsLatencyBenchmark("Simple", 15000, 15000)
	//bench.RunThroughputVsLatencyBenchmark("Paxos.BroadcastCall", 1000, 1000)
	//bench.RunThroughputVsLatencyBenchmark("Paxos.QuorumCall", 6000, 600)
	bench.RunThroughputVsLatencyBenchmark("Paxos.QuorumCallBroadcastOption", 6000, 600)
	//bench.RunThroughputVsLatencyBenchmark("PBFT.With.Gorums", 5000, 500)
	//bench.RunThroughputVsLatencyBenchmark("PBFT.Without.Gorums", 5000, 500)
	//bench.RunThroughputVsLatencyBenchmark("PBFT.NoOrder")
	//bench.RunThroughputVsLatencyBenchmark("PBFT.Order")
}*/

type Server struct {
	ID   int    `yaml:"id"`
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
}

// ServerEntry represents an entry in the servers list
type ServerEntry map[int]Server

// Config represents the configuration containing servers
type Config struct {
	Servers []ServerEntry `yaml:"servers"`
}

func getConfig() ServerEntry {
	data, err := os.ReadFile("conf.yaml")
	if err != nil {
		panic(err)
	}
	var c Config
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		panic(err)
	}
	srvs := make(map[int]Server, len(c.Servers))
	for _, srv := range c.Servers {
		for id, info := range srv {
			srvs[id] = info
		}
	}
	return srvs
}

type mappingType map[int]string

var mapping mappingType = map[int]string{
	0: bench.PaxosBroadcastCall,
	1: bench.PaxosQuorumCall,
	2: bench.PaxosQuorumCallBroadcastOption,
	3: bench.PBFTWithGorums,
	4: bench.PBFTWithoutGorums,
	5: bench.PBFTNoOrder,
	6: bench.PBFTOrder,
	7: bench.Simple,
}

func (m mappingType) String() string {
	ret := "\n"
	ret += "\t0: " + m[0] + "\n"
	ret += "\t1: " + m[1] + "\n"
	ret += "\t2: " + m[2] + "\n"
	ret += "\t3: " + m[3] + "\n"
	ret += "\t4: " + m[4] + "\n"
	ret += "\t5: " + m[5] + "\n"
	ret += "\t6: " + m[6] + "\n"
	ret += "\t7: " + m[7] + "\n"
	return ret
}

func main() {
	servers := getConfig()
	for id, srv := range servers {
		fmt.Printf("Server %v --> ID: %d, Address: %s, Port: %s\n", id, srv.ID, srv.Addr, srv.Port)
	}
	id := flag.Int("id", 0, "nodeID")
	runSrv := flag.Bool("server", false, "default: false")
	benchTypeIndex := flag.Int("run", 0, "type of benchmark to run"+mapping.String())
	memProfile := flag.Bool("mem", false, "create memory and cpu profile")
	numClients := flag.Int("clients", 0, "number of clients to run")
	clientBasePort := flag.Int("port", 0, "which base port to use for clients")
	withLogger := flag.Bool("log", false, "run with structured logger. Default: false")
	local := flag.Bool("local", true, "run servers locally. Default: true")
	throughput := flag.Int("throughput", 0, "target throughput")
	runs := flag.Int("runs", 0, "number of runs of each benchmark")
	steps := flag.Int("steps", 0, "number of increments on throughput vs latency benchmarks")
	dur := flag.Int("dur", 0, "number of seconds throughput vs latency benchmarks should run")
	broadcastID := flag.Uint64("BroadcastID", 0, "read BroadcastID from log")
	flag.Parse()

	if *broadcastID > 0 {
		readLog(*broadcastID)
		return
	}

	benchType, ok := mapping[*benchTypeIndex]
	if !ok {
		panic("invalid bench type")
	}

	if *runSrv {
		runServer(benchType, *id, servers, *withLogger)
	} else {
		runBenchmark(benchType, *throughput, *numClients, *clientBasePort, *steps, *runs, *dur, *local, servers, *memProfile)
	}
}

func runBenchmark(name string, throughput, numClients, clientBasePort, steps, runs, dur int, local bool, srvAddrs map[int]Server, memProfile bool) {
	var srvAddresses []string
	options := make([]bench.RunOption, 0)
	if !local {
		options = append(options, bench.RunExternal())
		if srvAddrs == nil {
			panic("srvAddrs cannot be nil when not running locally")
		}
		srvAddresses = make([]string, len(srvAddrs))
		for _, srv := range srvAddrs {
			srvAddresses[srv.ID] = fmt.Sprintf("%s:%s", srv.Addr, srv.Port)
		}
	}
	if numClients > 0 {
		options = append(options, bench.NumClients(numClients))
	}
	if clientBasePort > 0 {
		options = append(options, bench.ClientBasePort(clientBasePort))
	}
	if throughput > 0 {
		options = append(options, bench.MaxThroughput(throughput))
	}
	if steps > 0 {
		options = append(options, bench.Steps(steps))
	}
	if runs > 0 {
		options = append(options, bench.Runs(runs))
	}
	if dur > 0 {
		options = append(options, bench.Dur(dur))
	}
	if memProfile {
		options = append(options, bench.WithMemProfile())
	}
	bench.RunThroughputVsLatencyBenchmark(name, options...)
}

type BenchmarkServer interface {
	Start()
	Stop()
}

func runServer(benchType string, id int, srvAddrs map[int]Server, withLogger bool) {
	fmt.Println("Running server:", benchType)
	var logger *slog.Logger
	if withLogger {
		file, err := os.Create(fmt.Sprintf("log.%s:%s.json", srvAddrs[id].Addr, srvAddrs[id].Port))
		if err != nil {
			panic(err)
		}
		loggerOpts := &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}
		handler := slog.NewJSONHandler(file, loggerOpts)
		//handler := slog.NewJSONHandler(os.Stdout, loggerOpts)
		logger = slog.New(handler)
	}
	srvAddresses := make([]string, len(srvAddrs))
	for _, srv := range srvAddrs {
		srvAddresses[srv.ID] = fmt.Sprintf("%s:%s", srv.Addr, srv.Port)
	}
	var srv BenchmarkServer
	switch benchType {
	case bench.PaxosBroadcastCall:
		srv = paxosBroadcastCall.New(srvAddresses[id], srvAddresses, logger)
	case bench.PaxosQuorumCall:
		srv = paxosQuorumCall.New(srvAddresses[id], srvAddresses, logger)
	case bench.PaxosQuorumCallBroadcastOption:
		srv = paxosQuorumCallBroadcastOption.New(srvAddresses[id], srvAddresses, logger)
	case bench.PBFTWithGorums:
		srv = pbftWithGorums.New(srvAddresses[id], srvAddresses, logger)
	case bench.PBFTWithoutGorums:
		srv = pbftWithoutGorums.New(srvAddresses[id], srvAddresses)
	case bench.PBFTNoOrder:
		srv = pbftWithGorums.New(srvAddresses[id], srvAddresses, logger)
	case bench.PBFTOrder:
		srv = pbftOrder.New(srvAddresses[id], srvAddresses, logger)
	case bench.Simple:
		srv = simple.New(srvAddresses[id], srvAddresses, logger)
	}
	srv.Start()
	fmt.Println("Press any key to stop server")
	fmt.Scanln()
	srv.Stop()
}

type logEntry struct {
	Time   time.Time `json:"time"`
	Level  string    `json:"level"`
	Source struct {
		Function string `json:"function"`
		File     string `json:"file"`
		Line     int    `json:"line"`
	} `json:"source"`
	Msg         string `json:"msg"`
	BroadcastID uint64 `json:"BroadcastID"`
	Err         error  `json:"err"`
	Method      string `json:"method"`
	From        string `json:"from"`
}

func readLog(broadcastID uint64) {
	logFiles := []string{"log.127.0.0.1:5000.json", "log.127.0.0.1:5001.json", "log.127.0.0.1:5002.json"}
	for _, logFile := range logFiles {
		fmt.Println()
		fmt.Println("=============")
		fmt.Println("Reading:", logFile)
		fmt.Println()
		file, err := os.Open(logFile)
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(file)
		// optionally, resize scanner's capacity for lines over 64K, see next example
		for scanner.Scan() {
			var entry logEntry
			json.Unmarshal(scanner.Bytes(), &entry)
			if entry.BroadcastID == broadcastID {
				fmt.Println("msg:", entry.Msg, "err:", entry.Err)
			}
		}
	}
}
