package background

import (
	"sync/atomic"

	"github.com/hackborn/onefunc/sync"
)

// ---------------------------------------------------------
// CANCELLEABLE-TASK

// cancelleableTask includes a token that can be used for cancelling.
type cancelleableTask struct {
	task   Task
	cancel *sync.Cancel
}

func (t *cancelleableTask) Cancel() {
	if t.cancel != nil {
		t.cancel.Cancel()
	}
}

func (t *cancelleableTask) Unwrap() Task {
	return t.task
}

func (t *cancelleableTask) Remote(args RemoteArgs) {
	if !t.cancel.IsCancelled() {
		args.cancel = t.cancel
		t.task.Remote(args)
	}
}

func (t *cancelleableTask) Local() {
	if !t.cancel.IsCancelled() {
		t.task.Local()
	}
}

// ---------------------------------------------------------
// REQUIRED-TASK

// requiredTask is always run, even if the app is closing.
type requiredTask struct {
	task                  Task
	remoteDone, localDone atomic.Bool
}

func (t *requiredTask) Unwrap() Task {
	return t.task
}

func (t *requiredTask) Remote(args RemoteArgs) {
	t.task.Remote(args)
	t.remoteDone.Store(true)
}

func (t *requiredTask) Local() {
	t.task.Local()
	t.localDone.Store(true)
}

// ---------------------------------------------------------
// STREAMING-TASK

// streamingTask runs local immediately until remote is done.
type streamingTask struct {
	task Task
}

func (t *streamingTask) Unwrap() Task {
	return t.task
}

func (t *streamingTask) Remote(args RemoteArgs) {
	t.task.Remote(args)
}

func (t *streamingTask) Local() {
	t.task.Local()
}
