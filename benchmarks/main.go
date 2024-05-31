package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"

	bench "github.com/aleksander-vedvik/benchmark/benchmark"
	paxosBroadcastCall "github.com/aleksander-vedvik/benchmark/paxos.b/server"
	paxosQuorumCall "github.com/aleksander-vedvik/benchmark/paxosqc/server"
	paxosQuorumCallBroadcastOption "github.com/aleksander-vedvik/benchmark/paxosqcb/server"
	pbftOrder "github.com/aleksander-vedvik/benchmark/pbft.o/server"
	pbftWithoutGorums "github.com/aleksander-vedvik/benchmark/pbft.s/server"
	pbftWithGorums "github.com/aleksander-vedvik/benchmark/pbft/server"
	simple "github.com/aleksander-vedvik/benchmark/simple/server"
	"github.com/relab/gorums"
	"gopkg.in/yaml.v3"
)

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
	Clients []ServerEntry `yaml:"clients"`
}

func getConfig() (srvs, clients ServerEntry) {
	data, err := os.ReadFile("conf.yaml")
	if err != nil {
		panic(err)
	}
	var c Config
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		panic(err)
	}
	srvs = make(map[int]Server, len(c.Servers))
	for _, srv := range c.Servers {
		for id, info := range srv {
			srvs[id] = info
		}
	}
	clients = make(map[int]Server, len(c.Servers))
	for _, client := range c.Clients {
		for id, info := range client {
			clients[id] = info
		}
	}
	return srvs, clients
}

type mappingType map[int]string

var mapping mappingType = map[int]string{
	1: bench.PaxosBroadcastCall,
	2: bench.PaxosQuorumCall,
	3: bench.PaxosQuorumCallBroadcastOption,
	4: bench.PBFTWithGorums,
	5: bench.PBFTWithoutGorums,
	6: bench.PBFTNoOrder,
	7: bench.PBFTOrder,
	8: bench.Simple,
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
	fmt.Println("--------")
	fmt.Println("Servers")
	servers, clients := getConfig()
	for id, srv := range servers {
		fmt.Printf("\tServer %v --> ID: %d, Address: %s, Port: %s\n", id, srv.ID, srv.Addr, srv.Port)
	}
	fmt.Println("--------")
	fmt.Println("Clients")
	for id, client := range clients {
		fmt.Printf("\tClient %v --> ID: %d, Address: %s, Port: %s\n", id, client.ID, client.Addr, client.Port)
	}
	fmt.Println()

	id := flag.Int("id", -1, "nodeID")
	runSrv := flag.Bool("server", false, "default: false")
	benchTypeIndex := flag.Int("run", 0, "type of benchmark to run"+mapping.String())
	memProfile := flag.Bool("mem", false, "create memory and cpu profile")
	numClients := flag.Int("clients", 0, "number of clients to run")
	clientBasePort := flag.Int("port", 0, "which base port to use for clients")
	withLogger := flag.Bool("log", false, "run with structured logger. Default: false")
	local := flag.Bool("local", false, "run servers locally. Default: false")
	throughput := flag.Int("throughput", 0, "target throughput")
	runs := flag.Int("runs", 0, "number of runs of each benchmark")
	steps := flag.Int("steps", 0, "number of increments on throughput vs latency benchmarks")
	dur := flag.Int("dur", 0, "number of seconds throughput vs latency benchmarks should run")
	broadcastID := flag.Uint64("BroadcastID", 0, "read BroadcastID from log")
	flag.Parse()

	if *broadcastID > 0 {
		readLog(*broadcastID, *runSrv)
		return
	}

	if os.Getenv("SERVER") == "1" {
		*runSrv = true
	}

	if os.Getenv("THROUGHPUT") != "" {
		var err error
		*throughput, err = strconv.Atoi(os.Getenv("THROUGHPUT"))
		if err != nil {
			panic(err)
		}
	}

	if os.Getenv("STEPS") != "" {
		var err error
		*steps, err = strconv.Atoi(os.Getenv("STEPS"))
		if err != nil {
			panic(err)
		}
	}

	if os.Getenv("DUR") != "" {
		var err error
		*dur, err = strconv.Atoi(os.Getenv("DUR"))
		if err != nil {
			panic(err)
		}
	}

	srvID := *id
	if srvID < 0 && *runSrv {
		var err error
		srvID, err = strconv.Atoi(os.Getenv("ID"))
		if err != nil {
			panic(err)
		}
	}

	bT := *benchTypeIndex
	if bT <= 0 {
		var err error
		bT, err = strconv.Atoi(os.Getenv("BENCH"))
		if err != nil {
			panic(err)
		}
	}
	benchType, ok := mapping[bT]
	if !ok {
		panic("invalid bench type")
	}

	if *runSrv {
		runServer(benchType, srvID, servers, *withLogger, *memProfile)
	} else {
		runBenchmark(benchType, clients, *throughput, *numClients, *clientBasePort, *steps, *runs, *dur, *local, servers, *memProfile, *withLogger)
	}
}

func runBenchmark(name string, clients ServerEntry, throughput, numClients, clientBasePort, steps, runs, dur int, local bool, srvAddrs map[int]Server, memProfile, withLogger bool) {
	options := make([]bench.RunOption, 0)
	if withLogger {
		file, err := os.Create("log.Clients.json")
		if err != nil {
			panic(err)
		}
		loggerOpts := &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}
		handler := slog.NewJSONHandler(file, loggerOpts)
		logger := slog.New(handler)
		options = append(options, bench.WithLogger(logger))
	}
	var srvAddresses []string
	if !local {
		options = append(options, bench.RunExternal())
		if srvAddrs == nil {
			panic("srvAddrs cannot be nil when not running locally")
		}
		srvAddresses = make([]string, len(srvAddrs))
		for _, srv := range srvAddrs {
			srvAddresses[srv.ID] = fmt.Sprintf("%s:%s", srv.Addr, srv.Port)
		}
		options = append(options, bench.WithSrvAddrs(srvAddresses))
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
	// must be placed as the last option because it overwrites numClients
	if clients != nil {
		clientsMap := make(map[int]string, len(clients))
		for id, entry := range clients {
			clientsMap[id] = fmt.Sprintf("%s:%s", entry.Addr, entry.Port)
		}
		options = append(options, bench.WithClients(clientsMap))
	}
	bench.RunThroughputVsLatencyBenchmark(name, options...)
}

type BenchmarkServer interface {
	Start(local bool)
	Stop()
}

func runServer(benchType string, id int, srvAddrs map[int]Server, withLogger, memprofile bool) {
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

	if memprofile {
		runtime.GC()
		cpuProfile, err := os.Create(fmt.Sprintf("cpuprofile.%v", id))
		if err != nil {
			panic(err)
		}
		memProfile, err := os.Create(fmt.Sprintf("memprofile.%v", id))
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(cpuProfile)
		defer pprof.StopCPUProfile()
		defer pprof.WriteHeapProfile(memProfile)
	}
	srv.Start(false)
	//fmt.Println("Press any key to stop server")
	//fmt.Scanln()
	//srv.Stop()
}

type logEntry struct {
	gorums.LogEntry
	Source struct {
		Function string `json:"function"`
		File     string `json:"file"`
		Line     int    `json:"line"`
	} `json:"source"`
}

func readLog(broadcastID uint64, server bool) {
	if !server {
		fmt.Println()
		fmt.Println("=============")
		fmt.Println("Reading:", "log.Clients.json")
		fmt.Println()
		file, err := os.Open("log.Clients.json")
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(file)
		// optionally, resize scanner's capacity for lines over 64K, see next example
		for scanner.Scan() {
			var entry logEntry
			json.Unmarshal(scanner.Bytes(), &entry)
			/*if entry.Cancelled {
				fmt.Println("BroadcastID", entry.BroadcastID, "msg:", entry.Msg, "err:", entry.Err)
			}*/
			if entry.Level == "Info" {
				if entry.Err != nil {
					fmt.Println(entry.Msg, ", err:", entry.Err, "msgID", entry.MsgID, ", nodeAddr", entry.NodeAddr, ", MachineID", entry.MachineID)
				} else {
					fmt.Println(entry.Msg, "msgID", entry.MsgID, ", method:", entry.Method, ", nodeAddr", entry.NodeAddr, ", MachineID", entry.MachineID)
				}
			}
		}
		return
	}
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
			if broadcastID == 1 {
				fmt.Println("msg:", entry.Msg, "err:", entry.Err)
			}
		}
	}
}
