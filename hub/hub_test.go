package hub

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/hackborn/onefunc/cfg"
	oferrors "github.com/hackborn/onefunc/errors"
)

// go test . -update
var (
	update = flag.Bool("update", false, "update the golden files of this test")
)

func TestMain(m *testing.M) {
	flag.Parse()
	initFactories()
	os.Exit(m.Run())
}

// ---------------------------------------------------------
// TEST-BUILD-FINISHED

// TestBuildFinished tests the BuildFinisher lifecycle.
func TestBuildFinished(t *testing.T) {
	settings := strings.ReplaceAll(settings1, "$value$", doAddseValue)
	services, err := buildServices(settings)
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if _, err := Get[*testAddedService](services, testAddedServiceName); err != nil {
		t.Fatalf("Wanted added service but got err %v", err)
	}
	if ts, err := Get[*testRemovedService](services, testRemovedServiceName); err == nil {
		t.Fatalf("Removed service should have been removed but found %v %T", ts, ts)
	}
	// The removed service is gone but it made a change
	ts, err := Get[*testService](services, testServiceName)
	if err != nil {
		t.Fatalf("Missing service: %v %v", testServiceName, err)
	}
	if ts.value != removedServiceValue{
		t.Fatalf("Built but value is %v and want %v", ts.value, removedServiceValue)
	}
}

// ---------------------------------------------------------
// TEST-BUILD-AND-GET

// TestBuildAndGet tests building the services and getting a service.
func TestBuildAndGet(t *testing.T) {
	want := "xx"
	settings := strings.ReplaceAll(settings1, "$value$", want)
	services, err := buildServices(settings)
	if err != nil {
		t.Fatalf("err %v", err)
	}
	ts, err := Get[*testService](services, testServiceName)
	if err != nil {
		t.Fatalf("Missing service: %v %v", testServiceName, err)
	}
	if ts.value != want {
		t.Fatalf("Built but value is %v and want %v", ts.value, want)
	}
}

// ---------------------------------------------------------
// SERVICES

func buildServices(settingsStr string) (Services, error) {
	settings, err := cfg.NewSettings(cfg.WithString(settingsStr))
	oferrors.Panic(err)
	args := BuildArgs{}
	//	args.Steps = append(args.Steps, filetree.NewBuildServiceFs(settings.Subset("services/filetree"), myapp.FiletreeFrelName, fsys.FS()))
	return Build(args, settings.Subset("services"))
}

const testServiceName = "test"
const testAddedServiceName = "added"
const testRemovedServiceName = "removed"

type testService struct {
	value string
}

type testAddedService struct {
	s *testService
}

type testRemovedService struct {
	s *testService
}

func (s *testRemovedService) BuildFinished() {
	s.s.value = removedServiceValue
}

func initFactories() {
	fac := Factory{}
	fac.NewServiceFn = func(args BuildArgs, settings cfg.Settings) (any, error) {
		return &testService{value: settings.MustString("value", "")}, nil
	}
	fac.OnNewServiceFn = func(args OnNewServiceArgs) {
		if s, ok := args.Service.(*testService); ok {
			if s.value == doAddseValue {
			args.SetService(testAddedServiceName, &testAddedService{s: s})
			args.SetService(testRemovedServiceName, &testRemovedService{s: s})
			}
		}
	}
	RegisterFactory(testServiceName, fac)
}

// ---------------------------------------------------------
// CONST and VAR

const doAddseValue = "do_adds"
const removedServiceValue = "gone"

var settings1 = `{
  "services": {
    "test": {
      "value": "$value$"
    }
  }
}`
