package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"
)

type Config struct {
	Nodes []Node `json:"nodes"`
}

type Node struct {
	Id      int    `json:"id"`
	Address string `json:"address"`
}

func ReadConfig(path string) (int, *Config, error) {
	c, err := openFile(path)
	if err != nil {
		return -1, nil, err
	}
	id := os.Getenv("ID")
	node := getNode(id, c.Nodes)
	if node == nil {
		return -1, nil, errors.New("Environment variable id not set!")
	}
	return node.Id, c, nil
}

func openFile(path string) (*Config, error) {
	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	var c Config
	b, err := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func getNode(id string, nodes []Node) *Node {
	for _, node := range nodes {
		if strconv.Itoa(node.Id) == id {
			return &node
		}
	}
	return nil
}

func (c *Config) GetNodeIds() []int {
	ids := make([]int, 0, len(c.Nodes))
	for _, node := range c.Nodes {
		ids = append(ids, node.Id)
	}
	return ids
}

func (c *Config) GetNodeAddressess() map[int]string {
	addresses := make(map[int]string, 0)
	for _, node := range c.Nodes {
		addresses[node.Id] = node.Address
	}
	return addresses
}
