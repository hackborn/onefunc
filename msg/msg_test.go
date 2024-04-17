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
// TEST-SEQUENCE
func TestSequence(t *testing.T) {
	table := []struct {
		steps []seqStep
		want  []any
	}{
		{steps(sub("", "a"), channel("ca", "a"), pubChannel("ca", 10)), []any{capture{"a", 10}}},
		// Channels are updated after modifications to the topic tree.
		{steps(sub("s", "a"), channel("ca", "a"), unsub("s"), pubChannel("ca", 10)), nil},
	}
	for i, v := range table {
		state := newState()
		for _, step := range v.steps {
			step.Step(state)
		}
		if reflect.DeepEqual(v.want, state.captured) != true {
			t.Fatalf("TestSequence %v has \"%v\" but wants \"%v\"", i, state.captured, v.want)
		}
	}
}

// A whole bunch of cruft to handle creating test sequences

func newState() *seqState {
	r := &Router{}
	data := make(map[string]any)
	return &seqState{r: r, data: data}
}

func steps(steps ...seqStep) []seqStep {
	return steps
}

func sub(name, topic string) seqStep {
	return &_seqSub{name: name, topic: topic}
}

func unsub(name string) seqStep {
	return &_seqUnsub{name: name}
}

func channel(name, topic string) seqStep {
	return &_seqChannel{name: name, topic: topic}
}

func pubChannel(name string, value int) seqStep {
	return &_seqPubChannel{name: name, value: value}
}

type seqState struct {
	r        *Router
	data     map[string]any
	captured []any
}

func (s *seqState) addData(name string, value any) {
	if name == "" {
		panic("no name")
	}
	if _, ok := s.data[name]; ok {
		panic("name " + name + " exists")
	}
	s.data[name] = value
}

func getData[T any](s *seqState, name string) (T, bool) {
	if _v, ok := s.data[name]; ok {
		v, ok := _v.(T)
		return v, ok
	}
	var t T
	return t, false
}

func (s *seqState) handleInt(topic string, value int) {
	s.captured = append(s.captured, capture{topic: topic, value: value})
}

type capture struct {
	topic string
	value int
}

type seqStep interface {
	Step(*seqState)
}

type _seqSub struct {
	state *seqState
	name  string
	topic string
}

func (s *_seqSub) Step(state *seqState) {
	s.state = state
	sub := Sub(state.r, s.topic, s.handleInt)
	if s.name != "" {
		state.addData(s.name, sub)
	}
}

func (s *_seqSub) handleInt(topic string, value int) {
	s.state.handleInt(topic, value)
}

type _seqUnsub struct {
	name string
}

func (s *_seqUnsub) Step(state *seqState) {
	if sub, ok := getData[Subscription](state, s.name); ok {
		sub.Unsub()
	}
}

type _seqChannel struct {
	name  string
	topic string
}

func (s *_seqChannel) Step(state *seqState) {
	c := NewChannel[int](state.r, s.topic)
	state.addData(s.name, c)
}

type _seqPubChannel struct {
	name  string
	value int
}

func (s *_seqPubChannel) Step(state *seqState) {
	if c, ok := getData[Channel[int]](state, s.name); ok {
		c.Pub(s.value)
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
