package maps

import (
	"sync"
)

// NewPool answers a new generic pool object.
func NewPool[K comparable, V any](lock sync.Locker,
	interner PoolInterner[K, V]) *Pool[K, V] {
	cache := make(map[K]V)
	return &Pool[K, V]{interner: interner, lock: lock, cache: cache}
}

// PoolInterner specifies the operations necessary
// for a functioning pool. It can serve as an abstraction
// layer for pools that use an external type for their
// public API and an internal type for storage.
type PoolInterner[K comparable, V any] interface {
	// New allocates a new pool object.
	New() V

	// Key answers the key for the pool object.
	Key(v V) (K, bool)

	// Initialize resets the state of the pool object before use.
	Initialize(v V)
}

type Pool[K comparable, V any] struct {
	interner PoolInterner[K, V]
	emptyV   V

	lock  sync.Locker
	cache map[K]V
}

func (p *Pool[K, V]) Acquire() V {
	if v, ok := p.acquireLocked(); ok {
		p.interner.Initialize(v)
		return v
	}
	return p.interner.New()
}

func (p *Pool[K, V]) acquireLocked() (V, bool) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for k, v := range p.cache {
		delete(p.cache, k)
		return v, true
	}
	return p.emptyV, false
}

func (p *Pool[K, V]) Release(v V) {
	if k, ok := p.interner.Key(v); ok {
		p.lock.Lock()
		p.cache[k] = v
		p.lock.Unlock()
	}
}

func (p *Pool[K, V]) Close() error {
	return nil
}

// CacheSize reports the size of my current cache.
// Just for debugging.
func (p *Pool[K, V]) CacheSize() int {
	p.lock.Lock()
	defer p.lock.Unlock()

	return len(p.cache)
}
