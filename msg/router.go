package msg

import (
	"sync/atomic"

	ofstrings "github.com/hackborn/onefunc/strings"
	"github.com/hackborn/onefunc/sync"
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
//
// The current implementation is not designed to be highly optimized.
// The main optimization provided is the Channel, for clients that
// will repeatedly publish to the same topic. A Channel is pretty
// much just a straight function call on all handlers, so if subscription
// changes are minimal and everyone uses Channels then performance should
// be very good, but if you need frequent changes to subscriptions then
// performance will degrade.
//
// The current implementation is curiously partially-thread safe.
// The router itself is safe, but clients would need to correctly
// handle thread safety if they actually sent/handled messages across
// threads (and if you did that, you would want to use a chan). So
// in practice this should only be used in a single thread, and I might
// remove the locking at a future point.
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
	defer sync.Write(&r.mut).Unlock()
	r.r.Edit(pattern, fn)
	return sub
}

func (r *Router) unsub(pattern string, id int64) {
	fn := func(n int64, subs *routerSubscriptions) {
		subs.remove(id)
	}
	defer sync.Write(&r.mut).Unlock()
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
	defer sync.Read(&r.mut).Unlock()
	r.r.Visit(topic, fn)
}

func (r *Router) subsInit(pattern string, subs *routerSubscriptions) {
	subs.r = r
	subs.pattern = pattern
	if subs.subs == nil {
		subs.subs = make(map[int64]any)
	}
}

func mqttMatch(pattern, topic string) bool {
	return ofstrings.MqttMatch(pattern, topic)
}
