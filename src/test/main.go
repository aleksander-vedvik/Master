package main

import (
	"time"
)

var tProcessor = 5 * time.Microsecond

type msg struct {
	receiveChan chan int
	val         int
	id          uint64
	fastPath    bool
}
