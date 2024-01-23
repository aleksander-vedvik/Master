package main

type Nodes []*Node

func NewNodes() Nodes {
	return make(Nodes, 0, 3)
}

type Node struct {
	Id         int
	Status     Status
	Messages   int
	LastCmdRun string
	Data       Data
	Address    string
}

type Data struct {
	Val string
}

func NewNode(address string, id int) *Node {
	return &Node{
		Id:       id,
		Address:  address,
		Status:   Online,
		Messages: 0,
		Data: Data{
			Val: "",
		},
	}
}

type Payload struct {
	Address  string `json:"address"`
	Messages int    `json:"messages"`
	Method   string `json:"method"`
	Value    string `json:"value"`
}

func (nodes Nodes) UpdateNode(payload Payload) {
	for _, node := range nodes {
		if node.Address == payload.Address {
			node.Messages = payload.Messages
			node.LastCmdRun = payload.Method
			node.Data.Val = payload.Value
			break
		}
	}
}
