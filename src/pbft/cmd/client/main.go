package main

import (
	"fmt"
	"log"
	"time"

	"pbft/node"
)

var values = []string{"val 1", "val 2", "val 3"}

func main() {
	client2()
}

func client2() {
	time.Sleep(1 * time.Second)
	numServers := 1
	srvAddresses := make([]string, numServers)
	for i := range srvAddresses {
		srvAddresses[i] = fmt.Sprintf("localhost:%v", 5000+i)
	}
	client := node.NewStorageClient(srvAddresses, "client")
	log.Println("Created client...")
	log.Println("\t- Only writing to servers", srvAddresses)
	log.Println("Writing value", values[0])
	err := client.WriteValue(values[0])
	if err != nil {
		log.Println(err)
	}
	time.Sleep(5 * time.Second)
	fmt.Println()
	log.Println("Client done...")
	fmt.Println()
}
