package main

import (
	"log"

	"github.com/aleksander-vedvik/Master/gorums"
)

var values = []string{"val 1", "val 2", "val 3"}

func main() {
	srvAddresses := []string{"localhost:5000", "localhost:5001", "localhost:5002"}
	client := gorums.NewStorageClient(srvAddresses)
	log.Println("Created client...")
	val, err := client.ReadValue()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Reading inital value (should be empty): \"", val, "\"")
	log.Println("Writing values...")
	for _, val := range values {
		err = client.WriteValue(val)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("\t- Wrote:", val)
	}

	val, err = client.ReadValue()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Reading last value:", val)
}
