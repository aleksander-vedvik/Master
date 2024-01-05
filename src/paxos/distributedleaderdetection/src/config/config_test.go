package config

import (
	"os"
	"testing"
)

type N map[int]string

var serverAddresses = N{
	0: "127.0.0.1:5000",
	1: "127.0.0.1:5001",
	2: "127.0.0.1:5002",
	3: "127.0.0.1:5003",
}

func (n N) contains(key int, value string) bool {
	for k, v := range n {
		if key == k && value == v {
			return true
		}
	}
	return false
}

func TestConfig(t *testing.T) {
	idWant := 2
	err := os.Setenv("ID", "2")
	if err != nil {
		t.Fatalf("Could not set environment variable. Error: %s\n", err)
		return
	}
	idGot, c, err := ReadConfig("../config.json")
	if idWant != idGot {
		t.Fatalf("Wrong id! Got: %v, Want: %v\n, Error: %s", idGot, idWant, err)
		return
	}
	if err != nil {
		t.Fatalf("Could not read config. Error: %s\n", err)
		return
	}
	for _, node := range c.Nodes {
		if !serverAddresses.contains(node.Id, node.Address) {
			t.Fatalf("Wrong server addresss. NodeId: %v, Address: %s", node.Id, node.Address)
		}
	}
}
