package msg

import (
	"sync/atomic"
)

type Subscription interface {
	Unsub()
}

type _subscription struct {
	r     *Router
	topic string
	id    int64
}

func (s *_subscription) Unsub() {
	s.r.unsub(s.topic, s.id)
}

type subscriptions struct {
	changeId atomic.Int64
	subs     map[int64]any
	last     any // The last published value
}

func (s *subscriptions) add(r *Router, topic string, fn any) Subscription {
	id := s.changeId.Add(1)
	s.subs[id] = fn
	return &_subscription{r: r, topic: topic, id: id}
}
