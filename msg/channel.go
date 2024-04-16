package msg

import (
	"sync/atomic"

	"github.com/hackborn/onefunc/lock"
)

type Channel[T any] interface {
	Pub(value T)
}

// NewChannel answers a new channel on the given topic.
// Channels are an optimization when a client will be publishing
// multiple times. It's safe to add or remove subscriptions to
// the topic after creating a channel, the channel will be updated.
func NewChannel[T any](r *Router, topic string) Channel[T] {
	c := &_channel[T]{r: r, topic: topic}
	if r == nil {
		return c
	}
	subsFn := func(subs *subscriptions) {
		c.subs = subs
	}
	r.visit(topic, subsFn, nil)
	c.sync()
	return c
}

type _channel[T any] struct {
	r        *Router
	topic    string
	changeId atomic.Int64
	subs     *subscriptions
	fns      []HandlerFunc[T]
}

func (c *_channel[T]) Pub(value T) {
	c.sync()
	for _, fn := range c.fns {
		fn(c.topic, value)
	}
}

func (c *_channel[T]) sync() {
	if c.changeId.Load() == c.subs.changeId.Load() {
		return
	}
	c.fns = c.fns[:0]
	defer lock.Locker(&c.r.mut).Unlock()
	for _, _a := range c.subs.subs {
		if a, ok := _a.(HandlerFunc[T]); ok {
			c.fns = append(c.fns, a)
		}
	}
	c.changeId.Store(c.subs.changeId.Load())
}
