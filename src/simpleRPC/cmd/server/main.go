package main

import (
	"fmt"
	"log"

	"github.com/aleksander-vedvik/Master/storage"
)

func main() {
	srvAddresses := []string{"localhost:5000", "localhost:5001", "localhost:5002", "localhost:5003"}
	for _, srvAddr := range srvAddresses {
		srv := storage.NewStorageServer(srvAddresses, srvAddr)
		srv.StartServer()
	}
	log.Printf("Servers started. Listening on addresses: %s\n", srvAddresses)
	fmt.Scanln()
	log.Println("Servers stopped.")
}
