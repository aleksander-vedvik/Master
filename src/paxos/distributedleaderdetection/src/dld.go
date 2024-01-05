package main

import (
	"distributedleader/src/config"
	"distributedleader/src/defs"
	fd "distributedleader/src/failuredetector"
	ld "distributedleader/src/leaderdetector"
	"log"
	"time"
)

func main() {
	log.Println("Starting node...")
	ld, err := initLd()
	if err != nil {
		log.Fatalln("Could not start the node. Exiting with error code:", err)
		return
	}
	ld.Start()
	/*suspect, restore := make(chan int, 1), make(chan int, 1)
	fd, err := initFd(suspect, restore)
	if err != nil {
		log.Fatalln("Could not start the node. Exiting with error code:", err)
		return
	}
	go suspectrestorer(suspect, restore)
	fd.Start()*/
}

func initLd() (*ld.LeaderDetector, error) {
	id, config, err := config.ReadConfig("./src/config.json")
	if err != nil {
		return nil, err
	}
	ld := ld.NewLeaderDetector(id, config.GetNodeIds(), config.GetNodeAddressess())
	return ld, nil
}

func initFd(suspect, restore chan int) (*fd.EvtFailureDetector, error) {
	id, config, err := config.ReadConfig("./src/config.json")
	if err != nil {
		return nil, err
	}
	fd := fd.NewEvtFailureDetector(id, config.GetNodeIds(), config.GetNodeAddressess(), defs.NewSr(suspect, restore), time.Second)
	return fd, nil
}

func suspectrestorer(suspect, restore chan int) {
	for {
		select {
		case nodeId := <-suspect:
			log.Println("Suspected:", nodeId)
		case nodeId := <-restore:
			log.Println("Restored:", nodeId)
		}
	}
}
