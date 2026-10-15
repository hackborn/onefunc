package background

import (
	"flag"
	"fmt"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hackborn/onefunc/cfg"
	"github.com/hackborn/onefunc/sync"
)

// go test . -update
var (
	update = flag.Bool("update", false, "update the golden files of this test")
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

// ---------------------------------------------------------
// TEST-GET-TASK

func TestGetTask(t *testing.T) {
	runTestGetTask[*basicRunTask](t, &basicRunTask{})
	runTestGetTask[*requiredTask](t, Required(&basicRunTask{}))
	runTestGetTask[*basicRunTask](t, Required(&basicRunTask{}))
}

func runTestGetTask[T any](t *testing.T, task Task) {
	t.Helper()

	if _, ok := getTask[T](task); !ok {
		var want T
		t.Fatalf("No task found but want type %T", want)
	}
}

// ---------------------------------------------------------
// TEST-RUN

// TestRun starts a service, sends it a task, shuts down
// the service, and checks if the task is successful.
func TestRun(t *testing.T) {
	f := func(task runTask, startFn startRunTask) {
		t.Helper()

		s := newService(newSettings())
		if startFn != nil {
			startFn(s, task)
		} else {
			s.Run(task)
		}
		doUntil(time.Millisecond*100, func() error {
			s.Update()
			return task.Error()
		})
		s.Closing(nil)

		if haveErr := task.Error(); haveErr != nil {
			t.Fatalf("err: %v", haveErr)
		} else {
			//			t.FailNow()
		}
	}
	f(&basicRunTask{}, nil)
	// If you look at the cancel task you'll see that it only
	// verifies that the local won't get called on cancel,
	// not remote. Also it seems like the test can't really
	// fail, but still makes sure things are working.
	f(&cancelRunTask{}, startAndCancelTask)
	f(&streamingRunTask{neededCount: 5}, nil)
}

// ---------------------------------------------------------
// RUN START

type startRunTask func(s *service, t Task)

func startAndCancelTask(s *service, t Task) {
	info := s.Run(t)
	if tc, ok := t.(cancellableSetter); ok {
		tc.setCancel(info.Cancel)
	}
}

// ---------------------------------------------------------
// RUN TASKS

type runTask interface {
	Task
	// Check if this is in an error state. Note that this
	// needs to be thread-safe, as it is called locally
	// until the task completes.
	Error() error
}

type cancellableSetter interface {
	setCancel(*sync.Cancel)
}

type basicRunTask struct {
	remoteCalled atomic.Bool
	localCalled  atomic.Bool
}

func (t *basicRunTask) Remote(RemoteArgs) {
	t.remoteCalled.Store(true)
}

func (t *basicRunTask) Local() {
	t.localCalled.Store(true)
}

func (t *basicRunTask) Error() error {
	if t.remoteCalled.Load() == false {
		return fmt.Errorf("Remote never called")
	}
	if t.localCalled.Load() == false {
		return fmt.Errorf("Local never called")
	}
	return nil
}

type cancelRunTask struct {
	localCalled atomic.Bool
	cancel      *sync.Cancel
}

func (t *cancelRunTask) Remote(RemoteArgs) {
	if t.cancel != nil {
		t.cancel.Cancel()
	}
}

func (t *cancelRunTask) Local() {
	if t.cancel != nil && !t.cancel.IsCancelled() {
		t.localCalled.Store(true)
	}
}

func (t *cancelRunTask) Error() error {
	if t.cancel == nil {
		return fmt.Errorf("Missing cancel")
	}
	if !t.cancel.IsCancelled() {
		return fmt.Errorf("Not canceled")
	}
	if t.localCalled.Load() == true {
		return fmt.Errorf("Local was called (never cancelled?)")
	}
	return nil
}

func (t *cancelRunTask) setCancel(cancel *sync.Cancel) {
	t.cancel = cancel
}

type streamingRunTask struct {
	count, neededCount int
}

func (t *streamingRunTask) Remote(args RemoteArgs) {
	for range t.neededCount {
		args.Stream(&streamingCounter{t: t})
	}
}

func (t *streamingRunTask) Local() {
}

func (t *streamingRunTask) Error() error {
	if t.count != t.neededCount {
		return fmt.Errorf("Count has %v but wants %v", t.count, t.neededCount)
	}
	return nil
}

type streamingCounter struct {
	t *streamingRunTask
}

func (t *streamingCounter) Remote(args RemoteArgs) {
}

func (t *streamingCounter) Local() {
	if t.t != nil {
		t.t.count++
	}
}

// ---------------------------------------------------------
// WAITING

// doUntil will repeatedly run fn until it returns nil,
// unless dur is reach, at which point it will return
// the last error.
func doUntil(dur time.Duration, fn func() error) error {
	t := time.Now().Add(dur)
	for time.Now().Before(t) {
		if fn() == nil {
			return nil
		}
		time.Sleep(time.Millisecond * 1)
	}
	return fn()
}

// ---------------------------------------------------------
// BOILERPLATE

func newSettings() cfg.Settings {
	s, _ := cfg.NewSettings()
	return s
}
