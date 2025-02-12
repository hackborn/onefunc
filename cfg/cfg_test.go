package cfg

import (
	"embed"
	"os"
	"path"
	"reflect"
	"sort"
	"testing"

	"github.com/hackborn/onefunc/jacl"
	"github.com/hackborn/onefunc/math/geo"
)

func TestMain(m *testing.M) {
	setupTests()
	code := m.Run()
	os.Exit(code)
}

// ---------------------------------------------------------
// TEST-BOOL
func TestBool(t *testing.T) {
	optC := []Option{WithFS(dataFs, "testdata/c.json")}

	table := []struct {
		opts    []Option
		subset  []string
		path    string
		want    bool
		wantOk  bool
		wantErr error
	}{
		{nil, nil, "", false, false, nil},
		{optC, nil, "good", true, true, nil},
		{optC, nil, "bad", false, true, nil},
		{optC, nil, "goodstr1", true, true, nil},
		{optC, nil, "goodstr2", true, true, nil},
		{optC, nil, "list/a", true, true, nil},
	}
	for i, v := range table {
		s, haveErr := NewSettings(v.opts...)
		for _, path := range v.subset {
			s = s.Subset(path)
		}
		have, haveOk := s.Bool(v.path)
		if err := jacl.RunErr(haveErr, v.wantErr); err != nil {
			t.Fatalf("TestBool %v %v", i, err.Error())
		} else if haveOk != v.wantOk {
			t.Fatalf("TestBool %v has ok \"%v\" but wants ok \"%v\"", i, haveOk, v.wantOk)
		} else if have != v.want {
			t.Fatalf("TestBool %v has \"%v\" but wants \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-INT64
func TestInt64(t *testing.T) {
	optA := []Option{WithFS(dataFs, "testdata/a.json")}

	table := []struct {
		opts    []Option
		subset  []string
		path    string
		want    int64
		wantOk  bool
		wantErr error
	}{
		{nil, nil, "", 0, false, nil},
		{optA, nil, "age", 32, true, nil},
		{optA, nil, "run/count", 10, true, nil},
		{optA, []string{"run"}, "count", 10, true, nil},
	}
	for i, v := range table {
		s, haveErr := NewSettings(v.opts...)
		for _, path := range v.subset {
			s = s.Subset(path)
		}
		have, haveOk := s.Int64(v.path)
		if err := jacl.RunErr(haveErr, v.wantErr); err != nil {
			t.Fatalf("TestInt64 %v %v", i, err.Error())
		} else if haveOk != v.wantOk {
			t.Fatalf("TestInt64 %v has ok \"%v\" but wants ok \"%v\"", i, haveOk, v.wantOk)
		} else if have != v.want {
			t.Fatalf("TestInt64 %v has \"%v\" but wants \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-RECT
func TestRect(t *testing.T) {
	f := func(opts []Option, path string, want geo.RectI) {
		t.Helper()

		s, haveErr := NewSettings(opts...)
		havef, ok := s.RectF(path)
		have := geo.ConvertRect[float64, int](havef)
		if haveErr != nil {
			t.Fatalf("TestRect has err %v", haveErr)
		} else if !ok {
			t.Fatalf("TestRect has !ok for some reason")
		} else if !geo.RectsEqual(have, want) {
			t.Fatalf("TestRect wants %v but has %v", want, have)
		}
	}

	optC := []Option{WithFS(dataFs, "testdata/c.json")}
	f(optC, "rect1", geo.Rect(4, 8, 4, 12))
}

// ---------------------------------------------------------
// TEST-SLICES
func TestSlices(t *testing.T) {
	table := []struct {
		opts       []Option
		subset     []string
		path       string
		wantLength int
	}{
		{[]Option{WithFS(dataFs, "testdata/c.json")}, nil, "array1", 3},
		{[]Option{WithFS(dataFs, "testdata/c.json")}, nil, "array2", 4},
	}
	for i, v := range table {
		s, haveErr := NewSettings(v.opts...)
		for _, path := range v.subset {
			s = s.Subset(path)
		}
		have := s.Length(v.path)
		if err := jacl.RunErr(haveErr, nil); err != nil {
			t.Fatalf("TestSlices %v %v", i, err.Error())
		} else if have != v.wantLength {
			t.Fatalf("TestSlices %v has \"%v\" but wants \"%v\"", i, have, v.wantLength)
		}
	}
}

// ---------------------------------------------------------
// TEST-STRING
func TestString(t *testing.T) {
	table := []struct {
		opts    []Option
		subset  string
		path    string
		want    string
		wantErr error
	}{
		{nil, "", "", missingString, nil},
		{[]Option{WithFS(dataFs, "testdata/a.json")}, "", "a", "anna", nil},
		{[]Option{WithFS(dataFs, "testdata/b.json")}, "", "a", "ava", nil},
		{[]Option{WithFS(dataFs, "testdata/[{a-b}].json")}, "", "a", "ava", nil},
		{[]Option{WithEnv(EnvPattern("CFG_TESTDATA_*"))}, "", "CFG_TESTDATA_A", "ant", nil},
		{[]Option{WithEnv(EnvPrefix("CFG_TESTDATA_"))}, "", "A", "ant", nil},
		// Strings in an object in an array
		{[]Option{WithFS(dataFs, "testdata/c.json")}, "array1/0", "a", "b", nil},
		{[]Option{WithFS(dataFs, "testdata/c.json")}, "array1/1", "a", missingString, nil},
		{[]Option{WithFS(dataFs, "testdata/c.json")}, "array1/1", "c", "d", nil},
		{[]Option{WithFS(dataFs, "testdata/c.json")}, "array1/2", "a", "b", nil},
	}
	for i, v := range table {
		s, haveErr := NewSettings(v.opts...)
		if v.subset != "" {
			s = s.Subset(v.subset)
		}
		have := s.MustString(v.path, missingString)
		if err := jacl.RunErr(haveErr, v.wantErr); err != nil {
			t.Fatalf("TestString %v %v", i, err.Error())
		} else if have != v.want {
			t.Fatalf("TestString %v has \"%v\" but wants \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-STRINGS
func TestStrings(t *testing.T) {
	table := []struct {
		opts    []Option
		subset  []string
		path    string
		want    []string
		wantErr error
		ordered bool // For map testing, which will be unordered.
	}{
		{[]Option{WithFS(dataFs, "testdata/c.json")}, nil, "list", []string{"a", "b", "d"}, nil, true},
		{[]Option{WithFS(dataFs, "testdata/c.json")}, []string{"map"}, "", []string{"x", "y"}, nil, false},
		{[]Option{WithFS(dataFs, "testdata/c.json")}, nil, "smallmap", []string{"x"}, nil, true},
	}
	for i, v := range table {
		s, haveErr := NewSettings(v.opts...)
		for _, path := range v.subset {
			s = s.Subset(path)
		}
		have := s.Strings(v.path)
		if !v.ordered {
			sort.Strings(have)
			sort.Strings(v.want)
		}
		if err := jacl.RunErr(haveErr, v.wantErr); err != nil {
			t.Fatalf("TestStrings %v %v", i, err.Error())
		} else if !reflect.DeepEqual(have, v.want) {
			t.Fatalf("TestStrings %v has \"%v\" but wants \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-WITH-SETTINGS
func TestWithSettings(t *testing.T) {
	f := func(src, opt string, want string) {
		t.Helper()

		withOpt := []Option{WithFS(dataFs, path.Join("testdata", opt))}
		optS, haveErr := NewSettings(withOpt...)
		if haveErr != nil {
			panic(haveErr)
		}

		opts := []Option{WithFS(dataFs, path.Join("testdata", src))}
		opts = append(opts, WithSettings(optS))
		s, haveErr := NewSettings(opts...)
		s.Print()
		if haveErr != nil {
			panic(haveErr)
		}

		// TODO: Ahh, map compares. Need a good way of comparing
		// the result.
		//		if !reflect.DeepEqual(s, nil) {
		//			t.Fatalf("TestWith has %v\nbut wants\n%v", s, nil)
		//		}
	}

	f("a.json", "b.json", "c.json")
}

// ---------------------------------------------------------
// TEST-HEX-TO-UINT8
func TestHexToUint8(t *testing.T) {
	f := func(s string, idx int, fallback uint8, want uint8) {
		t.Helper()

		have := hexToUint8(s, idx, fallback)
		if want != have {
			t.Fatalf("TestHexToUint8 wants %v but has %v", want, have)
		}
	}
	f("s", 0, 10, 10)
	f("FF", 0, 10, 255)
	f("ff", 0, 10, 255)
}

// ---------------------------------------------------------
// LIFECYCLE

func setupTests() {
	os.Setenv("CFG_TESTDATA_A", "ant")
	os.Setenv("CFG_TESTDATA_B", "bear")
	os.Setenv("CFG_TESTDATA_C", "cat")
}

// ---------------------------------------------------------
// CONST and VAR

const (
	missingString = "~missing~"
)

//go:embed testdata/*
var dataFs embed.FS
