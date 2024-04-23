package msg

import (
	"sync/atomic"
)

func newRetained(match MatchFunc) *retained {
	all := make(map[string]last)
	return &retained{all: all, match: match}
}

type retained struct {
	// all is a map of topic to last published value.
	all   map[string]last
	match MatchFunc
}

func (r *retained) Retain(topic string, value any) {
	last := last{}
	last.value.Store(value)
	r.all[topic] = last
}

func (r *retained) Visit(pattern string, fn retainedVisitFunc) {
	for k, last := range r.all {
		if r.match(pattern, k) {
			if v := last.value.Load(); v != nil {
				fn(k, v)
			}
		}
	}
}

type last struct {
	value atomic.Value // The last published value
}
