package background

import (
	"slices"
	"sync/atomic"

	"github.com/hackborn/onefunc/cfg"
	"github.com/hackborn/onefunc/hub"
	"github.com/hackborn/onefunc/sync"
)

func newService(settings cfg.Settings) *service {
	c := make(chan Task, 512)
	l := sync.Mutex{}
	s := &service{c: c, l: &l,
		startUpdateFn: defaultUpdateFunc,
		endUpdateFn:   defaultUpdateFunc,
		runTaskFn:     defaultRunTaskFunc,
	}
	s.ls = &localService{s: s}
	workers := 4
	if settings.MustBool("debug/log", false) {
		timer := &debugTimer{}
		s.startUpdateFn = timer.startUpdate
		s.endUpdateFn = timer.endUpdate
		s.runTaskFn = timer.runLocal
	}
	for range workers {
		s.wg.Add(1)
		go s.loop()
	}
	return s
}

type service struct {
	closed       atomic.Bool
	c            chan Task
	wg           sync.WaitGroup
	required     []*requiredTask
	localHandled []Task
	streaming    []Task
	ls           Service
	cancels      []cancelEntry

	l       sync.Locker
	handled []Task

	// Debug support. So hacky. Wish I could cleanly separate it.
	startUpdateFn updateFunc
	endUpdateFn   updateFunc
	runTaskFn     runTaskFunc
}

// Closing shuts down before the official Close(), because
// other services might rely on the background service.
func (s *service) Closing(hub.Services) {
	if s.closed.Load() == false {
		s.cancelAll()
		s.closed.Store(true)
		close(s.c)
		s.wg.Wait()
		s.runRequired(s.required)
	}
}

func (s *service) Run(task Task) {
	// Required will be applied to everything wrapped by
	// the required; if the required is wrapped in something,
	// the parent gets lost. This is sorta weird, and points
	// to required being a flag, but that also has issues to
	// unwind, since required tracks state.
	if r, ok := getTask[*requiredTask](task); ok {
		s.required = append(s.required, r)
	}
	// If any part is streaming, the whole task is streaming.
	if _, ok := getTask[*streamingTask](task); ok {
		s.streaming = append(s.streaming, task)
	}
	if ca, ok := getTask[cancelAccess](task); ok {
		s.cancels = append(s.cancels, cancelEntry{t: task, ca: ca})
	}
	s.c <- task
}

func (s *service) RunOn(t Task, workerGroup int) {
	// Placeholder
	s.Run(t)
}

func (s *service) Update() {
	s.startUpdateFn()
	s.localHandled = s.popHandled(s.localHandled)
	if len(s.localHandled) < 1 {
		// This will be handled below, but it should
		// also be handled every update.
		s.updateStreaming()
		s.endUpdateFn()
		return
	}

	for _, t := range s.localHandled {
		// Umm what is this doing? tsk is assigned but
		// never used. Am I supposed to be using it for
		// Local()? But why? I'd want it to go through
		// the cancel check. Is this some remnant?
		tsk := t
		if ctsk, ok := tsk.(*cancelleableTask); ok {
			tsk = ctsk.task
		}
		s.runTaskFn(t)
		// Remove the task from streaming.
		s.streaming = slices.DeleteFunc(s.streaming, func(st Task) bool {
			return st == t
		})
		// Remove the cancel
		s.cancels = slices.DeleteFunc(s.cancels, func(e cancelEntry) bool {
			return e.t == t
		})
	}
	// Do this after the local handle so they aren't double-called.
	s.updateStreaming()
	s.required = slices.DeleteFunc(s.required, func(t *requiredTask) bool {
		return t.localDone.Load() == true
	})
	s.endUpdateFn()
}

func (s *service) runRequired(tasks []*requiredTask) {
	args := RemoteArgs{Local: s.ls}
	for _, t := range tasks {
		if t.remoteDone.Load() == false {
			t.Remote(args)
		}
		if t.localDone.Load() == false {
			t.Local()
		}
	}
}

func (s *service) loop() {
	defer s.wg.Done()

	for {
		task, more := <-s.c
		if !more {
			return
		}
		args := RemoteArgs{Local: s.ls}
		task.Remote(args)
		s.pushHandled(task)

		// I think the s.c == nil check is sufficient but I'm not
		// entirely positive about the safety, so I added the atomic.
		if s.c == nil || s.closed.Load() == true {
			break
		}
	}
}

func (s *service) updateStreaming() {
	for _, t := range s.streaming {
		s.runTaskFn(t)
	}
}

func (s *service) pushHandled(t Task) {
	defer sync.Lock(s.l).Unlock()
	s.handled = append(s.handled, t)
}

func (s *service) popHandled(replace []Task) []Task {
	replace = replace[:0]
	defer sync.Lock(s.l).Unlock()
	replace, s.handled = s.handled, replace
	return replace
}

type localService struct {
	s *service
}

func (s *localService) Run(t Task) {
	s.s.pushHandled(t)
}

func (s *localService) RunOn(t Task, g int) {
	s.Run(t)
}

// ---------------------------------------------------------
// CANCELLING

func (s *service) cancelAll() {
	for _, e := range s.cancels {
		e.ca.Cancel()
	}
	s.cancels = s.cancels[:0]
}

// cancelAccess is a hack interface that gets me access to
// a cancel token in the task.
// TODO: This should go away, and cancel tokens should just
// be automatically generated on Run, supplied to clients,
// and tracked internally.
type cancelAccess interface {
	Cancel()
}

type cancelEntry struct {
	t  Task
	ca cancelAccess
}
