package msg

import (
	"sync"
	"sync/atomic"

	"github.com/hackborn/onefunc/lock"
	ofstrings "github.com/hackborn/onefunc/strings"
)

func NewRouter() *Router {
	match := mqttMatch
	retained := newRetained(match)
	r := &Router{retained: retained}
	r.r.Init = r.subsInit
	r.r.Match = match
	return r
}

// Router provides message routing. Never instantiate it directly,
// use NewRouter() instead, which performs setup.
type Router struct {
	added    atomic.Int64
	deleted  atomic.Int64
	mut      sync.RWMutex
	r        MatchRouter[routerSubscriptions]
	retained *retained
}

// Subscribe to the handlerfunc. Note that clients can
// not call this function directly, they must go through Sub
// so that the function is properly recognized.
func (r *Router) sub(pattern string, value any) Subscription {
	var sub Subscription
	fn := func(n int64, subs *routerSubscriptions) {
		r.added.Add(1)
		sub = subs.add(value)
	}
	defer lock.Write(&r.mut).Unlock()
	r.r.Edit(pattern, fn)
	return sub
}

func (r *Router) unsub(pattern string, id int64) {
	fn := func(n int64, subs *routerSubscriptions) {
		subs.remove(id)
	}
	defer lock.Write(&r.mut).Unlock()
	if len(r.r.patterns) < 1 {
		return
	}
	r.r.Edit(pattern, fn)
}

func (r *Router) retain(topic string, value any) {
	if r.retained != nil {
		r.retained.Retain(topic, value)
	}
}

func (r *Router) visitRetained(pattern string, fn retainedVisitFunc) {
	if r.retained != nil {
		r.retained.Visit(pattern, fn)
	}
}

func (r *Router) readVisit(topic string, fn visitFunc[routerSubscriptions]) {
	defer lock.Read(&r.mut).Unlock()
	r.r.Visit(topic, fn)
}

func (r *Router) subsInit(pattern string, subs *routerSubscriptions) {
	subs.r = r
	subs.pattern = pattern
}

func mqttMatch(pattern, topic string) bool {
	return ofstrings.MqttMatch(pattern, topic)
}
