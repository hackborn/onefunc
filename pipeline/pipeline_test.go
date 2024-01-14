package pipeline

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------
// TEST-PARSER-SCAN
func TestParserScan(t *testing.T) {
	table := []struct {
		pipeline string
		want     string
		wantErr  error
	}{
		{`graph (na -> nb)`, `_:graph s:( _:na s:-> _:nb s:)`, nil},
		{`graph(na -> nb)`, `_:graph s:( _:na s:-> _:nb s:)`, nil},
		{`graph(na->nb)`, `_:graph s:( _:na s:-> _:nb s:)`, nil},
		//		{`graph ( na -> nb )`, `graph (na -> nb)`, nil},
		//		{`graph ( na s -> nb )`, `graph (nas -> nb)`, nil},
		//		{`graph ( na/s -> nb )`, `graph (na/s -> nb)`, nil},
		//		{`graph ( na -pa> nb )`, `graph (na -pa> nb)`, nil},
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
		{`graph (na -> nb)`, `graph (na -> nb)`, nil},
		{`graph ( na -> nb )`, `graph (na -> nb)`, nil},
		{`graph ( na s -> nb )`, `graph (nas -> nb)`, nil},
		{`graph ( na/s -> nb )`, `graph (na/s -> nb)`, nil},
		{`graph ( na -pa> nb )`, `graph (na -pa> nb)`, nil},
	}
	for i, v := range table {
		ast, haveErr := parse(v.pipeline)
		have := ast.print()

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
		h.b.WriteString("_")
	default:
		h.b.WriteString("?")
	}
	h.b.WriteString(":")
	h.b.WriteString(t.text)
}
