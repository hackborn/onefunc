package strings

import (
	_ "fmt"
	"io"
	"testing"

	"github.com/hackborn/onefunc/errors"
	"github.com/hackborn/onefunc/maps"
)

// ---------------------------------------------------------
// TEST-MQTT-MATCH
func TestMqttMatch(t *testing.T) {
	table := []struct {
		pattern string
		topic   string
		want    bool
	}{
		{"a", "a", true},
		{"a/b", "a/b", true},
		{"a/#", "a/b", true},
		{"a/#", "a/b/c", true},
		{"a", "b", false},
		{"a/b", "a/c", false},
		{"a/+/c", "a", false},
		{"a/+/c", "a/c", false},
		{"a/+/c", "a/b/c", true},
		{"a/+/c", "a/c/c", true},
		{"a/+/c", "a/c/c/c", false},
	}
	for i, v := range table {
		have := MqttMatch(v.pattern, v.topic)

		if v.want != have {
			t.Fatalf("TestMqttMatch %v has %v but wanted %v", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-STRING-WRITER-POOL
func TestStringWriterPool(t *testing.T) {
	type State struct {
		pool    Pool
		rawPool *maps.Pool[uint64, io.StringWriter]
		Stack   []io.StringWriter
	}
	type ActionFunc func(*State) error
	get := func(s *State) error {
		w := s.pool.Get(nil)
		s.Stack = append(s.Stack, w)
		return nil
	}
	put := func(s *State) error {
		size := len(s.Stack)
		if size > 0 {
			s.pool.Put(s.Stack[size-1])
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
		{[]ActionFunc{get}, "", 0, nil},
		{[]ActionFunc{get, put}, "", 1, nil},
		{[]ActionFunc{get, put, get}, "", 0, nil},
		{[]ActionFunc{get, put, put}, "", 1, nil},
	}
	for i, v := range table {
		state := &State{}
		state.rawPool = newLockingPool()
		state.pool = state.rawPool

		var haveErr error
		for _, action := range v.actions {
			haveErr = errors.First(haveErr, action(state))
		}
		have := ""
		haveLen := state.rawPool.CacheSize()

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
