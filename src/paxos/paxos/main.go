package main

import (
	"flag"
	"fmt"
	"paxos/server"
)

var srvAddrs = map[int]string{
	0: "127.0.0.1:5000",
	1: "127.0.0.1:5001",
	2: "127.0.0.1:5002",
}

func main() {
	id := flag.Int("id", 0, "nodeID")
	flag.Parse()

	srv := server.NewPaxosServer(*id, srvAddrs)
	srv.Start()

	fmt.Scanln()
}
