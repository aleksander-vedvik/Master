package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aleksander-vedvik/Master/storage/server"
)

func main() {
	startServers()
}

func startServers() {
	extra := "localhost:5003"
	numServers := 3
	srvAddresses := make([]string, numServers)
	for i := range srvAddresses {
		srvAddresses[i] = fmt.Sprintf("localhost:%v", 5000+i)
	}
	for _, srvAddr := range srvAddresses {
		peers := append(srvAddresses, extra)
		srv := server.NewStorageServer(srvAddr, peers)
		go srv.Start(srvAddr)
		go srv.Run()
	}
	log.Printf("Servers started. Listening on addresses: %s\n", srvAddresses)
	srv := server.NewStorageServer(extra, append(srvAddresses, extra))
	time.Sleep(11 * time.Second)
	fmt.Println()
	go srv.Start(extra)
	go srv.Run()
	fmt.Scanln()
	log.Println("Servers stopped.")
}
