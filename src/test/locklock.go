package main

import (
	"sync"
	"time"
)

type locklock struct {
	mut      sync.Mutex
	elements map[uint64]*locklockElem
}

func (ll *locklock) handle(m msg) func() int {
	ll.mut.Lock()
	if e, ok := ll.elements[m.id]; ok {
		ll.mut.Unlock()
		if m.fastPath {
			return func() int {
				return m.val
			}
		}
		return func() int {
			return e.handle(m)
		}
	}
	e := &locklockElem{}
	ll.elements[m.id] = e
	ll.mut.Unlock()
	return func() int {
		return e.handle(m)
	}
}

type locklockElem struct {
	mut sync.Mutex
	val int
}

func (e *locklockElem) handle(m msg) int {
	e.mut.Lock()
	defer e.mut.Unlock()
	time.Sleep(tProcessor)
	e.val = m.val
	return m.val
}
