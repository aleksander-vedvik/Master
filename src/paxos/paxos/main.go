package main

import (
	"flag"
	"fmt"
	"paxos/client"
	"paxos/server"
)

var srvAddrs = map[int]string{
	0: "127.0.0.1:5000",
	1: "127.0.0.1:5001",
	2: "127.0.0.1:5002",
}

func main() {
	id := flag.Int("id", 0, "nodeID")
	t := flag.Int("type", 0, "server = 0, client = 1")
	flag.Parse()

	if *t == 0 {
		srv := server.NewPaxosServer(*id, srvAddrs)
		srv.Start()

		fmt.Scanln()
	} else {
		addrs := make([]string, 0, len(srvAddrs))
		for _, addr := range srvAddrs {
			addrs = append(addrs, addr)
		}
		c := client.NewStorageClient(addrs)
		c.Write("test value")
	}
}
