package msg

import (
	"sync"
	"sync/atomic"

	"github.com/hackborn/onefunc/lock"
)

// Subscribe to the topic with the given function. Answer
// the subscription. Use the subscription to unsubscribe.
// The last message publshed to the topic will be immediately
// sent to the function.
func Sub[T any](r *Router, pattern string, fn HandlerFunc[T]) Subscription {
	subs, last := r.sub(pattern, fn)
	if lastT, ok := last.(T); ok {
		fn(pattern, lastT)
	}
	return subs
}

func Pub[T any](r *Router, topic string, value T) {
	if r == nil {
		return
	}
	//	subsFn := func(subs *subscriptions) {
	//		subs.last = value
	//	}
	// r.readVisit(topic, subsFn, func(a any) {
	r.readVisit(topic, func(pattern string, subs *_subscriptions) {
		for _, _h := range subs.subs {
			if h, ok := _h.(HandlerFunc[T]); ok {
				h(topic, value)
			}
		}
	})
}

type Router struct {
	added   atomic.Int64
	deleted atomic.Int64
	mut     sync.RWMutex
	r       MatchRouter[_subscriptions]
}

// Subscribe to the handlerfunc. Note that clients can
// not call this function directly, they must go through Sub
// so that the function is properly recognized.
func (r *Router) sub(pattern string, value any) (Subscription, any) {
	// Current design lets clients create the router without
	// going through an init func, but I still want initialization.
	r.r.Init = r.subsInit
	r.r.Match = globMatch

	var sub Subscription
	fn := func(n int64, subs *_subscriptions) {
		r.added.Add(1)
		sub = subs.add(value)
	}
	defer lock.Write(&r.mut).Unlock()
	r.r.Edit(pattern, fn)
	return sub, nil
}

func (r *Router) unsub(pattern string, id int64) {
	// Current design lets clients create the router without
	// going through an init func, but I still want initialization.
	r.r.Match = globMatch

	fn := func(n int64, subs *_subscriptions) {
		subs.remove(id)
	}
	defer lock.Write(&r.mut).Unlock()
	if len(r.r.patterns) < 1 {
		return
	}
	r.r.Edit(pattern, fn)
}

func (r *Router) readVisit(topic string, fn visitFunc[_subscriptions]) {
	// Current design lets clients create the router without
	// going through an init func, but I still want initialization.
	r.r.Match = globMatch

	defer lock.Read(&r.mut).Unlock()
	r.r.Visit(topic, fn)
}

func (r *Router) subsInit(pattern string, subs *_subscriptions) {
	subs.r = r
	subs.pattern = pattern
}
