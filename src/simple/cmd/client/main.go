package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aleksander-vedvik/Master/storage/client"
)

var values = []string{"val 1", "val 2", "val 3"}

func main() {
	c()
}

func c() {
	time.Sleep(2 * time.Second)
	srvAddresses := []string{"localhost:5000"}
	c := client.NewStorageClient(srvAddresses)
	fmt.Println()
	log.Println("Created client...")
	log.Println("\t- Only writing to servers", srvAddresses)
	log.Println("Writing value", values[0])
	err := c.WriteValue(values[0])
	if err != nil {
		log.Println(err)
	}
	fmt.Println()
	log.Println("Client received first response...")
	fmt.Println()
	time.Sleep(11 * time.Second)
	fmt.Println()
	log.Println("Writing value", values[1])
	err = c.WriteValue(values[1])
	if err != nil {
		log.Println(err)
	}
	fmt.Println()
	log.Println("Client done...")
	fmt.Println()
}
