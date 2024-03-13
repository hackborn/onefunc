package nxt

import (
	"sync/atomic"
)

func init() {
	runnerCounter.Store(0)
	wrapperCounter.Store(1 << 32)
}

var (
	runnerCounter  = &atomic.Int64{}
	wrapperCounter = &atomic.Int64{}
)
