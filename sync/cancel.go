package sync

import (
	"sync/atomic"
)

// Cancel is a simple thread-safe cancel state. It
// just wraps an atomic.Bool, providing a domain-specific
// syntax for clients that want it.
type Cancel struct {
	b atomic.Bool
}

func (c *Cancel) Cancel() {
	c.b.Store(true)
}

func (c *Cancel) IsCancelled() bool {
	return c.b.Load() == true
}

type CancelChecker interface {
	IsCancelled() bool
}

func ChainCancelChecker(a, b CancelChecker) CancelChecker {
	return &checkerChain{a: a, b: b}
}

type checkerChain struct {
	a, b CancelChecker
}

func (c *checkerChain) IsCancelled() bool {
	if c.a != nil && c.a.IsCancelled() {
		return true
	}
	if c.b != nil && c.b.IsCancelled() {
		return true
	}
	return false
}
