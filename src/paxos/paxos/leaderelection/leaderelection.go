package leaderelection

import (
	"paxos/failuredetector"
	pb "paxos/proto"
)

type process struct {
	rank      uint32
	suspected bool
}

type MonLeader struct {
	*failuredetector.FailureDetector
	processes     map[uint32]*process
	currentLeader *process
	config        *pb.Configuration
	leaderChan    chan string
	doneChan      chan struct{}
}

func New(c *pb.Configuration) *MonLeader {
	if c == nil {
		panic("config is nil")
	}
	return &MonLeader{
		config:          c,
		FailureDetector: failuredetector.New(c),
		leaderChan:      make(chan string, 10),
		doneChan:        make(chan struct{}),
	}
}

func (l *MonLeader) StartLeaderElection() {
	ids := l.config.NodeIDs()
	p := make(map[uint32]*process, len(ids))
	for _, id := range ids {
		p[id] = &process{
			rank:      id,
			suspected: false,
		}
	}
	l.processes = p
	l.FailureDetector.StartFailuredetector()
	go l.election()
}

func (l *MonLeader) StopLeaderElection() {
	close(l.doneChan)
	l.FailureDetector.StopFailuredetector()
}

func (l *MonLeader) Leaders() <-chan string {
	return l.leaderChan
}

func (l *MonLeader) election() {
	var id uint32
	for {
		l.elect(id)
		select {
		case <-l.doneChan:
			return
		case id = <-l.FailureDetector.Suspects():
			l.processes[id].suspected = true
		case id = <-l.FailureDetector.Restores():
			l.processes[id].suspected = false
		}
	}
}

func (l *MonLeader) elect(id uint32) {
	leader := l.maxRank()
	if l.currentLeader == nil || leader.rank != l.currentLeader.rank {
		l.currentLeader = leader
		var addr string
		for _, node := range l.config.Nodes() {
			if node.ID() == id {
				addr = node.Address()
				break
			}
		}
		l.leaderChan <- addr
	}
}

func (l *MonLeader) maxRank() *process {
	max := 1
	for _, process := range l.processes {
		if process.suspected {
			continue
		}
		if -int(process.rank) < max {
			return process
		}
	}
	panic("all processes suspected")
}
