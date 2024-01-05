package leaderdetector

import (
	"distributedleader/src/defs"
	fd "distributedleader/src/failuredetector"
	"log"
	"time"
)

type LeaderDetector struct {
	id        int
	nodeIds   []int
	suspected map[int]bool
	leader    int

	fd          *fd.EvtFailureDetector
	suspectChan chan int
	restoreChan chan int
}

func NewLeaderDetector(id int, nodeIDs []int, nodeAddresses map[int]string) *LeaderDetector {
	suspect, restore := make(chan int, 1), make(chan int, 1)
	sr := defs.NewSr(suspect, restore)
	failuredetector := fd.NewEvtFailureDetector(id, nodeIDs, nodeAddresses, sr, time.Second)

	ld := LeaderDetector{
		id:        id,
		nodeIds:   nodeIDs,
		suspected: make(map[int]bool),

		fd:          failuredetector,
		suspectChan: suspect,
		restoreChan: restore,
	}
	ld.leader = ld.maxRank()
	return &ld
}

func (ld *LeaderDetector) Start() {
	go ld.fd.Start()
	for {
		select {
		case nodeId := <-ld.suspectChan:
			ld.Suspect(nodeId)
		case nodeId := <-ld.restoreChan:
			ld.Restore(nodeId)
		}
	}
}

func (ld *LeaderDetector) Suspect(nodeId int) {
	log.Println("Suspected:", nodeId)
	ld.suspected[nodeId] = true
	if nodeId == ld.leader {
		ld.changeLeader()
	}
}

func (ld *LeaderDetector) Restore(nodeId int) {
	log.Println("Restored:", nodeId)
	s, ok := ld.suspected[nodeId]
	suspected := s && ok
	if suspected {
		ld.suspected[nodeId] = false
		delete(ld.suspected, nodeId) // Optional
	}
	if ld.leader < nodeId {
		ld.changeLeader()
	}
}

func (ld *LeaderDetector) maxRank() int {
	maxRank := -1
	for _, nodeId := range ld.nodeIds {
		s, ok := ld.suspected[nodeId]
		suspected := s && ok
		if (nodeId > maxRank) && !suspected {
			maxRank = nodeId
		}
	}
	return maxRank
}

func (ld *LeaderDetector) changeLeader() {
	old := ld.leader
	ld.leader = ld.maxRank()
	log.Printf("Changed leader from %v to %v\n", old, ld.leader)
}
