package main

import (
	bench "github.com/aleksander-vedvik/benchmark/benchmark"
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

func main() {
	//bench.RunSingleBenchmark("Simple")
	//bench.RunSingleBenchmark("PaxosQC")
	//bench.RunSingleBenchmark("Paxos")
	//bench.RunSingleBenchmark("PBFT")
	//bench.RunSingleBenchmark("PBFT.S")
	//bench.RunThroughputVsLatencyBenchmark("Paxos.BroadcastCall")
	//bench.RunThroughputVsLatencyBenchmark("Paxos.QuorumCall")
	bench.RunThroughputVsLatencyBenchmark("Paxos.QuorumCallBroadcastOption")
	//bench.RunThroughputVsLatencyBenchmark("PBFT.With.Gorums", 5000, 500)
	//bench.RunThroughputVsLatencyBenchmark("PBFT.Without.Gorums", 5000, 500)
	//bench.RunThroughputVsLatencyBenchmark("PBFT.NoOrder")
	//bench.RunThroughputVsLatencyBenchmark("PBFT.Order")
}

//func main() {
//id := flag.Int("id", 0, "nodeID")
//t := flag.Int("type", 0, "server = 0, client = 1")
//reqs := flag.Int("reqs", 0, "number of reqs")
//flag.Parse()

//if *t == 0 {
//srv := server.NewPaxosServer(*id, srvAddrs)
//srv.Start()

//time.Sleep(10 * time.Second)
//} else {
//time.Sleep(1 * time.Second)
//addrs := make([]string, 0, len(srvAddrs))
//for _, addr := range srvAddrs {
//addrs = append(addrs, addr)
//}
//c := client.NewStorageClient(*id, addrs)
//for i := 0; i < *reqs; i++ {
//c.Write(fmt.Sprintf("val %v, client %v", i, *id))
//}
//}
//}
