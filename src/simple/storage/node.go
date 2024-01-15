package storage

import (
	"github.com/relab/gorums"
)

type Node struct {
	node *gorums.GorumsNode
}

func NewNode() *Node {
	node := gorums.NewGorumsNode()

	node.RegisterClient(NewStorageClient([]string{}))
	node.RegisterServer(NewStorageServer([]string{}))

	node.NewConfiguration([]string{})

	return &Node{node}
}
