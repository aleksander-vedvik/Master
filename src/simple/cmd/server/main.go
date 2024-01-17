package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aleksander-vedvik/Master/storage"
)

func main() {
	//startDockerServer()
	//startServers()
	startServersTree()
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

	srv := storage.NewStorageServer("")
	addr := srv.StartServer(srvAddresses[*id-1])
	//srv.AddConfig(srvAddresses)
	log.Printf("Server started. Listening on address: %s\n", addr)
	fmt.Scanln()
	log.Println("Server stopped.")
}

func startServers() {
	numServers := 10
	srvAddresses := make([]string, numServers)
	for i := range srvAddresses {
		srvAddresses[i] = fmt.Sprintf("localhost:%v", 5000+i)
	}
	for _, srvAddr := range srvAddresses {
		srv := storage.NewStorageServer(srvAddr)
		_ = srv.StartServer(srvAddr)
		go srv.AddConfig(srvAddresses)
	}
	log.Printf("Servers started. Listening on addresses: %s\n", srvAddresses)
	fmt.Scanln()
	log.Println("Servers stopped.")
}

func startServersTree() {
	server1 := "localhost:5000"
	server2 := "localhost:5001"
	server3 := "localhost:5002"
	server4 := "localhost:5003"
	server5 := "localhost:5004"
	server6 := "localhost:5005"
	server7 := "localhost:5006"
	server8 := "localhost:5007"
	server9 := "localhost:5008"
	server10 := "localhost:5009"
	firstLayerTree := []string{server2, server3, server4}
	secondLayerTree := []string{server5, server6, server7}
	thirdLayerTree := []string{server8, server9, server10}

	createServer(server1, firstLayerTree)
	createServer(server3, []string{server1})
	createServer(server4, []string{server1})
	createServer(server2, secondLayerTree)
	createServer(server5, thirdLayerTree)
	createServer(server6, []string{server1})
	createServer(server7, []string{server1})
	createServer(server8, []string{server1})
	createServer(server9, []string{server1})
	createServer(server10, []string{server1})

	go func() {
		time.Sleep(1 * time.Second)
		for {
			time.Sleep(5 * time.Second)
			fmt.Println()
		}
	}()

	log.Printf("Servers started...")
	fmt.Scanln()
	log.Println("Servers stopped.")
}

func createServer(srvAddr string, srvAddresses []string) {
	srv := storage.NewStorageServer(srvAddr)
	_ = srv.StartServer(srvAddr)
	go srv.AddConfig(srvAddresses)
}

func startDockerServer() {
	srvAddresses := []string{"localhost:5000", "localhost:5001", "localhost:5002"}
	srv := storage.NewStorageServer("")
	srv.Start("0.0.0.0:8080")
	srv.AddConfig(srvAddresses)
	log.Println("Server stopped.")
}
