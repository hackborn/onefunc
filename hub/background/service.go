package background

import (
	"slices"

	"github.com/hackborn/onefunc/cfg"
	"github.com/hackborn/onefunc/hub"
	"github.com/hackborn/onefunc/sync"
)

func newService(settings cfg.Settings) *service {
	c := make(chan taskEntry, 512)
	l := sync.Mutex{}
	s := &service{c: c, l: &l,
		startUpdateFn: defaultUpdateFunc,
		endUpdateFn:   defaultUpdateFunc,
		runTaskFn:     defaultRunTaskFunc,
	}
	closer := NewChannelCloser[taskEntry](c)
	s.closers = append(s.closers, closer)
	workers := 4
	if settings.MustBool("debug/log", false) {
		timer := &debugTimer{}
		s.startUpdateFn = timer.startUpdate
		s.endUpdateFn = timer.endUpdate
		s.runTaskFn = timer.runLocal
	}
	for range workers {
		s.wg.Go(func() {
			s.loop(closer)
		})
	}
	return s
}

type service struct {
	closers      []Closer
	c            chan taskEntry
	wg           sync.WaitGroup
	required     []*requiredTask
	localHandled []taskEntry

	l         sync.Locker
	handled   []taskEntry
	streaming []taskEntry

	// Debug support. So hacky. Wish I could cleanly separate it.
	startUpdateFn updateFunc
	endUpdateFn   updateFunc
	runTaskFn     runTaskFunc
}

// Closing shuts down before the official Close(), because
// other services might rely on the background service.
func (s *service) Closing(hub.Services) {
	s.cancelAll()
	for _, closer := range s.closers {
		closer.Close()
	}
	s.wg.Wait()
	s.runRequired(s.required)
}

func (s *service) doGo(fn func(), c Closer) {
	if s.c != nil && fn != nil && c != nil {
		s.closers = append(s.closers, c)
		s.wg.Go(fn)
	}
}

func (s *service) Run(task Task) TaskInfo {
	// Required will be applied to everything wrapped by
	// the required; if the required is wrapped in something,
	// the parent gets lost. This is sorta weird, and points
	// to required being a flag, but that also has issues to
	// unwind, since required tracks state.
	if r, ok := getTask[*requiredTask](task); ok {
		s.required = append(s.required, r)
	}
	info := TaskInfo{Cancel: &sync.Cancel{}}
	s.c <- taskEntry{t: task, cancel: info.Cancel}
	return info
}

func (s *service) RunOn(t Task, workerGroup int) TaskInfo {
	// Placeholder
	return s.Run(t)
}

func (s *service) Update() {
	s.startUpdateFn()
	s.updateStreaming()
	s.localHandled = s.popHandled(s.localHandled)
	if len(s.localHandled) < 1 {
		s.endUpdateFn()
		return
	}
	for _, t := range s.localHandled {
		if !t.IsCancelled() {
			s.runTaskFn(t.t)
		}
	}
	s.required = slices.DeleteFunc(s.required, func(t *requiredTask) bool {
		return t.localDone.Load() == true
	})
	s.endUpdateFn()
}

func (s *service) runRequired(tasks []*requiredTask) {
	// Requireds don't support streaming.
	// (and apparently cancelling, though they should,
	// although you could make the argument that defeats
	// being required)
	args := RemoteArgs{}
	for _, t := range tasks {
		if t.remoteDone.Load() == false {
			t.Remote(args)
		}
		if t.localDone.Load() == false {
			t.Local()
		}
	}
}

func (s *service) loop(closer Closer) {
	for {
		task, more := <-s.c
		if !more {
			return
		}
		if !task.IsCancelled() {
			args := RemoteArgs{streamingFn: s.pushStreaming,
				cancel: task.cancel,
			}
			task.t.Remote(args)
		}
		s.pushHandled(task)

		if closer.IsClosed() {
			return
		}
	}
}

func (s *service) updateStreaming() {
	s.localHandled = s.popStreaming(s.localHandled)
	for _, te := range s.localHandled {
		if !te.IsCancelled() {
			s.runTaskFn(te.t)
		}
	}
}

func (s *service) pushHandled(t taskEntry) {
	defer sync.Lock(s.l).Unlock()
	s.handled = append(s.handled, t)
}

func (s *service) popHandled(replace []taskEntry) []taskEntry {
	replace = replace[:0]
	defer sync.Lock(s.l).Unlock()
	replace, s.handled = s.handled, replace
	return replace
}

func (s *service) pushStreaming(t taskEntry) {
	defer sync.Lock(s.l).Unlock()
	s.streaming = append(s.streaming, t)
}

func (s *service) popStreaming(replace []taskEntry) []taskEntry {
	replace = replace[:0]
	defer sync.Lock(s.l).Unlock()
	replace, s.streaming = s.streaming, replace
	return replace
}

// ---------------------------------------------------------
// CANCELLING

func (s *service) cancelAll() {
	defer sync.Lock(s.l).Unlock()
	s.cancelEntries(s.handled)
	s.cancelEntries(s.streaming)
}

func (s *service) cancelEntries(entries []taskEntry) {
	for _, e := range entries {
		if e.cancel != nil {
			e.cancel.Cancel()
		}
	}
}

// ---------------------------------------------------------
// TASK ENTRY

type taskEntry struct {
	t      Task
	cancel *sync.Cancel
}

func (e taskEntry) IsCancelled() bool {
	if e.cancel != nil {
		return e.cancel.IsCancelled()
	}
	return false
}

// ---------------------------------------------------------
// STREAMING

type taskHandlerFunc func(taskEntry)
