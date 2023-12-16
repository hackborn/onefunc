package strings

import (
	"io"
	"testing"
)

// ---------------------------------------------------------
// TEST-STRING-WRITER-POOL
func TestStringWriterPool(t *testing.T) {
	type State struct {
		pool  *pool
		Stack []io.StringWriter
	}
	type ActionFunc func(*State) error
	acquire := func(s *State) error {
		w := s.pool.Acquire()
		s.Stack = append(s.Stack, w)
		return nil
	}
	release := func(s *State) error {
		size := len(s.Stack)
		if size > 0 {
			s.pool.Release(s.Stack[size-1])
			s.Stack = s.Stack[:size-1]
		}
		return nil
	}

	table := []struct {
		actions []ActionFunc
		want    string
		wantLen int
		wantErr error
	}{
		{[]ActionFunc{acquire}, "", 0, nil},
		{[]ActionFunc{acquire, release}, "", 1, nil},
		{[]ActionFunc{acquire, release, acquire}, "", 0, nil},
		{[]ActionFunc{acquire, release, release}, "", 1, nil},
	}
	for i, v := range table {
		state := &State{}
		state.pool = newLockingPool()

		var haveErr error
		for _, action := range v.actions {
			haveErr = firstErr(haveErr, action(state))
		}
		have := ""
		haveLen := len(state.pool.builders)

		if v.wantErr == nil && haveErr != nil {
			t.Fatalf("TestStringWriterPool %v expected no error but has %v", i, haveErr)
		} else if v.wantErr != nil && haveErr == nil {
			t.Fatalf("TestStringWriterPool %v has no error but exptected %v", i, v.wantErr)
		} else if have != v.want {
			t.Fatalf("TestStringWriterPool %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		} else if haveLen != v.wantLen {
			t.Fatalf("TestStringWriterPool %v has len \"%v\" but wanted \"%v\"", i, haveLen, v.wantLen)
		}
	}
}
