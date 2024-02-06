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
