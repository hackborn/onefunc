package background

import (
	"fmt"
	"time"
)

// ---------------------------------------------------------
// FUNC

type updateFunc func()
type runTaskFunc func(Task)

func defaultUpdateFunc() {
}

func defaultRunTaskFunc(t Task) {
	t.Local()
}

// ---------------------------------------------------------
// DEBUG-TIMER

// debugTimer is an ugly way to monitor
// how long the update takes.
type debugTimer struct {
	entries []debugTimerEntry
	start   time.Time
}

type debugTimerEntry struct {
	task Task
	dur  time.Duration
}

func (t *debugTimer) startUpdate() {
	t.entries = t.entries[:0]
	t.start = time.Now()
}

func (t *debugTimer) endUpdate() {
	dur := time.Since(t.start)
	if dur > time.Millisecond*20 {
		for _, e := range t.entries {
			fmt.Printf("Background pacing %v for %T\n", e.dur, e.task)
		}
	}
}

func (t *debugTimer) runLocal(task Task) {
	prev := time.Now()
	task.Local()
	diff := time.Since(prev)
	t.entries = append(t.entries, debugTimerEntry{task: task, dur: diff})
}
