package main

/*import (
	"fmt"
	"log"
	"sync"
	"time"

	node "github.com/aleksander-vedvik/benchmark/pbft/client"
)

var values = []string{"val 1", "val 2", "val 3", "val 4", "val 5", "val 6"}

func main() {
	client2()
}

func client2() {
	time.Sleep(2 * time.Second)
	numServers := 1
	srvAddresses := make([]string, numServers)
	for i := range srvAddresses {
		srvAddresses[i] = fmt.Sprintf("localhost:%v", 5000+i)
	}
	client := node.New("127.0.0.1:8080", srvAddresses)
	fmt.Println()
	log.Println("Created client...")
	log.Println("\t- Only writing to servers", srvAddresses)
	fmt.Println()
	log.Println("Writing values:", values[0], values[1], values[2])
	for i, val := range values {
		if i >= 3 {
			break
		}
		err := client.WriteValue(val)
		if err != nil {
			log.Println(err)
		}
	}
	fmt.Println()
	log.Println("Writing values concurrently:", values[3], values[4], values[5])
	var wg sync.WaitGroup
	for i := 3; i < len(values); i++ {
		wg.Add(1)
		go func(i int) {
			client.WriteValue(values[i])
			wg.Done()
		}(i)
	}
	wg.Wait()
	time.Sleep(1 * time.Second)
	fmt.Println()
	log.Println("Client done...")
	fmt.Println()
}
*/
