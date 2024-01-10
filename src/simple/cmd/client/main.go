package main

import (
	"log"
	"time"

	"github.com/aleksander-vedvik/Master/storage"
)

var values = []string{"val 1", "val 2", "val 3"}

func main() {
	time.Sleep(1 * time.Second)
	srvAddresses := []string{"localhost:5000", "localhost:5001", "localhost:5002"}
	client := storage.NewStorageClient(srvAddresses)
	log.Println("Created client...")
	val, err := client.ReadValue()
	if err != nil {
		log.Println(err)
	}
	log.Println("Reading inital value (should be empty): \"", val, "\"")
	log.Println("Writing values...")
	for _, val := range values {
		err = client.WriteValue(val)
		if err != nil {
			log.Println(err)
		}
		log.Println("\t- Wrote:", val)
	}

	val, err = client.ReadValue()
	if err != nil {
		log.Println(err)
	}
	log.Println("Reading last value:", val)
}
