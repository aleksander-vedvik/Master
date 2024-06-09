package main

import (
	"context"
	"fmt"
	"reliablebroadcast/client"
	"reliablebroadcast/server"
)

func main() {
	clientAddr := "127.0.0.1:8080"
	srvAddrs := []string{"127.0.0.1:5000", "127.0.0.1:5001", "127.0.0.1:5002"}
	c := client.New(clientAddr, srvAddrs, len(srvAddrs))

	for _, addr := range srvAddrs {
		srv := server.New(addr, srvAddrs)
		srv.Start()
	}

	resp, err := c.Broadcast(context.Background(), "test")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\nWrote \"%s\" to the servers!", resp.Data)
}
