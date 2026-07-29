package hub

import (
	"cmp"
	"fmt"
	"io"
	"iter"
	"slices"

	"github.com/hackborn/onefunc/hub/internal/dep"
)

type Services interface {
	All() iter.Seq2[string, any]
	Close() error
}

func Get[T any](_s Services, name string) (T, error) {
	if s, ok := _s.(*_services); !ok {
		var t T
		return t, fmt.Errorf("Invalid services")
	} else if service, ok := s.all[name].service.(T); !ok {
		var t T
		return t, fmt.Errorf("No service named %v type %T", name, t)
	} else {
		return service, nil
	}
}

func All[T any](all Services) iter.Seq2[string, T] {
	return func(yield func(string, T) bool) {
		for k, v := range all.All() {
			if s, ok := v.(T); ok {
				if !yield(k, s) {
					return
				}
			}
		}
	}
}

type _services struct {
	all map[string]serviceEntry
}

func (s *_services) All() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		for k, v := range s.all {
			if !yield(k, v.service) {
				return
			}
		}
	}
}

func (s *_services) Close() error {
	close := make([]closeEntry, 0, len(s.all))
	for name, a := range s.all {
		closing, _ := a.service.(Closing)
		closer, _ := a.service.(io.Closer)
		close = append(close, closeEntry{name: name,
			dependencies: a.depedencies,
			closing:      closing,
			closer:       closer,
			order:        a.closingOrder})
	}
	// Ignore the error. It should panic in dev, log in
	// production, but not sure the best way to handle that.
	close, _ = dep.Sort(close)
	// Reverse, because the whole point is to do our work
	// before anyone we depend on.
	slices.Reverse(close)

	// Give everyone a chance to perform a pre-Close
	// while all the services are running.
	for _, e := range close {
		if e.closing != nil {
			e.closing.Closing(s)
		}
	}
	// Close on all interested services.
	var err error
	for _, e := range close {
		if e.closer != nil {
			err = cmp.Or(err, e.closer.Close())
		}
	}

	return err
}

func (a *_services) moveTo(key string, b *_services) {
	if s, ok := a.all[key]; ok {
		delete(a.all, key)
		b.all[key] = s
	}
}

func newServices() *_services {
	all := make(map[string]serviceEntry)
	return &_services{all: all}
}

type serviceEntry struct {
	service      any
	depedencies  []string
	closingOrder int
}

type closeEntry struct {
	name         string
	dependencies []string
	closing      Closing
	closer       io.Closer
	order        int
}

func (e closeEntry) Name() string {
	return e.name
}

func (e closeEntry) Dependencies() []string {
	return e.dependencies
}
