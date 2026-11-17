package background

import (
	"slices"

	"github.com/hackborn/onefunc/cfg"
	"github.com/hackborn/onefunc/hub"
	"github.com/hackborn/onefunc/sync"
)

func newService(settings cfg.Settings) *service {
	c := make(chan RunningTask, 512)
	l := sync.Mutex{}
	s := &service{c: c, l: &l,
		startUpdateFn: defaultUpdateFunc,
		endUpdateFn:   defaultUpdateFunc,
		runTaskFn:     defaultRunLocalFunc,
	}
	closer := NewChannelCloser[RunningTask](c)
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
			s.loop(s.c, closer, nil)
		})
	}
	return s
}

type service struct {
	closers      []Closer
	c            chan RunningTask
	channels     []channel
	wg           sync.WaitGroup
	required     []*requiredTask
	localHandled []RunningTask
	// localArgs stores any client data. Ideally it is only accessed
	// in the main thread, but in practice some UI libraries force access
	// in worker threads, so this is immutable after service construction,
	// and any data should be thread-safe.
	localArgs LocalArgs
	// runLocalFn allows clients to replace how the task is run
	// locally. If nil, the task is sent to the main thread and
	// Local() is called in the Update(). If not nil, this function
	// is called in the worker thread that handled the remote, and
	// the task is NOT sent to the main thread. This is here because
	// some libraries have their own method for running in the main
	// thread and discourage the existence of an Update() loop.
	runLocalFn RunLocalFunc

	l         sync.Locker
	handled   []RunningTask
	streaming []RunningTask

	// Debug support. So hacky. Wish I could cleanly separate it.
	startUpdateFn updateFunc
	endUpdateFn   updateFunc
	runTaskFn     RunLocalFunc
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

func (s *service) doAddChannel(args AddChannelArgs) int {
	id := len(s.channels)
	c := channel{c: make(chan RunningTask, 512),
		localArgs: args.LocalData,
	}
	s.channels = append(s.channels, c)
	closer := NewChannelCloser[RunningTask](c.c)
	s.closers = append(s.closers, closer)
	s.wg.Go(func() {
		s.loop(c.c, closer, args.RemoteData)
	})
	return id
}

func (s *service) doGo(fn func(), c Closer) {
	if s.c != nil && fn != nil && c != nil {
		s.closers = append(s.closers, c)
		s.wg.Go(fn)
	}
}

func (s *service) doSendLocalFunc() SendLocalFunc {
	return func(rt RunningTask) {
		if s.runLocalFn == nil {
			s.pushHandled(rt)
		} else {
			s.runLocalFn(rt.task, s.localArgs)
		}
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
	s.c <- RunningTask{task: task, cancel: info.Cancel}
	return info
}

func (s *service) RunOn(task Task, workerGroup int) TaskInfo {
	if workerGroup < 0 || workerGroup >= len(s.channels) {
		return s.Run(task)
	}

	// Required will be applied to everything wrapped by
	// the required; if the required is wrapped in something,
	// the parent gets lost. This is sorta weird, and points
	// to required being a flag, but that also has issues to
	// unwind, since required tracks state.
	if r, ok := getTask[*requiredTask](task); ok {
		s.required = append(s.required, r)
	}
	info := TaskInfo{Cancel: &sync.Cancel{}}
	s.channels[workerGroup].c <- RunningTask{task: task, cancel: info.Cancel}
	return info
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
			s.runTaskFn(t.task, s.localArgs)
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
			t.Local(s.localArgs)
		}
	}
}

func (s *service) loop(c <-chan RunningTask, closer Closer, remoteData any) {
	// Universal arg params
	args := RemoteArgs{streamingFn: s.pushStreaming}
	if remoteData != nil {
		args.data = append(args.data, remoteData)
	}

	for {
		task, more := <-c
		if !more {
			return
		}
		if !task.IsCancelled() {
			args.cancel = task.cancel
			task.task.Remote(args)
		}
		if !task.IsCancelled() {
			if s.runLocalFn == nil {
				s.pushHandled(task)
			} else {
				s.runLocalFn(task.task, s.localArgs)
			}
		}

		if closer.IsClosed() {
			return
		}
	}
}

func (s *service) updateStreaming() {
	s.localHandled = s.popStreaming(s.localHandled)
	for _, te := range s.localHandled {
		if !te.IsCancelled() {
			s.runTaskFn(te.task, s.localArgs)
		}
	}
}

func (s *service) pushHandled(t RunningTask) {
	defer sync.Lock(s.l).Unlock()
	s.handled = append(s.handled, t)
}

func (s *service) popHandled(replace []RunningTask) []RunningTask {
	replace = replace[:0]
	defer sync.Lock(s.l).Unlock()
	replace, s.handled = s.handled, replace
	return replace
}

func (s *service) pushStreaming(t RunningTask) {
	defer sync.Lock(s.l).Unlock()
	s.streaming = append(s.streaming, t)
}

func (s *service) popStreaming(replace []RunningTask) []RunningTask {
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

func (s *service) cancelEntries(entries []RunningTask) {
	for _, e := range entries {
		if e.cancel != nil {
			e.cancel.Cancel()
		}
	}
}

// ---------------------------------------------------------
// CHANNELS

type channel struct {
	c         chan RunningTask
	localArgs any
}

// ---------------------------------------------------------
// STREAMING

type taskHandlerFunc func(RunningTask)
