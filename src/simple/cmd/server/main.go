package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/aleksander-vedvik/Master/storage"
)

func main() {
	//startDockerServer()
	startServers()
	//simpleServer()
}

func simpleServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func startServer() {
	srvAddresses := []string{"localhost:5000", "localhost:5001", "localhost:5002"}
	id := flag.Int("id", 0, "id of server")
	flag.Parse()

	srv := storage.NewStorageServer()
	addr := srv.StartServer(srvAddresses[*id-1])
	//srv.AddConfig(srvAddresses)
	log.Printf("Server started. Listening on address: %s\n", addr)
	fmt.Scanln()
	log.Println("Server stopped.")
}

func startServers() {
	srvAddresses := []string{"localhost:5000", "localhost:5001", "localhost:5002", "localhost:5003"}
	for _, srvAddr := range srvAddresses {
		srv := storage.NewStorageServer()
		_ = srv.StartServer(srvAddr)
		go srv.AddConfig(srvAddresses)
	}
	log.Printf("Servers started. Listening on addresses: %s\n", srvAddresses)
	fmt.Scanln()
	log.Println("Servers stopped.")
}

func startDockerServer() {
	srvAddresses := []string{"localhost:5000", "localhost:5001", "localhost:5002"}
	srv := storage.NewStorageServer()
	srv.Start("0.0.0.0:8080")
	srv.AddConfig(srvAddresses)
	log.Println("Server stopped.")
}
