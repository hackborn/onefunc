package background

import (
	"sync/atomic"
)

// This file provides a menchanism for clients to add their own
// background routines to the service, managed by the service.
// This can only be done during the Open() stage of the
// service lifecycle.

// OpenerService is only available during Open(). It allows clients
// to add managed background tasks.
type OpenerService interface {
	// Go will run the func on the background thread. Closer
	// will be closed when the service ends, and the func
	// should check the closed state and return when closed.
	Go(func(), Closer)
}

type Closer interface {
	Close()
	IsClosed() bool
}

// ChannelCloser is a convenience for the most common case:
// Clients have a channel and ends when it closes.
type ChannelCloser[T any] struct {
	closed atomic.Bool
	c      chan T
}

func NewChannelCloser[T any](c chan T) *ChannelCloser[T] {
	return &ChannelCloser[T]{c: c}
}

func (c *ChannelCloser[T]) Close() {
	c.closed.Store(true)
	close(c.c)
}

func (c *ChannelCloser[T]) IsClosed() bool {
	// The closed state is used so clients don't have to
	// wait for the channel to drain before they know to
	// return. They can if they want, but most won't.
	return c.c == nil || c.closed.Load() == true
}
