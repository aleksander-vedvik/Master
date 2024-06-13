package server

import (
	"context"
	"log/slog"
	"sync"
	"time"

	pb "github.com/aleksander-vedvik/benchmark/paxos.bc/proto"

	"github.com/relab/gorums"
)

type Proposer struct {
	mut       sync.Mutex
	id        uint32
	peers     []string
	rnd       uint32
	ctx       context.Context
	stopFunc  context.CancelFunc
	view      *pb.Configuration
	adu       uint32
	broadcast func(req *pb.AcceptMsg, opts ...gorums.BroadcastOption)
}

func NewProposer(id uint32, peers []string, rnd uint32, view *pb.Configuration, broadcast func(req *pb.AcceptMsg, opts ...gorums.BroadcastOption)) *Proposer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Proposer{
		id:        id,
		peers:     peers,
		rnd:       rnd,
		ctx:       ctx,
		stopFunc:  cancel,
		view:      view,
		broadcast: broadcast,
	}
}

func (p *Proposer) Start() {
	p.runPhaseOne()
}

func (p *Proposer) Stop() {
	p.stopFunc()
}

func (p *Proposer) runPhaseOne() {
	p.mut.Lock()
	defer p.mut.Unlock()
start:
	select {
	case <-p.ctx.Done():
		slog.Info("phase one: stopping...")
		return
	default:
	}
	p.setNewRnd()
	prepareCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	promiseMsg, err := p.view.Prepare(prepareCtx, &pb.PrepareMsg{
		Rnd:  p.rnd,
		Slot: 0,
	})
	cancel()
	if err != nil {
		select {
		case <-time.After(5 * time.Second):
			slog.Error("phase one: error...", "err", err)
			goto start
		case <-p.ctx.Done():
			return
		}
	}
	for i, slot := range promiseMsg.Slots {
		p.broadcast(&pb.AcceptMsg{
			Rnd:  p.rnd,
			Val:  slot.Value,
			Slot: uint32(i),
		})
	}
	p.adu = uint32(len(promiseMsg.Slots))
}

func (p *Proposer) setNewRnd() {
	numSrvs := uint32(len(p.peers))
	p.rnd -= p.rnd % numSrvs
	p.rnd += p.id + numSrvs
}
