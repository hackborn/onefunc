package jacl

import (
	"fmt"
	"testing"
)

// ---------------------------------------------------------
// TEST-RUN
func TestRun(t *testing.T) {
	f := func(dst any, wantErr error, want ...string) {
		t.Helper()

		haveErr := Run(dst, want...)
		if wantErr == nil && haveErr != nil {
			t.Fatalf("TestRun expected no error but has %v", haveErr)
		} else if wantErr != nil && haveErr == nil {
			t.Fatalf("TestRun has no error but expected %v", wantErr)
		}
	}
	f(Field{Name: "a"}, nil, `Name=a`)
	f([]Field{{Name: "a"}}, nil, `0/Name=a`)
	f(Field{}, nil, `{type}="Field"`)
	f(&Field{}, nil, `{type}="*Field"`)
	f(map1, nil, `a/Name=blip`)
	f(map1, nil, `b/Name=bop`, `b/IV=10`)
	f(map1, nil, `ct/Name=":)"`, `ct/BV=true`)
	f(map2, nil, `a/b/Name=dash`)
	f(map2, nil, `a/""/Name=found`)
	// Errors
	f([]Field{{Name: "a"}}, fmt.Errorf("out of range"), `1/Name=a`)
}

// ---------------------------------------------------------
// TEST-RUN-OPTS
func TestRunOpts(t *testing.T) {
	table := []struct {
		opts    Opts
		dst     any
		exprs   []string
		wantErr error
	}{
		// QUOTES
		// Single quotes are passed down in the string.
		{Opts{}, Field{Name: `'a'`}, []string{`Name='a'`}, nil},
		// Two single quotes are reduced to double quotes with the default opts.
		{Opts{}, Field{Name: `a"b`}, []string{`Name="a''b"`}, nil},
		// Single-to-double quotes around the entire phrase are not removed,
		// and the comparison term must have the double quotes.
		{Opts{}, Field{Name: `"a"`}, []string{`Name="''a''"`}, nil},
		// Can disable the single-to-double quotes
		{Opts{RawValues: true}, Field{Name: `a''`}, []string{`Name="a''"`}, nil},
	}
	for i, v := range table {
		haveErr := RunOpts(v.opts, v.dst, v.exprs...)

		if v.wantErr == nil && haveErr != nil {
			t.Fatalf("TestRunOpts %v expected no error but has %v", i, haveErr)
		} else if v.wantErr != nil && haveErr == nil {
			t.Fatalf("TestRunOpts %v has no error but expected %v", i, v.wantErr)
		}
	}
}

// ---------------------------------------------------------
// SUPPORT

type Field struct {
	Name string
	IV   int
	SV   string
	BV   bool
}

var (
	map1 = map[string]Field{
		"a":  {Name: "blip", SV: "bloop"},
		"b":  {Name: "bop", IV: 10},
		"ct": {Name: ":)", BV: true},
		"cf": {Name: ":(", BV: false},
	}

	map2 = map[string]map[string]Field{
		"a": {
			"":  {Name: "found"},
			"b": {Name: "dash"},
		},
	}
)
