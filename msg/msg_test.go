package msg

import (
	"reflect"
	"testing"
)

// ---------------------------------------------------------
// TEST-PUBLISH-INT
func TestPublishInt(t *testing.T) {
	table := []struct {
		topic   string
		message int
		want    []any
	}{
		{"a", 1, []any{1}},
	}
	for i, v := range table {
		r := &Router{}
		sub := &subscription{}
		Sub(r, v.topic, sub.receiveInt)
		Pub(r, v.topic, v.message)

		if reflect.DeepEqual(v.want, sub.captured) != true {
			t.Fatalf("TestPublishInt %v has \"%v\" but wants \"%v\"", i, sub.captured, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-CHANNEL-INT
func TestChannelInt(t *testing.T) {
	table := []struct {
		topic   string
		message int
		want    []any
	}{
		{"a", 1, []any{1}},
	}
	for i, v := range table {
		r := &Router{}
		sub := &subscription{}
		Sub(r, v.topic, sub.receiveInt)
		c := NewChannel[int](r, v.topic)
		c.Pub(v.message)

		if reflect.DeepEqual(v.want, sub.captured) != true {
			t.Fatalf("TestChannelInt %v has \"%v\" but wants \"%v\"", i, sub.captured, v.want)
		}
	}
}

// ---------------------------------------------------------
// HANDLERS

type subscription struct {
	captured []any
}

func (s *subscription) receiveInt(topic string, value int) {
	s.captured = append(s.captured, value)
}

// ---------------------------------------------------------
// CONST and VAR
