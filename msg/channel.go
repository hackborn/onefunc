package msg

type Channel[T any] interface {
	Pub(value T)
}

// NewChannel answers a new channel on the given topic.
// Channels are an optimization when a client will be publishing
// multiple times. It's safe to add or remove subscriptions to
// the topic after creating a channel, the channel will be updated.
func NewChannel[T any](r *Router, topic string) Channel[T] {
	c := &_channel[T]{r: r,
		topic:     topic,
		addedId:   r.added.Load() - 1,
		deletedId: r.deleted.Load() - 1}
	if r == nil {
		return c
	}
	return c
}

type _channel[T any] struct {
	r     *Router
	topic string
	// added and deleted let me know whenever the router state has changed,
	// so I can rebuild my handlers. It'd be nice to be smarter but this works for now.
	addedId   int64
	deletedId int64
	handlers  []HandlerFunc[T]
}

func (c *_channel[T]) Pub(value T) {
	c.sync()
	for _, fn := range c.handlers {
		fn(c.topic, value)
	}
	c.r.retain(c.topic, value)
}

func (c *_channel[T]) sync() {
	added := c.r.added.Load()
	deleted := c.r.deleted.Load()
	if c.addedId == added && c.deletedId == deleted {
		return
	}
	c.addedId = added
	c.deletedId = deleted
	c.handlers = c.handlers[:0]
	visitFn := func(pattern string, subs *routerSubscriptions) {
		for _, _h := range subs.subs {
			if h, ok := _h.(HandlerFunc[T]); ok {
				c.handlers = append(c.handlers, h)
			}
		}
	}
	c.r.readVisit(c.topic, visitFn)
}
