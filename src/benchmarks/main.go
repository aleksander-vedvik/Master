package main

import (
	"flag"
	"fmt"
	"os"

	bench "github.com/aleksander-vedvik/benchmark/benchmark"
	pbftServer "github.com/aleksander-vedvik/benchmark/pbft/server"
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
	t := flag.Int("server", 0, "no = 0, server = 1")
	b := flag.Int("run", 0, "type of benchmark to run"+mapping.String())
	memProfile := flag.Bool("mem", false, "create memory and cpu profile")
	numClients := flag.Int("clients", 0, "number of clients to run")
	clientBasePort := flag.Int("port", 0, "which base port to use for clients")
	wL := flag.Bool("leaderDetection", true, "run without leader detection. Default: true")
	local := flag.Bool("local", true, "run servers locally. Default: true")
	throughput := flag.Int("throughput", 0, "target throughput")
	flag.Parse()

	benchType, ok := mapping[*b]
	if !ok {
		panic("invalid bench type")
	}

	if *t == 0 {
		runBenchmark(benchType, *throughput, *numClients, *clientBasePort, *local, servers, *memProfile)
	} else {
		runServer(*b, *id, servers, *wL)
	}
}

func runBenchmark(name string, throughput, numClients, clientBasePort int, local bool, srvAddrs map[int]Server, memProfile bool) {
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
		options = append(options, bench.ThroughputIncr(throughput/10))
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

func runServer(benchType int, id int, srvAddrs map[int]Server, withoutLeaderDetection bool) {
	srvAddresses := make([]string, len(srvAddrs))
	for _, srv := range srvAddrs {
		srvAddresses[srv.ID] = fmt.Sprintf("%s:%s", srv.Addr, srv.Port)
	}
	var srv BenchmarkServer
	switch benchType {
	case 0:
	case 1:
	case 2:
	case 3:
		srv = pbftServer.New(srvAddresses[id], srvAddresses, withoutLeaderDetection)
	case 4:
	case 5:
	case 6:
	case 7:
	}
	srv.Start()
	fmt.Println("Press any key to stop server")
	fmt.Scanln()
	srv.Stop()
}
