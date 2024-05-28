package main

import (
	"context"
	"sync"
	"testing"
)

var p = manyReqsFewMsgsOrder

type params struct {
	numReqs      int
	shardBuffer  int
	sendBuffer   int
	numMsgs      int
	withFastPath bool
}

// LC, LL
var fewReqsManyMsgsNoOrder = params{
	numReqs:      10,
	shardBuffer:  500,
	sendBuffer:   20,
	numMsgs:      2000,
	withFastPath: true,
}

// LC
var fewReqsManyMsgsOrder = params{
	numReqs:      10,
	shardBuffer:  500,
	sendBuffer:   20,
	numMsgs:      2000,
	withFastPath: false,
}

// LL, LC
var manyReqsFewMsgsNoOrder = params{
	numReqs:      500,
	shardBuffer:  500,
	sendBuffer:   20,
	numMsgs:      10,
	withFastPath: true,
}

// LL, LC
var manyReqsFewMsgsOrder = params{
	numReqs:      10000,
	shardBuffer:  500,
	sendBuffer:   20,
	numMsgs:      20,
	withFastPath: false,
}

// LL, (LC)
// all fast
var moderateReqsFewMsgsNoOrder = params{
	numReqs:      100,
	shardBuffer:  500,
	sendBuffer:   20,
	numMsgs:      20,
	withFastPath: true,
}

// CC, LC
// all fast
var moderateReqsFewMsgsOrder = params{
	numReqs:      100,
	shardBuffer:  500,
	sendBuffer:   20,
	numMsgs:      20,
	withFastPath: false,
}

func BenchmarkChannelChannel(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cc := &channelchannel{
		elements: make(map[uint64]*channelchannelElem),
		ctx:      ctx,
		sendChan: make(chan msg, p.shardBuffer),
	}
	go cc.handle(p.sendBuffer)

	b.ResetTimer()
	b.Run("None", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			for id := 0; id < p.numReqs; id++ {
				for j := 0; j < p.numMsgs; j++ {
					wg.Add(1)
					go func(id uint64, i, j int) {
						defer wg.Done()
						resChan := make(chan int, 1)
						m := msg{id: id, val: i, receiveChan: resChan, fastPath: p.withFastPath}
						if j%5 == 0 {
							m.fastPath = false
						}
						cc.sendChan <- m
						<-resChan
					}(uint64(id), i, j)
				}
			}
			wg.Wait()
		}
	})
}

func BenchmarkLockChannel(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lc := &lockchannel{
		elements: make(map[uint64]*lockchannelElem),
		ctx:      ctx,
	}

	b.ResetTimer()
	b.Run("Shard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			for id := 0; id < p.numReqs; id++ {
				for j := 0; j < p.numMsgs; j++ {
					wg.Add(1)
					go func(id uint64, i, j int) {
						defer wg.Done()
						resChan := make(chan int, 1)
						m := msg{id: id, val: i, receiveChan: resChan, fastPath: p.withFastPath}
						if j%5 == 0 {
							m.fastPath = false
						}
						lc.handle(m, p.sendBuffer)
						<-resChan
					}(uint64(id), i, j)
				}
			}
			wg.Wait()
		}
	})
}

func BenchmarkChannelLock(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := &channellock{
		elements: make(map[uint64]*channellockElem),
		ctx:      ctx,
		sendChan: make(chan msg, p.shardBuffer),
	}
	go r.handle()

	b.ResetTimer()
	b.Run("Processor", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			for id := 0; id < p.numReqs; id++ {
				for j := 0; j < p.numMsgs; j++ {
					wg.Add(1)
					go func(id uint64, i, j int) {
						defer wg.Done()
						resChan := make(chan int, 1)
						m := msg{id: id, val: i, receiveChan: resChan, fastPath: p.withFastPath}
						if j%5 == 0 {
							m.fastPath = false
						}
						r.sendChan <- m
						<-resChan
					}(uint64(id), i, j)
				}
			}
			wg.Wait()
		}
	})
}

func BenchmarkLockLock(b *testing.B) {
	l := &locklock{
		elements: make(map[uint64]*locklockElem),
	}

	b.ResetTimer()
	b.Run("Both", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			for id := 0; id < p.numReqs; id++ {
				for j := 0; j < p.numMsgs; j++ {
					wg.Add(1)
					go func(id uint64, i, j int) {
						defer wg.Done()
						m := msg{id: id, val: i, fastPath: p.withFastPath}
						if j%5 == 0 {
							m.fastPath = false
						}
						res := l.handle(m)
						_ = res()
					}(uint64(id), i, j)
				}
			}
			wg.Wait()
		}
	})
}
