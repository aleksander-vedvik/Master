package server

import (
	"context"
	"log/slog"
	"sync"
	"time"

	pb "github.com/aleksander-vedvik/benchmark/paxos/proto"

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
	slots     map[uint32]*pb.PromiseSlot // slots: is the internal data structure maintained by the acceptor to remember the slots
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
	//p.runPhaseTwo()
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
	//slog.Info("phase one: starting...")
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
	//adu := p.adu
	for i, slot := range promiseMsg.Slots {
		//if slot.Slot > adu {
		//adu = slot.Slot
		//}
		//if s, ok := p.slots[slot.Slot]; ok {
		//if s.Final {
		//continue
		//}
		//}
		p.broadcast(&pb.AcceptMsg{
			Rnd:  p.rnd,
			Val:  slot.Value,
			Slot: uint32(i),
		})
	}
	p.adu = uint32(len(promiseMsg.Slots))
	//slog.Info("phase one: finished")
}

//func (p *Proposer) runPhaseTwo() {
//slog.Info("phase two: started...")
//defer slog.Info("phase two: finished")
//for {
//select {
//case <-p.ctx.Done():
//return
//default:
//}
//p.mut.Lock()
//for _, req := range p.clientReqs {
//p.maxSeenSlot++
//p.broadcast(&pb.AcceptMsg{
//Rnd:  p.rnd,
//Slot: p.maxSeenSlot,
//Val:  req.message,
//}, req.broadcastID)
//}
//p.clientReqs = make([]*clientReq, 0)
//p.mut.Unlock()
//}
//}

func (p *Proposer) setNewRnd() {
	//p.mut.Lock()
	//defer p.mut.Unlock()
	numSrvs := uint32(len(p.peers))
	p.rnd -= p.rnd % numSrvs
	p.rnd += p.id + numSrvs
}
