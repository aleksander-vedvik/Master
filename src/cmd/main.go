package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/aleksander-vedvik/Master/protos"
	"github.com/aleksander-vedvik/Master/server"
	"github.com/relab/gorums"
)

func main() {
	fmt.Println("Hello, world!")
}

func ExampleStorageServer(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	gorumsSrv := gorums.NewServer()
	srv := server.StorageSrv{State: &pb.State{}}
	pb.RegisterStorageServer(gorumsSrv, &srv)
	gorumsSrv.Serve(lis)
}
