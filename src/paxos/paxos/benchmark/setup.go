package main

import (
	"fmt"
	"paxos/client"
	"paxos/server"
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
