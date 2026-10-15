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
	Local()
}

type RemoteArgs struct {
	// cancel is an optional cancel token.
	cancel *sync.Cancel

	// streamingFn handles sending tasks to the main thread.
	streamingFn taskHandlerFunc
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
		a.streamingFn(taskEntry{t: t, cancel: a.cancel})
	}
}

type TaskInfo struct {
	Cancel *sync.Cancel
}
