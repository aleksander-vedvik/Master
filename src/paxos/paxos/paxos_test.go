package main

import (
	"fmt"
	"paxos/client"
	"paxos/server"
	"testing"
)

func createSrvs(numSrvs int) ([]*server.PaxosServer, map[int]string, func()) {
	srvAddrs := make(map[int]string)
	for i := 0; i < numSrvs; i++ {
		srvAddrs[i] = fmt.Sprintf("127.0.0.1:%v", 5000+i)
	}
	srvs := make([]*server.PaxosServer, numSrvs)
	for id := range srvAddrs {
		srvs[id] = server.NewPaxosServer(id, srvAddrs, true)
		srvs[id].Start()
	}
	return srvs, srvAddrs, func() {
		for id := range srvAddrs {
			srvs[id].Stop()
		}
	}
}

func createClients(numClients int, srvAddrs map[int]string, known ...int) ([]*client.StorageClient, func()) {
	addrs := make([]string, 0, len(srvAddrs))
	for id, addr := range srvAddrs {
		if len(known) > 0 {
			for _, k := range known {
				if id == k {
					addrs = append(addrs, addr)
				}
			}
		} else {
			addrs = append(addrs, addr)
		}
	}
	clients := make([]*client.StorageClient, numClients)
	for i := 0; i < numClients; i++ {
		clients[i] = client.NewStorageClient(i, addrs)
	}
	return clients, func() {
		for id := range clients {
			clients[id].Stop()
		}
	}
}

func BenchmarkPaxos(b *testing.B) {
	// go test -bench=BenchmarkPaxos -benchmem -count=2 -run=^# -benchtime=100x
	srvs, srvAddrs, cleanup := createSrvs(3)
	defer cleanup()

	clients, stop := createClients(1, srvAddrs, 1)
	defer stop()

	for i, client := range clients {
		client.Write(fmt.Sprintf("val %v, client %v", i, i))
	}

	for _, srv := range srvs {
		srv.PrintStats()
	}
	b.ResetTimer()
	for c, client := range clients {
		b.RunParallel(func(pb *testing.PB) {
			for i := 0; pb.Next(); i++ {
				client.Write(fmt.Sprintf("val %v, client %v", i, c))
			}
		})
	}
	b.StopTimer()

	fmt.Println()
	fmt.Println()
	for _, srv := range srvs {
		//fmt.Println(srv.Status())
		srv.PrintStats()
	}
	fmt.Println()
	fmt.Println()

	// create 3 servers
	// create 5 servers
	// create 7 servers

	// run 1000 reqs with
	// 	- 1 client
	// 	- 5 clients
	// 	- 10 clients

	// run 10 000 reqs with
	// 	- 1 client
	// 	- 5 clients
	// 	- 10 clients

	// run 50 000 reqs with
	// 	- 1 client
	// 	- 5 clients
	// 	- 10 clients

	// measure average, mean, min, and max processing time
	// measure throughput: reqs/sec
	// measure successful and failed reqs
}
