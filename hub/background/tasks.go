package background

import (
	"sync/atomic"
)

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
