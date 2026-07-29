package background

import (
	"flag"
	"os"
	"testing"
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
	runTestGetTask[*fakeTask](t, &fakeTask{})
	runTestGetTask[*requiredTask](t, Required(&fakeTask{}))
	runTestGetTask[*fakeTask](t, Required(&fakeTask{}))
	runTestGetTask[*fakeTask](t, Required(&cancelleableTask{task: &fakeTask{}}))
}

func runTestGetTask[T any](t *testing.T, task Task) {
	t.Helper()

	if _, ok := getTask[T](task); !ok {
		var want T
		t.Fatalf("No task found but want type %T", want)
	}
}

// ---------------------------------------------------------
// SUPPORT

type fakeTask struct {
}

func (t *fakeTask) Remote(RemoteArgs) {
}

func (t *fakeTask) Local() {
}
