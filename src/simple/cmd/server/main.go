package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aleksander-vedvik/Master/storage"
)

func main() {
	//startDockerServer()
	startServers()
}

func startServer() {
	srvAddresses := []string{"localhost:5000", "localhost:5001", "localhost:5002"}
	id := flag.Int("id", 0, "id of server")
	flag.Parse()

	srv := storage.NewStorageServer()
	addr := srv.StartServer(srvAddresses[*id-1])
	log.Printf("Server started. Listening on address: %s\n", addr)
	fmt.Scanln()
	log.Println("Server stopped.")
}

func startServers() {
	srvAddresses := []string{"localhost:5000", "localhost:5001", "localhost:5002", "localhost:5003"}
	for _, srvAddr := range srvAddresses {
		srv := storage.NewStorageServer()
		_ = srv.StartServer(srvAddr)
	}
	log.Printf("Servers started. Listening on addresses: %s\n", srvAddresses)
	fmt.Scanln()
	log.Println("Servers stopped.")
}

func startDockerServer() {
	srv := storage.NewStorageServer()
	srv.Start("localhost:8080")
	log.Println("Server stopped.")
}
