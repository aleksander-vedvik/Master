package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	paxos "paxos"
	pb "paxos/proto"

	gorums "github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var (
		saddrs           = flag.String("addrs", "", "server addresses separated by ','")
		numClientRequest = flag.Int("clientRequest", 0, "numberOfClientRequests")
		clientId         = flag.String("clientId", "", "Client Id, different for each client")
	)

	flag.Usage = func() {
		log.Printf("Usage: %s [OPTIONS]\nOptions:", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	addrs := strings.Split(*saddrs, ",")
	if len(addrs) == 0 {
		log.Fatalln("no server addresses provided")
	}

	n := *numClientRequest
	if n == 0 {
		log.Fatalln("no client requests are provided")
	}
	clientRequests := make([]string, 0)
	for k := 0; k < n; k++ {
		clientRequests = append(clientRequests, fmt.Sprint(k))
	}

	// start a initial proposer
	ClientStart(addrs, clientRequests, clientId)
}

// ClientStart creates the configuration with the list of replicas addresses, which are read from the
// commandline. From the list of clientRequests, send each request to the configuration and
// wait for the reply. Upon receiving the reply send the next request.
func ClientStart(addrs []string, clientRequests []string, clientId *string) {
	log.Printf("Connecting to %d Paxos replicas: %v", len(addrs), addrs)
	config, mgr := createConfiguration(addrs)
	defer mgr.Close()

	latencyArray := make([]float64, 0)
	for index, request := range clientRequests {
		req := pb.Value{ClientID: *clientId, ClientSeq: uint32(index), ClientCommand: request}
		_, latency := doSendRequest(config, &req)

		//log.Printf("response: %v\t for the client request: %v", resp, &req)

		latencyArray = append(latencyArray, latency.Seconds())
	}
	log.Println("Done sending requests")
	// TODO(BM): Insert logic for saving data to file
	f, err := os.Create(fmt.Sprintf("client=%v.csv", *clientId))
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err = f.WriteString("latency (s)\n")

	if err != nil {
		log.Fatal(err)
	}

	for _, data := range latencyArray {
		_, err = f.WriteString(fmt.Sprintf("%v\n", data))

		if err != nil {
			log.Fatal(err)
		}
	}
}

// Internal: doSendRequest can send requests to paxos servers by quorum call and
// for the response from the quorum function.
func doSendRequest(config *pb.Configuration, value *pb.Value) (*pb.Response, *time.Duration) {
	waitTimeForRequest := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), waitTimeForRequest)
	defer cancel()

	// BM: I assume this is where it makes the most sense to meassure the latency
	startTime := time.Now()
	resp, err := config.ClientHandle(ctx, value)
	latency := time.Since(startTime)
	log.Printf("latency: %v", latency)

	if err != nil {
		log.Fatalf("ClientHandle quorum call error: %v", err)
	}
	if resp == nil {
		log.Println("response is nil")
	}
	return resp, &latency
}

// createConfiguration creates the gorums configuration with the list of addresses.
func createConfiguration(addrs []string) (configuration *pb.Configuration, manager *pb.Manager) {
	mgr := pb.NewManager(gorums.WithDialTimeout(5*time.Second),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(), // block until connections are made
			grpc.WithTransportCredentials(insecure.NewCredentials()), // disable TLS
		),
	)
	quorumSize := (len(addrs)-1)/2 + 1
	qfspec := paxos.NewPaxosQSpec(quorumSize)
	config, err := mgr.NewConfiguration(qfspec, gorums.WithNodeList(addrs))
	if err != nil {
		log.Fatalf("Error in forming the configuration: %v\n", err)
		return nil, nil
	}
	return config, mgr
}
