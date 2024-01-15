package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aleksander-vedvik/Master/storage"
)

var values = []string{"val 1", "val 2", "val 3"}

func main() {
	client2()
}

func client1() {
	time.Sleep(1 * time.Second)
	srvAddresses := []string{"localhost:5000"}
	client := storage.NewStorageClient(srvAddresses)
	log.Println("Created client...")
	log.Println("\t- Only writing to servers", srvAddresses)
	val, err := client.ReadValue()
	if err != nil {
		log.Println(err)
	}
	log.Println("Reading inital value (should be empty): \"", val, "\"")
	fmt.Println()
	time.Sleep(5 * time.Second)
	fmt.Println()
	log.Println("Writing values...")
	for _, val := range values {
		err = client.WriteValue(val)
		if err != nil {
			log.Println(err)
		}
		log.Println("\t- Wrote:", val)
	}
	fmt.Println()
	time.Sleep(5 * time.Second)
	fmt.Println()
	val, err = client.ReadValue()
	if err != nil {
		log.Println(err)
	}
	log.Println("Reading last value:", val)
	log.Println("Client done...")
	fmt.Println()
}

func client2() {
	time.Sleep(1 * time.Second)
	srvAddresses := []string{"localhost:5000"}
	client := storage.NewStorageClient(srvAddresses)
	log.Println("Created client...")
	log.Println("\t- Only writing to servers", srvAddresses)
	log.Println("Writing value", values[0])
	err := client.WriteValue(values[0])
	if err != nil {
		log.Println(err)
	}
	time.Sleep(10 * time.Second)
	fmt.Println()
	val, err := client.ReadValue()
	if err != nil {
		log.Println(err)
	}
	log.Println("Reading last value:", val)
	log.Println("Client done...")
	fmt.Println()
}
