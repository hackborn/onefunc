package strings

import (
	"io"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/hackborn/onefunc/errors"
	"github.com/hackborn/onefunc/maps"
)

// ---------------------------------------------------------
// Pool interface describes an object that gives and
// receives ownership of string writers.
type Pool interface {
	// Acquire takes a writer out of the pool.
	Get(oferrors.Block) io.StringWriter

	// Release puts a writer in the pool.
	Put(io.StringWriter)

	// Close puts all writers back in the master pool.
	Close() error
}

func newLockingPool() *maps.Pool[uint64, io.StringWriter] {
	lock := &sync.Mutex{}
	interner := &stringPoolInterner{}
	return maps.NewPool[uint64, io.StringWriter](lock, interner)
}

type stringPoolInterner struct {
	count atomic.Uint64
}

func (p *stringPoolInterner) Key(w io.StringWriter) (uint64, bool) {
	if b, ok := w.(*stringBuilder); ok {
		return b.id, true
	}
	return 0, false
}

func (p *stringPoolInterner) New() io.StringWriter {
	id := p.count.Add(1)
	var sb strings.Builder
	return &stringBuilder{id: id, sb: &sb}
}

func (p *stringPoolInterner) OnGet(w io.StringWriter, eb oferrors.Block) {
	if b, ok := w.(*stringBuilder); ok {
		b.Reset()
		b.eb = eb
		if b.eb == nil {
			b.eb = nullErrorBlock
		}
	}
}

func (p *stringPoolInterner) OnPut(w io.StringWriter) {
	if b, ok := w.(*stringBuilder); ok {
		b.Reset()
	}
}

// GetWriter removes and answers a new string writer from the global pool.
func GetWriter(eb oferrors.Block) io.StringWriter {
	return globalPool.Get(eb)
}

// PutWriter places a writer into the global pool.
func PutWriter(w io.StringWriter) {
	globalPool.Put(w)
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
	eb  oferrors.Block
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
	b.eb.Add(err)
	return n, err
}

func (b *stringBuilder) Reset() {
	b.sb.Reset()
	b.err = nil
	b.eb = nil
}

// ---------------------------------------------------------
// CONST and VAR

var globalPool = newLockingPool()

var nullErrorBlock = &oferrors.NullBlock{}
