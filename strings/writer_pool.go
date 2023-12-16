package strings

import (
	"io"
	"strings"
	"sync"
	"sync/atomic"
)

// ---------------------------------------------------------
// Pool interface describes an object that gives and
// receives ownership of string writers.
type Pool interface {
	// Acquire takes a writer out of the pool.
	Acquire() io.StringWriter

	// Release puts a writer in the pool.
	Release(io.StringWriter)

	// Close puts all writers back in the master pool.
	Close() error
}

// ---------------------------------------------------------
// allocStringWriter can create new string writers.
type allocStringWriter interface {
	New() io.StringWriter
}

// ---------------------------------------------------------
// pool gives and receives io.StringWriters.
type pool struct {
	lock     sync.Locker
	alloc    allocStringWriter
	builders map[uint64]*stringBuilder
}

func (p *pool) Acquire() io.StringWriter {
	sb := p.acquireLocked()
	if sb == nil {
		return p.alloc.New()
	}
	sb.Reset()
	return sb
}

func (p *pool) acquireLocked() *stringBuilder {
	p.lock.Lock()
	defer p.lock.Unlock()

	for k, v := range p.builders {
		delete(p.builders, k)
		return v
	}
	return nil
}

func (p *pool) Release(w io.StringWriter) {
	if b, ok := w.(*stringBuilder); ok {
		p.lock.Lock()
		p.builders[b.id] = b
		p.lock.Unlock()
	}
}

// newLockingPool returns a new stand-alone pool that
// locks and allocates new writers.
func newLockingPool() *pool {
	lock := &sync.Mutex{}
	alloc := &newPoolAllocator{}
	builders := make(map[uint64]*stringBuilder)
	return &pool{lock: lock, alloc: alloc, builders: builders}
}

// newPoolAllocator is used to create new builders.
type newPoolAllocator struct {
	count atomic.Uint64
}

func (a *newPoolAllocator) New() io.StringWriter {
	id := a.count.Add(1)
	var sb strings.Builder
	return &stringBuilder{id: id, sb: &sb}
}

// AcquireWriter removes and answers a new string writer from the global pool.
func AcquireWriter() io.StringWriter {
	return globalPool.Acquire()
}

// ReleaseWriter places a writer into the global pool.
func ReleaseWriter(w io.StringWriter) {
	globalPool.Release(w)
}

func String(w io.StringWriter) string {
	if b, ok := w.(*stringBuilder); ok {
		return b.String()
	}
	return ""
}

func StringErr(w io.StringWriter) error {
	if b, ok := w.(*stringBuilder); ok {
		return b.err
	}
	return nil
}

// stringBuilder is a small convenience on strings.Builder that
// tracks any errors generated.
type stringBuilder struct {
	id  uint64
	sb  *strings.Builder
	err error
}

func (b *stringBuilder) String() string {
	return b.sb.String()
}

func (b *stringBuilder) WriteString(s string) (int, error) {
	n, err := b.sb.WriteString(s)
	if b.err == nil {
		b.err = err
	}
	return n, err
}

func (b *stringBuilder) Reset() {
	b.sb.Reset()
	b.err = nil
}

// ---------------------------------------------------------
// CONST and VAR

var globalPool = newLockingPool()
