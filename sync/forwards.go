package sync

import (
	"sync"
)

// The sync package is intended to be a drop-in replacement for the
// system sync, just adding some minor quality-of-life improvements.
// NOTE: Obviously this is a fragile design, and is relying on the
// standard library never introducing a conflicting name. Would not
// be surprised if I unwound this decision in the future.
type Cond = sync.Cond
type Locker = sync.Locker
type Map = sync.Map
type Mutex = sync.Mutex
type Pool = sync.Pool
type Once = sync.Once
type RWMutex = sync.RWMutex
type WaitGroup = sync.WaitGroup

func NewCond(l Locker) *Cond {
	return sync.NewCond(l)
}
