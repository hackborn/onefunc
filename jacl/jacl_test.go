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
// SUPPORT

type Field struct {
	Name string
	IV   int
	SV   string
}
