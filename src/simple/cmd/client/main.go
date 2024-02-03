package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aleksander-vedvik/Master/storage"
)

var values = []string{"val 1", "val 2", "val 3"}

func main() {
	client()
}

func client() {
	time.Sleep(2 * time.Second)
	srvAddresses := []string{"localhost:5000"}
	client := storage.NewStorageClient(srvAddresses)
	fmt.Println()
	log.Println("Created client...")
	log.Println("\t- Only writing to servers", srvAddresses)
	log.Println("Writing value", values[0])
	err := client.WriteValue(values[0])
	if err != nil {
		log.Println(err)
	}
	fmt.Println()
	log.Println("Client done...")
	fmt.Println()
}
