package dep

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hackborn/onefunc/jacl"
)

// go test -bench=.

// ---------------------------------------------------------
// TEST-DEPSORT
func TestDepSort(t *testing.T) {
	f := func(items []depItem, wantNames string, wantErr error) {
		t.Helper()

		have, haveErr := Sort(items)
		if err := jacl.RunErr(haveErr, wantErr); err != nil {
			t.Fatalf("Has err %v but wants %v", haveErr, wantErr)
		} else if wantErr == nil {
			haveNames := ""
			for _, item := range have {
				haveNames += item.name
			}
			if haveNames != wantNames {
				t.Fatalf("Has %v but wants %v", haveNames, wantNames)
			}
		}
	}
	// TODO: This won't work, the alo is based on maps
	// so there will be some randomness. Need to work that out.
	f(depList("A", "+B", "B", "C", "+B", "D", "+A"), "BACD", nil)
	f(depList("X", "+Y", "Y", "+Z", "Z", "+X"), "", fmt.Errorf("cycle"))
}

// ---------------------------------------------------------
// TYPES

// depList makes a depItem slice by taking any strings
// that start with "+", stripping it off, and making them
// dependents of the current item.
func depList(ss ...string) []depItem {
	ans := []depItem{}
	var cur *depItem
	for _, s := range ss {
		if strings.HasPrefix(s, "+") {
			if cur != nil {
				cur.deps = append(cur.deps, strings.TrimPrefix(s, "+"))
			}
		} else {
			if cur != nil {
				ans = append(ans, *cur)
			}
			cur = &depItem{name: s}
		}
	}
	if cur != nil {
		ans = append(ans, *cur)
	}
	return ans
}

type depItem struct {
	name string
	deps []string
}

func (d depItem) Name() string {
	return d.name
}

func (d depItem) Dependencies() []string {
	return d.deps
}
