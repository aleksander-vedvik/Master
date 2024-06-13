package main

import (
	"fmt"
	"os"

	bench "github.com/aleksander-vedvik/benchmark/benchmark"
	"gopkg.in/yaml.v3"
)

type Server struct {
	ID   int    `yaml:"id"`
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
}

// ServerEntry represents an entry in the servers list
type ServerEntry map[int]Server

// Config represents the configuration containing servers
type Config struct {
	Servers []ServerEntry `yaml:"servers"`
	Clients []ServerEntry `yaml:"clients"`
}

func getConfig() (srvs, clients ServerEntry) {
	confType := os.Getenv("CONF")
	if confType == "" {
		confType = "local"
	}
	confPath := fmt.Sprintf("conf.%s.yaml", confType)
	data, err := os.ReadFile(confPath)
	if err != nil {
		panic(err)
	}
	var c Config
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		panic(err)
	}
	srvs = make(map[int]Server, len(c.Servers))
	for _, srv := range c.Servers {
		for id, info := range srv {
			srvs[id] = info
		}
	}
	clients = make(map[int]Server, len(c.Servers))
	for _, client := range c.Clients {
		for id, info := range client {
			clients[id] = info
		}
	}
	return srvs, clients
}

type mappingType map[int]string

var mapping mappingType = map[int]string{
	1: bench.PaxosBroadcastCall,
	2: bench.PaxosQuorumCall,
	3: bench.PaxosQuorumCallBroadcastOption,
	4: bench.PBFTWithGorums,
	5: bench.PBFTWithoutGorums,
}

func (m mappingType) String() string {
	ret := "\n"
	ret += "\t1: " + m[1] + "\n"
	ret += "\t2: " + m[2] + "\n"
	ret += "\t3: " + m[3] + "\n"
	ret += "\t4: " + m[4] + "\n"
	ret += "\t5: " + m[5] + "\n"
	return ret
}
