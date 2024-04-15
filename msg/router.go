package msg

import (
	"sync"

	"github.com/hackborn/onefunc/lock"
)

// Subscribe to the topic with the given function. Answer
// the subscription. Use the subscription to unsubscribe.
// The last message publshed to the topic will be immediately
// sent to the function.
func Sub[T any](r *Router, topic string, fn HandlerFunc[T]) Subscription {
	subs, last := r.sub(topic, fn)
	if lastT, ok := last.(T); ok {
		fn(topic, lastT)
	}
	return subs
}

func Pub[T any](r *Router, topic string, value T) {
	if r == nil {
		return
	}
	r.pub(topic, value, func(a any) {
		if c, ok := a.(HandlerFunc[T]); ok {
			c(topic, value)
		}
	})
}

type Router struct {
	mut sync.Mutex
	all map[string]*subscriptions
}

func (r *Router) sub(topic string, value any) (Subscription, any) {
	defer lock.Locker(&r.mut).Unlock()
	r.validate()
	subs := r.validateSubscriptions(topic)
	return subs.add(r, topic, value), subs.last
}

func (r *Router) unsub(topic string, id int64) {
	defer lock.Locker(&r.mut).Unlock()
	if r.all != nil {
		if subs := r.all[topic]; subs != nil {
			delete(subs.subs, id)
		}
	}
}

func (r *Router) pub(topic string, value any, fn visitFunc) {
	defer lock.Locker(&r.mut).Unlock()
	subs := r.validateSubscriptions(topic)
	subs.last = value
	for _, s := range subs.subs {
		fn(s)
	}
}

func (r *Router) validate() {
	if r.all == nil {
		r.all = make(map[string]*subscriptions)
	}
}

func (r *Router) validateSubscriptions(topic string) *subscriptions {
	if subs := r.all[topic]; subs != nil {
		return subs
	}
	all := make(map[int64]any)
	subs := &subscriptions{subs: all}
	r.all[topic] = subs
	return subs
}
