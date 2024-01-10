package main

import (
	"math/rand"
	"time"
)

type Node struct {
	Id       uint32
	Status   Status
	Messages int
	Round    int
	Data     Data
	Address  string
}

type Nodes []*Node

func CreateNodes(srvAddrs []string) Nodes {
	nodes := make([]*Node, 0)
	for _, addr := range srvAddrs {
		node := &Node{
			Status:   Online,
			Messages: 0,
			Round:    0,
			Address:  addr,
		}
		nodes = append(nodes, node)
	}
	n := Nodes(nodes)
	return n
}

func (n Nodes) AddIDs(ids []uint32, addresses []string) {
	for i, node := range n {
		node.Id = ids[i]
		node.Address = addresses[i]
	}
	/*for i, address := range addresses {
		for _, node := range n {
			if node.Address == address {
				node.Id = ids[i]
				break
			}
		}
	}*/
}

func (n Nodes) GetNode(id uint32) *Node {
	for _, node := range n {
		if node.Id == id {
			return node
		}
	}
	return nil
}

func (n Nodes) GetAddresses() []string {
	addrs := make([]string, len(n))
	for i, node := range n {
		addrs[i] = node.Address + ":8080"
	}
	return addrs
}

func (n Nodes) GetIds() []uint32 {
	ids := make([]uint32, len(n))
	for i, node := range n {
		ids[i] = node.Id
	}
	return ids
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
			Status:   status,
			Messages: 0,
			Round:    0,
		}
		go node.DoSomething()
		nodes = append(nodes, node)
	}
	return nodes
}
