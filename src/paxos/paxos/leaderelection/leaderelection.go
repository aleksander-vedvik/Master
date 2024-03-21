package leaderelection

import (
	"paxos/failuredetector"
	pb "paxos/proto"
)

type process struct {
	rank      uint32
	suspected bool
	addr      string
}

type MonLeader struct {
	*failuredetector.FailureDetector
	processes     map[uint32]*process
	currentLeader *process
	config        *pb.Configuration
	leaderChan    chan string
	doneChan      chan struct{}
}

func New(c *pb.Configuration, id uint32) *MonLeader {
	if c == nil {
		panic("config is nil")
	}
	return &MonLeader{
		config:          c,
		FailureDetector: failuredetector.New(c, id),
		leaderChan:      make(chan string, 10),
		doneChan:        make(chan struct{}),
	}
}

func (l *MonLeader) StartLeaderElection() {
	nodes := l.config.Nodes()
	p := make(map[uint32]*process, len(nodes))
	for _, node := range nodes {
		p[node.ID()] = &process{
			rank:      node.ID(),
			suspected: false,
			addr:      node.Address(),
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
		l.elect()
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

func (l *MonLeader) elect() {
	leader := l.maxRank()
	if leader == nil {
		l.leaderChan <- ""
		return
	}
	if l.currentLeader == nil || leader.rank != l.currentLeader.rank {
		l.currentLeader = leader
		l.leaderChan <- leader.addr
	}
}

func (l *MonLeader) maxRank() *process {
	max := 1
	var p *process = nil
	for _, process := range l.processes {
		if process.suspected {
			continue
		}
		if -int(process.rank) < max {
			max = -int(process.rank)
			p = process
		}
	}
	return p
}
