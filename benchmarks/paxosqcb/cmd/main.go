package main

/*import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aleksander-vedvik/benchmark/paxosqcb/server"
)

func main() {
	var addr string
	var srvAddrs []string

	peers := os.Getenv("PEERS")
	if peers == "" {
		port := flag.Int("port", 0, "listening port")
		flag.Parse()

		if *port == 0 {
			return
		}

		srvAddrs = []string{"127.0.0.1:5000", "127.0.0.1:5001", "127.0.0.1:5002"}
		addr = fmt.Sprintf("127.0.0.1:%v", *port)
	} else {
		srvAddrs = strings.Split(peers, ",")
		addr = "0.0.0.0:8080"
	}

	srv := server.New(addr, srvAddrs)
	srv.Start()
	defer srv.Stop()
	fmt.Scanln()
}
*/
