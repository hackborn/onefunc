package slices

import (
	"reflect"
	"testing"
)

// ---------------------------------------------------------
// TEST-ARRAY-FROM
func TestArrayFrom(t *testing.T) {
	table := []struct {
		src  []any
		fn   func(any) any
		want []any
	}{
		{nil, func(any) any { return "a" }, []any{}},
		{[]any{"a", "b"}, func(any) any { return "a" }, []any{"a", "a"}},
		{[]any{"a", "b"}, func(a any) any { return a }, []any{"a", "b"}},
	}
	for i, v := range table {
		have := ArrayFrom(v.src, v.fn)

		if !reflect.DeepEqual(have, v.want) {
			t.Fatalf("TestArrayFrom %v has %v but wants %v", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-POP
func TestPop(t *testing.T) {
	f := func(s []int, wantS []int, wantItem int) {
		t.Helper()

		haveS, haveItem := Pop(s)
		if haveItem != wantItem {
			t.Fatalf("TestPop has %v but wants %v", haveItem, wantItem)
		} else if reflect.DeepEqual(haveS, wantS) == false {
			t.Fatalf("TestPop has slice %v but wants %v", haveS, wantS)
		}
	}

	f([]int{1, 2, 3}, []int{1, 2}, 3)
}
