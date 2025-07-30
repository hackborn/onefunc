package msg

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
	"text/scanner"

	oferrors "github.com/hackborn/onefunc/errors"
)

// ---------------------------------------------------------
// TEST-PUBLISH-INT

func TestPublishInt(t *testing.T) {
	f := func(topic string, message int, want []any) {
		t.Helper()

		r := NewRouter()
		sub := &captureSubscription{}
		Sub(r, topic, sub.receiveInt)
		Pub(r, topic, message)

		if reflect.DeepEqual(want, sub.captured) != true {
			t.Fatalf("has \"%v\" but wants \"%v\"", sub.captured, want)
		}
	}
	f("a", 1, []any{1})
}

// ---------------------------------------------------------
// TEST-CHANNEL-INT

func TestChannelInt(t *testing.T) {
	f := func(topic string, message int, want []any) {
		t.Helper()

		sub := &captureSubscription{}
		r := NewRouter()
		Sub(r, topic, sub.receiveInt)
		c := NewChannel[int](r, topic)
		c.Pub(message)

		if reflect.DeepEqual(want, sub.captured) != true {
			t.Fatalf("has \"%v\" but wants \"%v\"", sub.captured, want)
		}
	}
	f("a", 1, []any{1})
}

// ---------------------------------------------------------
// TEST-CHANNEL-INTERFACE

func TestChannelInterface(t *testing.T) {
	f := func(topic string, message Namer, want []Namer) {
		t.Helper()

		sub := &captureSubscription{}
		r := NewRouter()
		Sub(r, topic, sub.receiveNamer)
		c := NewChannel[Namer](r, topic)
		c.Pub(message)

		if reflect.DeepEqual(want, sub.namerCaptured) != true {
			t.Fatalf("has \"%v\" but wants \"%v\"", sub.namerCaptured, want)
		}
	}
	// Used to cause crash.
	f("a", nil, []Namer{nil})
	// TODO: Need a different way to test for this to work
	//	f("a", &_namer{n: "hi"}, []Namer{&_namer{n: "hi"}})
}

type Namer interface {
	Name() string
}

type _namer struct {
	n string
}

func (n *_namer) Name() string {
	return n.n
}

// ---------------------------------------------------------
// TEST-SEQUENCE
func TestSequence(t *testing.T) {
	f := func(test string) {
		t.Helper()

		tests := loadSeqTests(test)
		for _, test := range tests {
			state := newState()
			for _, step := range test.Steps {
				step.Step(state)
			}
			if reflect.DeepEqual(state.want, state.captured) != true {
				t.Fatalf("test %v has \"%v\" but wants \"%v\"", state.name, state.captured, state.want)
			}
		}
	}
	f("sequences_1.txt")
	f("sequences_2.txt")
}

// ---------------------------------------------------------
// SEQUENCE SUPPORT
// A whole bunch of cruft to handle creating test sequences

type seqTest struct {
	Steps []seqStep
}

func newState() *seqState {
	r := NewRouter()
	data := make(map[string]any)
	return &seqState{r: r, data: data}
}

type seqState struct {
	r        *Router
	name     string
	want     []any
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

func (s *seqState) handleString(topic string, value string) {
	s.captured = append(s.captured, capture{Topic: topic, Value: value})
}

type capture struct {
	Topic string
	Value string
}

type seqStep interface {
	Step(*seqState)
	SetParam(name string, value any)
}

func getStringParam(nameA, nameB string, value any) (string, bool) {
	if nameA != nameB {
		return "", false
	}
	v, ok := value.(string)
	return v, ok
}

type _seqName struct {
	Id string
}

func (s *_seqName) Step(state *seqState) {
	state.name = s.Id
	// fmt.Println("sequence", s.Id)
}

func (s *_seqName) SetParam(name string, value any) {
	if v, ok := getStringParam("id", name, value); ok {
		s.Id = v
	}
}

type _seqSub struct {
	state *seqState
	Id    string
	Topic string
}

func (s *_seqSub) Step(state *seqState) {
	s.state = state
	sub := Sub(state.r, s.Topic, s.handleString)
	if s.Id != "" {
		state.addData(s.Id, sub)
	}
}

func (s *_seqSub) SetParam(name string, value any) {
	if v, ok := getStringParam("id", name, value); ok {
		s.Id = v
	} else if v, ok := getStringParam("topic", name, value); ok {
		s.Topic = v
	}
}

func (s *_seqSub) handleString(topic string, value string) {
	s.state.handleString(topic, value)
}

type _seqUnsub struct {
	Id string
}

func (s *_seqUnsub) Step(state *seqState) {
	if sub, ok := getData[Subscription](state, s.Id); ok {
		sub.Unsub()
	}
}

func (s *_seqUnsub) SetParam(name string, value any) {
	if v, ok := getStringParam("id", name, value); ok {
		s.Id = v
	}
}

type _seqChannel struct {
	Id    string
	Topic string
}

func (s *_seqChannel) Step(state *seqState) {
	c := NewChannel[string](state.r, s.Topic)
	state.addData(s.Id, c)
}

func (s *_seqChannel) SetParam(name string, value any) {
	if v, ok := getStringParam("id", name, value); ok {
		s.Id = v
	} else if v, ok := getStringParam("topic", name, value); ok {
		s.Topic = v
	}
}

type _seqPub struct {
	Topic string
	Value string
}

func (s *_seqPub) Step(state *seqState) {
	Pub(state.r, s.Topic, s.Value)
}

func (s *_seqPub) SetParam(name string, value any) {
	if v, ok := getStringParam("topic", name, value); ok {
		s.Topic = v
	} else if v, ok := getStringParam("value", name, value); ok {
		s.Value = v
	}
}

type _seqPubChannel struct {
	Id    string
	Value string
}

func (s *_seqPubChannel) Step(state *seqState) {
	if c, ok := getData[Channel[string]](state, s.Id); ok {
		c.Pub(s.Value)
	}
}

func (s *_seqPubChannel) SetParam(name string, value any) {
	if v, ok := getStringParam("id", name, value); ok {
		s.Id = v
	} else if v, ok := getStringParam("value", name, value); ok {
		s.Value = v
	}
}

type _seqWant struct {
	Want []any
}

func (s *_seqWant) Step(state *seqState) {
	state.want = s.Want
}

func (s *_seqWant) SetParam(name string, value any) {
	if v, ok := value.(string); ok {
		s.Want = append(s.Want, capture{Topic: name, Value: v})
	}
}

// ---------------------------------------------------------
// LOADING

func loadSeqTests(name string) []seqTest {
	var lexer scanner.Scanner
	lexer.Init(strings.NewReader(load(name)))
	lexer.Whitespace ^= 1 << '\n'
	lexer.Mode = scanner.ScanChars | scanner.ScanComments | scanner.ScanFloats | scanner.ScanIdents | scanner.ScanInts | scanner.ScanRawStrings | scanner.ScanStrings
	lexer.Error = func(s *scanner.Scanner, msg string) {
		oferrors.LogFatal(fmt.Errorf("loadSeqTests error: %v", msg))
	}
	b := &seqTestBuilder{}
	for tok := lexer.Scan(); tok != scanner.EOF; tok = lexer.Scan() {
		switch tok {
		case '\n':
			b.flush()
		default:
			b.handle(lexer.TokenText())
		}
	}
	return b.make()
}

func load(name string) string {
	return loadDataFile(os.DirFS("."), name)
}

func loadDataFile(fsys fs.FS, name string) string {
	dat, err := fs.ReadFile(fsys, path.Join("testdata", name))
	oferrors.LogFatal(err)
	return string(dat)
}

type seqTestBuilder struct {
	curTest     *seqTest
	curStep     seqStep
	curKey      string
	needsAssign bool

	built []seqTest
}

func (b *seqTestBuilder) handle(t string) {
	if b.curTest == nil {
		b.curTest = &seqTest{}
	}
	if b.curStep == nil {
		b.handleNewStep(t)
	} else {
		switch t {
		case "=":
			b.needsAssign = true
		default:
			if b.needsAssign {
				if b.curKey == "" {
					oferrors.LogFatal(fmt.Errorf("assign with no LHS"))
				}
				key := strings.Trim(b.curKey, "\"")
				value := strings.Trim(t, "\"")
				b.curStep.SetParam(key, value)
				b.curKey = ""
				b.needsAssign = false
			} else {
				b.curKey = strings.ToLower(t)
			}
		}
	}
}

func (b *seqTestBuilder) handleNewStep(t string) {
	if fn, ok := newSeqStepMap[strings.ToLower(t)]; ok {
		b.curStep = fn()
	} else {
		oferrors.LogFatal(fmt.Errorf("No seq step named %v", t))
	}
	b.curKey = "id"
}

func (b *seqTestBuilder) flush() {
	// Double new lines start a new test, so flush the current
	// step on one new line, and the current test on two.
	if b.curStep != nil {
		b.curTest.Steps = append(b.curTest.Steps, b.curStep)
		b.curStep = nil
	} else if b.curTest != nil {
		b.built = append(b.built, *b.curTest)
		b.curTest = nil
	}
	b.needsAssign = false
}

func (b *seqTestBuilder) make() []seqTest {
	b.flush()
	b.flush()
	return b.built
}

// ---------------------------------------------------------
// HANDLERS

type captureSubscription struct {
	captured      []any
	namerCaptured []Namer
}

func (s *captureSubscription) receiveInt(topic string, value int) {
	s.captured = append(s.captured, value)
}

func (s *captureSubscription) receiveNamer(topic string, value Namer) {
	s.namerCaptured = append(s.namerCaptured, value)
}

// ---------------------------------------------------------
// FUNC

type newSeqStepFunc func() seqStep

// ---------------------------------------------------------
// CONST and VAR

var newSeqStepMap = map[string]newSeqStepFunc{
	"name":       func() seqStep { return &_seqName{} },
	"sub":        func() seqStep { return &_seqSub{} },
	"unsub":      func() seqStep { return &_seqUnsub{} },
	"pub":        func() seqStep { return &_seqPub{} },
	"channel":    func() seqStep { return &_seqChannel{} },
	"pubchannel": func() seqStep { return &_seqPubChannel{} },
	"want":       func() seqStep { return &_seqWant{} },
}
