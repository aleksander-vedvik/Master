package main

import (
	"math/rand"
	"time"
)

type Node struct {
	Id       int
	Status   Status
	Messages int
	Round    int
	Data     Data
}

func (n *Node) DoSomething() {
	go n.work()
	for {
		time.Sleep(time.Duration(3+rand.Intn(7)) * time.Second)
		if n.Status == Online {
			n.Status = Offline
		} else {
			n.Status = Online
		}
	}
}

func (n *Node) work() {
	for {
		time.Sleep(1 * time.Second)
		if n.Status == Online {
			n.Messages += 300 + rand.Intn(700)
			n.Round += rand.Intn(3)
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
			Id:       i + 1,
			Status:   status,
			Messages: 0,
			Round:    0,
		}
		go node.DoSomething()
		nodes = append(nodes, node)
	}
	return nodes
}

func GetNode(id int, nodes []*Node) *Node {
	for _, node := range nodes {
		if node.Id == id {
			return node
		}
	}
	return nil
}
