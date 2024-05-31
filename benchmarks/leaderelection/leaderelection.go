package leaderelection

import (
	fd "github.com/aleksander-vedvik/benchmark/leaderelection/failuredetector"
)

type process struct {
	rank      uint32
	suspected bool
	addr      string
}

type MonLeader[U fd.Node, T fd.Config[U, V], V fd.Heartbeat] struct {
	*fd.FailureDetector[U, T, V]
	processes     map[uint32]*process
	currentLeader *process
	leaderChan    chan string
	doneChan      chan struct{}
}

func New[U fd.Node, T fd.Config[U, V], V fd.Heartbeat](c T, id uint32, createHeartbeat func(id uint32) V) *MonLeader[U, T, V] {
	if len(c.NodeIDs()) <= 0 {
		panic("config is nil")
	}
	return &MonLeader[U, T, V]{
		FailureDetector: fd.New(c, id, createHeartbeat),
		leaderChan:      make(chan string, 10),
		doneChan:        make(chan struct{}),
	}
}

func (l *MonLeader[U, T, V]) StartLeaderElection() {
	nodes := l.Config.Nodes()
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

func (l *MonLeader[U, T, V]) StopLeaderElection() {
	close(l.doneChan)
	l.FailureDetector.StopFailuredetector()
}

func (l *MonLeader[U, T, V]) Leaders() <-chan string {
	return l.leaderChan
}

func (l *MonLeader[U, T, V]) election() {
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

func (l *MonLeader[U, T, V]) elect() {
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

func (l *MonLeader[U, T, V]) maxRank() *process {
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
