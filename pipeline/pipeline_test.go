package pipeline

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	setupTests()
	code := m.Run()
	shutdownTests()
	os.Exit(code)
}

// ---------------------------------------------------------
// TEST-PARSER-SCAN
func TestParserScan(t *testing.T) {
	table := []struct {
		pipeline string
		want     string
		wantErr  error
	}{
		{`graph (na -> nb)`, `graph s:( na s:- s:> nb s:)`, nil},
		{`graph(na -> nb)`, `graph s:( na s:- s:> nb s:)`, nil},
		{`graph(na->nb)`, `graph s:( na s:- s:> nb s:)`, nil},
		{`graph(na/a->nb)`, `graph s:( na s:/ a s:- s:> nb s:)`, nil},
		{`graph (na (S=add))`, `graph s:( na s:( S s:= add s:) s:)`, nil},
	}
	for i, v := range table {
		p := newParser()
		h := &fmtTokenHandler{}
		haveErr := p.scan(v.pipeline, h)
		have := h.b.String()

		if v.wantErr == nil && haveErr != nil {
			t.Fatalf("TestTokenizer %v expected no error but has %v", i, haveErr)
		} else if v.wantErr != nil && haveErr == nil {
			t.Fatalf("TestTokenizer %v has no error but exptected %v", i, v.wantErr)
		} else if have != v.want {
			t.Fatalf("TestTokenizer %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-PARSER
func TestParser(t *testing.T) {
	table := []struct {
		pipeline string
		want     string
		wantErr  error
	}{
		{`graph (na)`, `graph (na)`, nil},
		{`graph (na -> nb)`, `graph (na -> nb)`, nil},
		{`graph ( na -> nb )`, `graph (na -> nb)`, nil},
		{`graph ( na s -> nb )`, `graph (nas -> nb)`, nil},
		{`graph ( na/s -> nb )`, `graph (na/s -> nb)`, nil},
		{`graph ( na -pa> nb )`, `graph (na -pa> nb)`, nil},
		{`graph ( na/a -> nb )`, `graph (na/a -> nb)`, nil},
		{`graph (na(S=f))`, `graph (na) vars (na/S=f)`, nil},
	}
	for i, v := range table {
		ast, haveErr := parse(v.pipeline)
		have := ast.print()

		if v.wantErr == nil && haveErr != nil {
			t.Fatalf("TestParser %v expected no error but has %v", i, haveErr)
		} else if v.wantErr != nil && haveErr == nil {
			t.Fatalf("TestParser %v has no error but exptected %v", i, v.wantErr)
		} else if have != v.want {
			t.Fatalf("TestParser %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-RUN-STRING
func TestRunString(t *testing.T) {
	table := []struct {
		pipeline string
		input    string
		want     string
		wantErr  error
	}{
		{`graph (na(S=!))`, ``, `!`, nil},
		{`graph (na(S=!))`, `hi`, `hi!`, nil},
		// XXX This doesn't work because we aren't generating unique names for nodes, but clearly this should be supported
		//		{`graph (na(S=!) -> na(S=?))`, ``, `!?`, nil},
		{`graph (na/a(S=!) -> na/b(S=?))`, ``, `!?`, nil},
	}
	for i, v := range table {
		have, haveErr := runAsString(v.pipeline, v.input)

		if v.wantErr == nil && haveErr != nil {
			t.Fatalf("TestRunString %v expected no error but has %v", i, haveErr)
		} else if v.wantErr != nil && haveErr == nil {
			t.Fatalf("TestRunString %v has no error but exptected %v", i, v.wantErr)
		} else if have != v.want {
			t.Fatalf("TestRunString %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		}
	}
}

func runAsString(expr, input string) (string, error) {
	p, err := Compile(expr)
	if err != nil {
		return "", err
	}
	data := &stringData{s: input}
	ri := NewInput(data)
	ro, err := Run(p, &ri)
	if err != nil {
		return "", err
	}
	if ro == nil || len(ro.Pins) < 1 {
		return "", nil
	}
	for _, pin := range ro.Pins {
		switch pt := pin.(type) {
		case *stringData:
			return pt.s, nil
		}
	}
	return "", fmt.Errorf("No stringData pin")
}

// ---------------------------------------------------------
// BENCHMARK-PARSER
func BenchmarkParser(b *testing.B) {
	const input string = `graph (na -> nb)`
	for n := 0; n < b.N; n++ {
		parse(input)
	}
}

// ---------------------------------------------------------
// TOKENS

type fmtTokenHandler struct {
	b strings.Builder
}

func (h *fmtTokenHandler) HandleToken(t token) {
	if t.tt == whitespaceToken {
		return
	}

	if h.b.Len() > 0 {
		h.b.WriteString(" ")
	}
	switch t.tt {
	case stringToken:
		h.b.WriteString("s")
	case floatToken:
		h.b.WriteString("f")
	case intToken:
		h.b.WriteString("i")
	case identToken:
		//	h.b.WriteString("_")
		// No prefix on idents, just to be easier to read
		h.b.WriteString(t.text)
		return
	default:
		h.b.WriteString("?")
	}
	h.b.WriteString(":")
	h.b.WriteString(t.text)
}

func (h *fmtTokenHandler) HandleVars(s varState) {
}

func (h *fmtTokenHandler) Pushed(base *baseHandler) {
}

// ---------------------------------------------------------
// PIN DATA

type stringData struct {
	s string
}

func (d *stringData) PinName() string {
	return ""
}

// ---------------------------------------------------------
// NODES

type nodeNa struct {
	S string
}

func (n *nodeNa) Run(s *State, input RunInput) (*RunOutput, error) {
	out := RunOutput{}
	for _, p := range input.Pins {
		switch pt := p.(type) {
		case *stringData:
			out.Pins = append(out.Pins, &stringData{s: pt.s + n.S})
		}
	}
	return &out, nil
}

// ---------------------------------------------------------
// LIFECYCLE

func setupTests() {
	RegisterNode("na", func() Node {
		return &nodeNa{}
	})
}

func shutdownTests() {
	reg = newRegistry()
}
