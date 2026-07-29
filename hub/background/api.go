package background

import (
	"fmt"

	"github.com/hackborn/onefunc/sync"

	"github.com/hackborn/onefunc/hub"
)

// Run will run the task in the background.
func Run(all hub.Services, task Task) error {
	if s, err := hub.Get[*service](all, ServiceName); err == nil {
		s.Run(task)
		return nil
	} else {
		return fmt.Errorf("No "+ServiceName+" service: %w", err)
	}
}

// RunRequired will run the task with the guarantee that the task
// is run, even if the app is shut down.
func RunRequired(all hub.Services, task Task) error {
	return Run(all, Required(task))
}

// Cancelleable wraps a task and provides a cancel token.
// TODO: Cancelleable shouldn't be a task; every Run()
// should generate a cancel token that is provided to clients
// and the task chain.
func Cancelleable(t Task) (Task, *sync.Cancel) {
	cancel := &sync.Cancel{}
	return &cancelleableTask{task: t, cancel: cancel}, cancel
}

// Required makes a task required, so it will run even if
// the app is shutting down.
func Required(t Task) Task {
	return &requiredTask{task: t}
}

// Streaming causes a task Local() to be called
// immediately and until Remote() finishes.
func Streaming(t Task) Task {
	return &streamingTask{task: t}
}

type Service interface {
	Run(Task)
	// Run on a specific group, as defined by the settings.
	// If the group doesn't exist then run on the default 0 group.
	RunOn(t Task, workerGroup int)
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
	// Local is a Service for sending any Tasks to the
	// main (i.e. local) thread. Long-running operations
	// might use it to report progress.
	Local Service

	// cancel is an optional cancel token.
	cancel *sync.Cancel
}

func (a RemoteArgs) IsCancelled() bool {
	if a.cancel == nil {
		return false
	}
	return a.cancel.IsCancelled()
}
