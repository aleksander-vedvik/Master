package main

import (
	"context"
	"sync"
	"time"
)

type channellock struct {
	sendChan chan msg
	elements map[uint64]*channellockElem
	ctx      context.Context
}

func (cl *channellock) handle() {
	for {
		select {
		case m := <-cl.sendChan:
			if e, ok := cl.elements[m.id]; ok {
				if m.fastPath {
					m.receiveChan <- m.val
					continue
				}
				go e.handle2(m, cl.ctx)
			} else {
				e := &channellockElem{}
				cl.elements[m.id] = e
				go e.handle2(m, cl.ctx)
			}
		case <-cl.ctx.Done():
			return
		}
	}
}

type channellockElem struct {
	mut sync.Mutex
	val int
}

func (e *channellockElem) handle2(m msg, ctx context.Context) {
	e.mut.Lock()
	defer e.mut.Unlock()
	time.Sleep(tProcessor)
	e.val = m.val
	select {
	case m.receiveChan <- m.val:
	case <-ctx.Done():
	}
}
