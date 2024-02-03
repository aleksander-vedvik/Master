package main

import (
	"fmt"
	"log"

	"github.com/aleksander-vedvik/Master/storage/server"
)

func main() {
	startServers()
}

func startServers() {
	numServers := 3
	srvAddresses := make([]string, numServers)
	for i := range srvAddresses {
		srvAddresses[i] = fmt.Sprintf("localhost:%v", 5000+i)
	}
	for _, srvAddr := range srvAddresses {
		peers := append(srvAddresses, "localhost:5010")
		srv := server.NewStorageServer(srvAddr, peers)
		go srv.Start(srvAddr)
		go srv.Run()
	}
	log.Printf("Servers started. Listening on addresses: %s\n", srvAddresses)
	fmt.Scanln()
	log.Println("Servers stopped.")
}
