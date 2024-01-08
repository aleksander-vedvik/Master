package main

import (
	"math/rand"
	"time"
)

type Node struct {
	Id       int
	Status   Status
	Messages int
	Data     Data
}

func (n *Node) DoSomething() {
	for {
		time.Sleep(time.Duration(3+rand.Intn(7)) * time.Second)
		if n.Status == Online {
			n.Status = Offline
		} else {
			n.Status = Online
		}
	}
}

type Data struct {
}

func TestNodes() []*Node {
	nodes := make([]*Node, 0)
	for i := 0; i < 9; i++ {
		status := Online
		if i%4 == 0 {
			status = Offline
		}
		node := &Node{
			Id:     i + 1,
			Status: status,
		}
		go node.DoSomething()
		nodes = append(nodes, node)
	}
	return nodes
}
