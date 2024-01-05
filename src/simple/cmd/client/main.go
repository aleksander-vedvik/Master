package main

import "github.com/aleksander-vedvik/Master/gorums"

func main() {
	srvAddresses := []string{"localhost:5000", "localhost:5001", "localhost:5002"}
	client := gorums.NewStorageClient(srvAddresses)
	client.WriteValue("test")
}
