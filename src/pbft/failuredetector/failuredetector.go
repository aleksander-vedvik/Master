package failuredetector

import (
	"context"
	"sync"
	"time"

	pb "pbft/protos"

	"github.com/relab/gorums"
)

type process struct {
	id          uint32
	alive       bool
	suspected   bool
	suspectChan chan<- uint32
	restoreChan chan<- uint32
}

type FailureDetector struct {
	id             uint32
	mu             sync.Mutex
	processes      map[uint32]*process
	config         *pb.Configuration
	delta          time.Duration
	sendHeartbeats func(context.Context, *pb.Heartbeat, ...gorums.CallOption)
	suspectChan    chan uint32
	restoreChan    chan uint32
	doneChan       chan struct{}
}

func New(c *pb.Configuration, id uint32) *FailureDetector {
	suspectChan := make(chan uint32, 10)
	restoreChan := make(chan uint32, 10)
	return &FailureDetector{
		id:             id,
		config:         c,
		sendHeartbeats: c.Ping,
		suspectChan:    suspectChan,
		restoreChan:    restoreChan,
		doneChan:       make(chan struct{}),
		delta:          3 * time.Second,
	}
}

func (fd *FailureDetector) StartFailuredetector() {
	ids := fd.config.NodeIDs()
	p := make(map[uint32]*process, len(ids))
	for _, id := range ids {
		p[id] = &process{
			id:          id,
			alive:       true,
			suspected:   false,
			suspectChan: fd.suspectChan,
			restoreChan: fd.restoreChan,
		}
	}
	fd.processes = p
	go func() {
		for {
			select {
			case <-fd.doneChan:
				return
			default:
			}
			fd.timeout()
		}
	}()
}

func (fd *FailureDetector) StopFailuredetector() {
	close(fd.doneChan)
}

func (fd *FailureDetector) Suspects() <-chan uint32 {
	return fd.suspectChan
}

func (fd *FailureDetector) Restores() <-chan uint32 {
	return fd.restoreChan
}

func (fd *FailureDetector) timeout() {
	fd.mu.Lock()
	for _, p := range fd.processes {
		if !p.alive && !p.suspected {
			p.suspected = true
			p.suspect()
		} else if p.alive && p.suspected {
			p.suspected = false
			p.restore()
		}
		p.alive = false
	}
	fd.mu.Unlock()
	req := &pb.Heartbeat{
		Id: fd.id,
	}
	fd.sendHeartbeats(context.Background(), req, gorums.WithNoSendWaiting())
	time.Sleep(fd.delta)
}

func (p *process) suspect() {
	p.suspectChan <- p.id
}

func (p *process) restore() {
	p.restoreChan <- p.id
}

func (fd *FailureDetector) Ping(id uint32) {
	fd.mu.Lock()
	defer fd.mu.Unlock()
	if p, ok := fd.processes[id]; ok {
		p.alive = true
	}
}
