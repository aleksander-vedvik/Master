package failuredetector

import (
	"context"
	"time"

	pb "pbft/protos"

	"log/slog"

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
	processes      map[uint32]*process
	config         *pb.Configuration
	delta          time.Duration
	sendHeartbeats func(context.Context, *pb.HeartbeatReq, ...gorums.CallOption)
	suspectChan    chan uint32
	restoreChan    chan uint32
	doneChan       chan struct{}
}

func New(c *pb.Configuration) *FailureDetector {
	suspectChan := make(chan uint32, 10)
	restoreChan := make(chan uint32, 10)
	return &FailureDetector{
		config:         c,
		sendHeartbeats: c.Heartbeat,
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
	req := &pb.HeartbeatReq{
		From: fd.id,
	}
	fd.sendHeartbeats(context.Background(), req)
	time.Sleep(fd.delta)
}

func (p *process) suspect() {
	p.suspectChan <- p.id
}

func (p *process) restore() {
	p.restoreChan <- p.id
}

func (fd *FailureDetector) Heartbeat(request *pb.HeartbeatReq) {
	slog.Info("failure detector heartbeat")
	if p, ok := fd.processes[request.From]; ok {
		p.alive = true
	}
}
