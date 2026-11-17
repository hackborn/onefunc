package background

import (
	"fmt"

	"github.com/hackborn/onefunc/sync"

	"github.com/hackborn/onefunc/hub"
)

// Run will run the task in the background.
func Run(all hub.Services, task Task) (TaskInfo, error) {
	if s, err := hub.Get[*service](all, ServiceName); err == nil {
		return s.Run(task), nil
	} else {
		return TaskInfo{}, fmt.Errorf("No "+ServiceName+" service: %w", err)
	}
}

// RunRequired will run the task with the guarantee that the task
// is run, even if the app is shut down.
func RunRequired(all hub.Services, task Task) (TaskInfo, error) {
	return Run(all, Required(task))
}

// Required makes a task required, so it will run even if
// the app is shutting down.
func Required(t Task) Task {
	return &requiredTask{task: t}
}

type Service interface {
	Run(Task) TaskInfo
	// Run on a specific group, as defined by the settings.
	// If the group doesn't exist then run on the default 0 group.
	RunOn(t Task, workerGroup int) TaskInfo
}

// Task defines a single background task.
type Task interface {
	// Remote is run on the background thread.
	// Run tasks on local if you want to report
	// intermediary data.
	Remote(RemoteArgs)

	// Local is run on the main thread.
	Local(LocalArgs)
}

type RemoteArgs struct {
	// cancel is an optional cancel token.
	cancel *sync.Cancel

	// streamingFn handles sending tasks to the main thread.
	streamingFn taskHandlerFunc

	// data stores data supplied by clients. Access it by type with Get()
	data []any
}

func (a RemoteArgs) IsCancelled() bool {
	if a.cancel == nil {
		return false
	}
	return a.cancel.IsCancelled()
}

// Stream is used to send a task to the local thread
// n number of times during remote execution.
func (a RemoteArgs) Stream(t Task) {
	if a.streamingFn != nil {
		a.streamingFn(RunningTask{task: t, cancel: a.cancel})
	}
}

// Get answers the data type, or nil if none.
func (a RemoteArgs) Get[T any]() *T {
	for _, _d := range a.data {
		if d, ok := _d.(*T); ok {
			return d
		}
	}
	return nil
}

type LocalArgs struct {
	// data stores data supplied by clients. Access it by type with Get()
	data []any
}

// Get answers the data type, or nil if none.
func (a LocalArgs) Get[T any]() *T {
	for _, _d := range a.data {
		if d, ok := _d.(*T); ok {
			return d
		}
	}
	return nil
}

type TaskInfo struct {
	Cancel *sync.Cancel
}

// RunningTask is used to hook into the new pluggable
// run-local system.
type RunningTask struct {
	task   Task
	cancel *sync.Cancel
}

func NewRunningTask(task Task, cancel *sync.Cancel) RunningTask {
	return RunningTask{task: task, cancel: cancel}
}

// WithTask answers the myself with my task replaced, but my cancel intact.
// This feels overengineered because I got backed into a corner where
// some systems expect the original cancel struct and some transform it
// into a cancelchecker, and I'm not sure how to reconcile.
func (t RunningTask) WithTask(task Task) RunningTask {
	return RunningTask{task: task, cancel: t.cancel}
}

func (t RunningTask) IsCancelled() bool {
	if t.cancel != nil {
		return t.cancel.IsCancelled()
	}
	return false
}

// ---------------------------------------------------------
// FUNC

// RunLocalFunc can be installed on the service to provide
// a custom local handler. It runs in a worker thread.
type RunLocalFunc func(Task, LocalArgs)

// SendLocalFunc is used to expose the RunLocalFunc to clients.
// By default it will send the task to the main thread to run,
// but if a RunLocalFunc is installed it will use that in the
// worker thread. This satisfies some UI lib requirements.
type SendLocalFunc func(RunningTask)
