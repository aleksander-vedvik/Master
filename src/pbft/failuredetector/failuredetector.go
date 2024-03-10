package failuredetector

import (
	"time"

	pb "pbft/protos"

	"github.com/relab/gorums"
)

type process struct {
	alive    bool
	detected bool
}

type Failuredetector struct {
	processes map[uint32]*process
	delta     time.Duration
	doneChan  chan struct{}
}

func New() *Failuredetector {
	return &Failuredetector{}
}

func (fd *Failuredetector) StartFailuredetector() {
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

func (fd *Failuredetector) StopFailuredetector() {
	close(fd.doneChan)
}

func (fd *Failuredetector) timeout() {
	for _, p := range fd.processes {
		if !p.alive && !p.detected {
			p.detected = true
			p.crash()
		}
		p.sendHeartbeat()
		p.alive = false
	}
	time.Sleep(fd.delta)
}

func (p *process) crash() {

}

func (p *process) sendHeartbeat() {

}

func (fd *Failuredetector) Heartbeat(ctx gorums.ServerCtx, request *pb.PrePrepareRequest, broadcast *pb.Broadcast) {

}
