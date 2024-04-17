package server

import (
	"context"
	"log/slog"
	pb "paxos/proto"
	"sync"
	"time"

	"github.com/relab/gorums"
)

type Proposer struct {
	id          uint32
	peers       []string
	rnd         uint32
	mut         sync.Mutex
	ctx         context.Context
	stopFunc    context.CancelFunc
	view        *pb.Configuration
	clientReqs  []*clientReq
	slots       map[uint32]*pb.PromiseSlot // slots: is the internal data structure maintained by the acceptor to remember the slots
	maxSeenSlot uint32
	broadcast   func(req *pb.AcceptMsg, broadcastID uint64, opts ...gorums.BroadcastOption)
}

func NewProposer(id uint32, peers []string, rnd uint32, view *pb.Configuration, broadcast func(req *pb.AcceptMsg, broadcastID uint64, opts ...gorums.BroadcastOption)) *Proposer {
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
	//p.runPhaseTwo()
}

func (p *Proposer) Stop() {
	p.stopFunc()
}

func (p *Proposer) runPhaseOne() {
	p.mut.Lock()
start:
	select {
	case <-p.ctx.Done():
		p.mut.Unlock()
		slog.Info("phase one: stopping...")
		return
	default:
	}
	//slog.Info("phase one: starting...")
	p.mut.Unlock()
	p.setNewRnd()
	rnd := p.rnd
	maxSeenSlot := p.maxSeenSlot
	prepareCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	promiseMsg, err := p.view.Prepare(prepareCtx, &pb.PrepareMsg{
		Rnd:  rnd,
		Slot: maxSeenSlot,
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
	p.mut.Lock()
	maxSlot := p.maxSeenSlot
	for _, slot := range promiseMsg.Slots {
		if slot.Slot > maxSlot {
			maxSlot = slot.Slot
		}
		if s, ok := p.slots[slot.Slot]; ok {
			if s.Final {
				continue
			}
		}
		p.slots[slot.Slot] = slot
	}
	p.maxSeenSlot = maxSlot
	//slog.Info("phase one: finished")
	p.mut.Unlock()
}

func (p *Proposer) runPhaseTwo() {
	slog.Info("phase two: started...")
	defer slog.Info("phase two: finished")
	for {
		select {
		case <-p.ctx.Done():
			return
		default:
		}
		p.mut.Lock()
		for _, req := range p.clientReqs {
			p.maxSeenSlot++
			p.broadcast(&pb.AcceptMsg{
				Rnd:  p.rnd,
				Slot: p.maxSeenSlot,
				Val:  req.message,
			}, req.broadcastID)
		}
		p.clientReqs = make([]*clientReq, 0)
		p.mut.Unlock()
	}
}

func (p *Proposer) setNewRnd() {
	p.mut.Lock()
	defer p.mut.Unlock()
	numSrvs := uint32(len(p.peers))
	p.rnd -= p.rnd % numSrvs
	p.rnd += p.id + numSrvs
}
