package test

import (
	"log"

	"github.com/aleksander-vedvik/Master/storage"
	"github.com/relab/gorums"
)

func CreateConfig() {
	quorum := NewGorumsManager(configuration)

	resp, err := quorum.Write("test")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
}

func configuration(q *GorumsQSpec) {
	srvAddresses := []string{"localhost:5000", "localhost:5001", "localhost:5002"}
	q.QuorumSize = len(srvAddresses)
	q.QSpec = storage.NewQSpec(len(srvAddresses))
	q.Nodes = gorums.WithNodeList(srvAddresses)
	q.OnlyRunWhenQuorum = false
}
