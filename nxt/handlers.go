package nxt

import (
	"sync"

	"github.com/hackborn/onefunc/lock"
)

func wrapperFunc(f Factory, wd *_wrapperData, next Handler) Handler {
	id := wrapperCounter.Add(1)
	w := _wrapper{id: id,
		next:    next,
		factory: f,
		wd:      wd}
	return w.Handle
}

func newWrapperData() *_wrapperData {
	locker := &sync.Mutex{}
	data := make(map[any]any)
	return &_wrapperData{locker: locker, data: data}
}

type _wrapperData struct {
	locker sync.Locker
	data   map[any]any
}

func (d *_wrapperData) Get(id any) (any, bool) {
	defer lock.Locker(d.locker).Unlock()
	ans, ok := d.data[id]
	return ans, ok
}

func (d *_wrapperData) Set(id any, data any) {
	defer lock.Locker(d.locker).Unlock()
	d.data[id] = data
}

type _wrapper struct {
	id      int64
	next    Handler
	factory Factory
	wd      *_wrapperData
}

func (w *_wrapper) Handle(args HandlerArgs, event any) {
	args.Data = nil
	id := args.runnerId + w.id
	if d, ok := w.wd.Get(id); ok {
		args.Data = d
	} else if w.factory != nil {
		d := w.factory.NewData()
		w.wd.Set(id, d)
		args.Data = d
	}
	w.next(args, event)
}

// NilNewHandlerFunc is used to return a nil handler.
func NilNewHandlerFunc(NewHandlerArgs) Handler {
	return nil
}

// NullHandler is used when an engine isn't supplied a handler.
func NullHandler() Handler {
	return func(args HandlerArgs, event any) {
	}
}

// ManyHandler fans an event to multiple handlers.
func ManyHandler(handlers []Handler) Handler {
	return func(args HandlerArgs, event any) {
		for _, h := range handlers {
			h(args, event)
		}
	}
}
