package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/relab/gorums"
)

type logEntry struct {
	gorums.LogEntry
	Source struct {
		Function string `json:"function"`
		File     string `json:"file"`
		Line     int    `json:"line"`
	} `json:"source"`
}

func readLog(broadcastID uint64, server bool) {
	if !server {
		fmt.Println()
		fmt.Println("=============")
		fmt.Println("Reading:", "log.Clients.json")
		fmt.Println()
		file, err := os.Open("./logs/log.Clients.json")
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var entry logEntry
			json.Unmarshal(scanner.Bytes(), &entry)
			if entry.Level == "Info" {
				if entry.Err != nil {
					fmt.Println(entry.Msg, ", err:", entry.Err, "msgID", entry.MsgID, ", nodeAddr", entry.NodeAddr, ", MachineID", entry.MachineID)
				} else {
					fmt.Println(entry.Msg, "msgID", entry.MsgID, ", method:", entry.Method, ", nodeAddr", entry.NodeAddr, ", MachineID", entry.MachineID)
				}
			}
		}
		return
	}
	logFiles := []string{"./logs/log.10.0.0.5:5000.json", "./logs/log.10.0.0.6:5001.json", "./logs/log.10.0.0.7:5002.json", "./logs/log.10.0.0.8:5003.json"}
	for _, logFile := range logFiles {
		fmt.Println()
		fmt.Println("=============")
		fmt.Println("Reading:", logFile)
		fmt.Println()
		file, err := os.Open(logFile)
		if err != nil {
			panic(err)
		}
		started := 0
		stopped := 0
		bProcessors := make(map[uint64]struct{})
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var entry logEntry
			json.Unmarshal(scanner.Bytes(), &entry)
			if entry.BroadcastID == broadcastID {
				fmt.Println("msg:", entry.Msg, "method:", entry.Method, "from:", entry.From)
			}
			if broadcastID == 2 && entry.Msg == "processor: started" {
				if _, ok := bProcessors[entry.BroadcastID]; ok {
					panic("duplicate broadcastIDs")
				}
				bProcessors[entry.BroadcastID] = struct{}{}
				started++
			}
			if broadcastID != 1 {
				continue
			}
			if entry.Msg == "processor: started" {
				bProcessors[entry.BroadcastID] = struct{}{}
				started++
			}
			if entry.Msg == "processor: stopped" {
				delete(bProcessors, entry.BroadcastID)
				stopped++
			}
		}
		fmt.Println()
		fmt.Println("processors:", bProcessors)
		fmt.Println("processors started:", started)
		fmt.Println("processors stopped:", stopped)
	}
}
