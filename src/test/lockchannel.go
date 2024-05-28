package main

import (
	"context"
	"sync"
	"time"
)

type lockchannel struct {
	//mut      sync.Mutex
	mut      sync.RWMutex
	elements map[uint64]*lockchannelElem
	ctx      context.Context
}

/*func (lc *lockchannel) handle(m msg, buffer int) {
	lc.mut.Lock()
	if e, ok := lc.elements[m.id]; ok {
		lc.mut.Unlock()
		if m.fastPath {
			m.receiveChan <- m.val
			return
		}
		e.sendChan <- m
	} else {
		e := &lockchannelElem{
			sendChan: make(chan msg, buffer),
			ctx:      lc.ctx,
		}
		lc.elements[m.id] = e
		lc.mut.Unlock()
		go e.handle(m)
	}
}*/

func (lc *lockchannel) handle(m msg, buffer int) {
	lc.mut.RLock()
	if e, ok := lc.elements[m.id]; ok {
		lc.mut.RUnlock()
		if m.fastPath {
			m.receiveChan <- m.val
			return
		}
		e.sendChan <- m
	} else {
		lc.mut.RUnlock()
		lc.mut.Lock()
		if e, ok := lc.elements[m.id]; ok {
			lc.mut.Unlock()
			e.sendChan <- m
			return
		}
		e := &lockchannelElem{
			sendChan: make(chan msg, buffer),
			ctx:      lc.ctx,
		}
		lc.elements[m.id] = e
		lc.mut.Unlock()
		go e.handle(m)
	}
}

type lockchannelElem struct {
	sendChan chan msg
	ctx      context.Context
	val      int
}

func (e *lockchannelElem) handle(m msg) {
	time.Sleep(tProcessor)
	e.val = m.val
	m.receiveChan <- m.val
	for {
		select {
		case m := <-e.sendChan:
			time.Sleep(tProcessor)
			e.val = m.val
			m.receiveChan <- m.val
		case <-e.ctx.Done():
			return
		}
	}
}
