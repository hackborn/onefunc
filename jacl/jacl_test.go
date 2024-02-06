package jacl

import (
	"fmt"
	"testing"
)

// ---------------------------------------------------------
// TEST-RUN
func TestRun(t *testing.T) {
	table := []struct {
		dst     any
		exprs   []string
		wantErr error
	}{
		{Field{Name: "a"}, []string{`Name=a`}, nil},
		{[]Field{{Name: "a"}}, []string{`0/Name=a`}, nil},
		{Field{}, []string{`{type}="Field"`}, nil},
		{&Field{}, []string{`{type}="*Field"`}, nil},
		{map1, []string{`a/Name=blip`}, nil},
		// Errors
		{[]Field{{Name: "a"}}, []string{`1/Name=a`}, fmt.Errorf("out of range")},
	}
	for i, v := range table {
		haveErr := Run(v.dst, v.exprs...)

		if v.wantErr == nil && haveErr != nil {
			t.Fatalf("TestRun %v expected no error but has %v", i, haveErr)
		} else if v.wantErr != nil && haveErr == nil {
			t.Fatalf("TestRun %v has no error but exptected %v", i, v.wantErr)
		}
	}
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
			t.Fatalf("TestRunOpts %v has no error but exptected %v", i, v.wantErr)
		}
	}
}

// ---------------------------------------------------------
// SUPPORT

type Field struct {
	Name string
	IV   int
	SV   string
}

var (
	map1 = map[string]Field{
		"a": Field{Name: "blip", SV: "bloop"},
	}
)
