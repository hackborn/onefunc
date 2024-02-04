package pipeline

import (
	"fmt"
	"os"
	"slices"
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
		{`graph ( na s -> nb )`, `graph (na s -> nb)`, nil},
		{`graph ( na/s -> nb )`, `graph (na/s -> nb)`, nil},
		{`graph ( na -pa> nb )`, `graph (na -pa> nb)`, nil},
		{`graph ( na/a -> nb )`, `graph (na/a -> nb)`, nil},
		{`graph (na(S=f))`, `graph (na) vars (na/S=f)`, nil},
		{`graph (na/a na/b)`, `graph (na/a na/b)`, nil},
		{`graph (na(S="f"))`, `graph (na) vars (na/S=f)`, nil},
		{`graph (na(S=$cat))`, `graph (na) vars (na/S=$cat) env ($cat)`, nil},
		{`graph (na1 -> na3 na2->na3 )`, `graph (na1 -> na3 na2 -> na3)`, nil},
		{`graph (na1/1(S=a) na1/2 )`, `graph (na1/1 na1/2) vars (na1/1/S=a)`, nil},
		{`graph (na1 -> na3(S=a) na2 -> na3 )`, `graph (na1 -> na3 na2 -> na3) vars (na3/S=a)`, nil},
		{`graph (na1 -> na3(S=a) na2 -> na3(S=b) )`, `graph (na1 -> na3 na2 -> na3) vars (na3/S=b)`, nil},
		{`graph (na) env (Path=$Path)`, `graph (na) env (Path=$Path)`, nil},
		// Errors
		{`graph (na`, ``, newSyntaxError("")},
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
		env      map[string]any
		want     []string
		wantErr  error
	}{
		{`graph (na(S=!))`, ``, nil, []string{`!`}, nil},
		{`graph (na(S=!))`, `hi`, nil, []string{`hi!`}, nil},
		// XXX This doesn't work because we aren't generating unique names for nodes, but clearly this should be supported
		//		{`graph (na(S=!) -> na(S=?))`, ``, nil, []string{`!?`}, nil},
		{`graph (na/a(S=!) -> na/b(S=?))`, ``, nil, []string{`!?`}, nil},
		{`graph (na/a(S=!) na/b(S=?))`, ``, nil, []string{`!`, `?`}, nil},
		{`graph (na/a(S=a) na/b(S=b) na/c(S=c) na/d(S=d))`, ``, nil, []string{`a`, `b`, `c`, `d`}, nil},
		{`graph (nb(S1=a, S2=b))`, ``, nil, []string{`ab`}, nil},
		{`graph (nc(S=!))`, `hi`, nil, []string{`hi!`}, nil},
		{`graph (na(S=x) -> nc(S=y))`, `hi`, nil, []string{`hixy`}, nil},
		{`graph (na(S=$env))`, `hi`, map[string]any{`$env`: `!`}, []string{`hi!`}, nil},
		{`graph (na1(S=a) -> na3(S=!) )`, ``, nil, []string{`a!`}, nil},
		{`graph (na1(S=a) -> na3(S=!) na2(S=b) -> na3(S=!))`, ``, nil, []string{`a!`, `b!`}, nil},
		{`graph (na1(S=a) -> na3(S=!) na2(S=b) -> na3)`, ``, nil, []string{`a!`, `b!`}, nil},
		{`graph (na1(S=a) -> na3(S=a) na2(S=b) -> na3(S=!))`, ``, nil, []string{`a!`, `b!`}, nil},
		// ERRORS
		// No env var
		{`graph (na(S=$env))`, `hi`, nil, []string{}, fmt.Errorf("missing env var")},
	}
	for i, v := range table {
		have, haveErr := runAsString(v.pipeline, v.input, v.env)

		if v.wantErr == nil && haveErr != nil {
			t.Fatalf("TestRunString %v expected no error but has %v", i, haveErr)
		} else if v.wantErr != nil && haveErr == nil {
			t.Fatalf("TestRunString %v has no error but exptected %v", i, v.wantErr)
		} else if slices.Compare(have, v.want) != 0 {
			t.Fatalf("TestRunString %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		}
	}
}

func runAsString(expr, input string, env map[string]any) ([]string, error) {
	pin := Pin{Payload: &stringData{s: input}}
	ri := NewInput(pin)
	ro, err := RunExpr(expr, &ri, env)
	if err != nil {
		return nil, err
	}
	if ro == nil || len(ro.Pins) < 1 {
		return nil, nil
	}
	var ans []string
	for _, pin := range ro.Pins {
		switch pt := pin.Payload.(type) {
		case *stringData:
			ans = append(ans, pt.s)
		}
	}
	return ans, nil
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
// BENCHMARK-RUN-AS-STRING
func BenchmarkRunAsString(b *testing.B) {
	const input string = `graph (na(S=!))`
	for n := 0; n < b.N; n++ {
		runAsString(input, "hi", nil)
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

// ---------------------------------------------------------
// NODES

// nodeNa has a string value and adds it to any incoming stringData.
type nodeNa struct {
	S string
}

func (n *nodeNa) Run(s *State, input RunInput) (*RunOutput, error) {
	out := RunOutput{}
	//	fmt.Println("Input len", len(input.Pins))
	for _, p := range input.Pins {
		switch pt := p.Payload.(type) {
		case *stringData:
			out.Pins = append(out.Pins, Pin{Payload: &stringData{s: pt.s + n.S}})
		}
	}
	return &out, nil
}

// nodeNb has multiple string values that are appended to incoming stringData.
type nodeNb struct {
	S1 string
	S2 string
}

func (n *nodeNb) Run(s *State, input RunInput) (*RunOutput, error) {
	out := RunOutput{}
	for _, p := range input.Pins {
		switch pt := p.Payload.(type) {
		case *stringData:
			out.Pins = append(out.Pins, Pin{Payload: &stringData{s: pt.s + n.S1 + n.S2}})
		}
	}
	return &out, nil
}

// nodeNc accumulates string values without producing output.
// When it receives a FlushData{}, it adds its string value to
// the accumulation and sends out data.
type nodeNc struct {
	S string

	accum string
}

func (n *nodeNc) Run(s *State, input RunInput) (*RunOutput, error) {
	for _, p := range input.Pins {
		switch pt := p.Payload.(type) {
		case *stringData:
			n.accum += pt.s
		}
	}
	return nil, nil
}

func (n *nodeNc) Flush(s *State) (*RunOutput, error) {
	out := RunOutput{}
	out.Pins = append(out.Pins, Pin{Payload: &stringData{s: n.accum + n.S}})
	return &out, nil
}

// ---------------------------------------------------------
// LIFECYCLE

func setupTests() {
	RegisterNode("na", func() Runner {
		return &nodeNa{}
	})
	RegisterNode("nb", func() Runner {
		return &nodeNb{}
	})
	RegisterNode("nc", func() Runner {
		return &nodeNc{}
	})

	// Aliases
	RegisterNode("na1", func() Runner {
		return &nodeNa{}
	})
	RegisterNode("na2", func() Runner {
		return &nodeNa{}
	})
	RegisterNode("na3", func() Runner {
		return &nodeNa{}
	})
}

func shutdownTests() {
	reg = newRegistry()
}

var genErr = fmt.Errorf("generic")
