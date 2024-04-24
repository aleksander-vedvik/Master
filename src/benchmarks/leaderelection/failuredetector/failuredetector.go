package failuredetector

import (
	"context"
	"sync"
	"time"

	"github.com/relab/gorums"
)

type process struct {
	id          uint32
	alive       bool
	suspected   bool
	suspectChan chan<- uint32
	restoreChan chan<- uint32
}

type Heartbeat interface {
	GetId() uint32
}

type Node interface {
	ID() uint32
	Address() string
}

type Config[U Node, V Heartbeat] interface {
	Ping(context.Context, V, ...gorums.CallOption)
	NodeIDs() []uint32
	Nodes() []U
}

type FailureDetector[U Node, T Config[U, V], V Heartbeat] struct {
	id              uint32
	mu              sync.Mutex
	processes       map[uint32]*process
	Config          T
	delta           time.Duration
	createHeartbeat func(id uint32) V
	suspectChan     chan uint32
	restoreChan     chan uint32
	doneChan        chan struct{}
}

func New[U Node, T Config[U, V], V Heartbeat](c T, id uint32, createHeartbeat func(id uint32) V) *FailureDetector[U, T, V] {
	suspectChan := make(chan uint32, 10)
	restoreChan := make(chan uint32, 10)
	return &FailureDetector[U, T, V]{
		id:              id,
		Config:          c,
		createHeartbeat: createHeartbeat,
		suspectChan:     suspectChan,
		restoreChan:     restoreChan,
		doneChan:        make(chan struct{}),
		delta:           3 * time.Second,
	}
}

func (fd *FailureDetector[U, T, V]) StartFailuredetector() {
	ids := fd.Config.NodeIDs()
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

func (fd *FailureDetector[U, T, V]) StopFailuredetector() {
	close(fd.doneChan)
}

func (fd *FailureDetector[U, T, V]) Suspects() <-chan uint32 {
	return fd.suspectChan
}

func (fd *FailureDetector[U, T, V]) Restores() <-chan uint32 {
	return fd.restoreChan
}

func (fd *FailureDetector[U, T, V]) timeout() {
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
	req := fd.createHeartbeat(fd.id)
	ctx, cancel := context.WithTimeout(context.Background(), fd.delta)
	defer cancel()
	fd.Config.Ping(ctx, req, gorums.WithNoSendWaiting())
	time.Sleep(fd.delta)
}

func (p *process) suspect() {
	p.suspectChan <- p.id
}

func (p *process) restore() {
	p.restoreChan <- p.id
}

func (fd *FailureDetector[U, T, V]) Ping(id uint32) {
	fd.mu.Lock()
	defer fd.mu.Unlock()
	if p, ok := fd.processes[id]; ok {
		p.alive = true
	}
}
