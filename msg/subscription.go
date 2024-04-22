package msg

import (
	"sync/atomic"
)

type Subscription interface {
	Unsub()
}

type _subscription struct {
	r       *Router
	pattern string
	id      int64
}

func (s *_subscription) Unsub() {
	s.r.unsub(s.pattern, s.id)
}

type _subscriptions struct {
	r       *Router
	pattern string
	change  atomic.Int64
	subs    map[int64]any
	last    any // The last published value
}

func (s *_subscriptions) add(fn any) Subscription {
	if s.subs == nil {
		s.subs = make(map[int64]any)
	}
	id := s.change.Add(1)
	s.subs[id] = fn
	return &_subscription{r: s.r, pattern: s.pattern, id: id}
}

func (s *_subscriptions) remove(id int64) {
	if s.subs == nil {
		return
	}
	if _, ok := s.subs[id]; ok {
		delete(s.subs, id)
		s.r.deleted.Add(1)
	}
}
