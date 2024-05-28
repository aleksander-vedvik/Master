package main

import (
	"context"
	"time"
)

type channelchannel struct {
	sendChan chan msg
	elements map[uint64]*channelchannelElem
	ctx      context.Context
}

func (cc *channelchannel) handle(buffer int) {
	for {
		select {
		case m := <-cc.sendChan:
			if e, ok := cc.elements[m.id]; ok {
				if m.fastPath {
					m.receiveChan <- m.val
					continue
				}
				e.sendChan <- m
			} else {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				e := &channelchannelElem{
					sendChan: make(chan msg, buffer),
					ctx:      ctx,
				}
				go e.handle(m)
				cc.elements[m.id] = e
			}
		case <-cc.ctx.Done():
			return
		}
	}
}

type channelchannelElem struct {
	sendChan chan msg
	ctx      context.Context
	val      int
}

func (e *channelchannelElem) handle(m msg) {
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
