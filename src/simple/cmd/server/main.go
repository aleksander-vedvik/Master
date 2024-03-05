package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/aleksander-vedvik/Master/storage/server"
)

func initREMOVE() {
	protoc := `protoc -I=$(go list -m -f {{.Dir}} github.com/relab/gorums):. \
	--go_out=paths=source_relative:. \
	--gorums_out=paths=source_relative:. \
	../Master/src/simple/protos/storage.proto`

	cmd := exec.Command(protoc)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	fmt.Println(string(stdout))
}

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
		go srv.Start()
		go srv.Run()
	}
	log.Printf("Servers started. Listening on addresses: %s\n", srvAddresses)
	srv := server.NewStorageServer(extra, append(srvAddresses, extra))
	time.Sleep(11 * time.Second)
	fmt.Println()
	go srv.Start()
	go srv.Run()
	fmt.Scanln()
	log.Println("Servers stopped.")
}
