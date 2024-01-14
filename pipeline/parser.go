package pipeline

import (
	"fmt"
	"strings"
	"text/scanner"

	"github.com/hackborn/onefunc/errors"
	ofstrings "github.com/hackborn/onefunc/strings"
)

type tokenType int

const (
	stringToken tokenType = iota
	floatToken
	intToken
	identToken
	whitespaceToken
	charToken
)

type token struct {
	tt   tokenType
	text string
}

func parse(input string) (astPipeline, error) {
	p := newParser()
	return p.parse(input)
}

func newParser() *parser {
	return &parser{}
}

type parser struct {
	errors.FirstBlock
}

func (p *parser) parse(input string) (astPipeline, error) {
	h := &baseHandler{}
	err := p.scan(input, h)
	if err != nil {
		return astPipeline{}, err
	}
	return h.astPipeline, nil
}

// scan lexes the input string, sending the tokens to the handler.
func (p *parser) scan(input string, h tokenHandler) error {
	var lexer scanner.Scanner
	lexer.Init(strings.NewReader(input))
	lexer.Whitespace = 0
	lexer.Mode = scanner.ScanChars | scanner.ScanComments | scanner.ScanFloats | scanner.ScanIdents | scanner.ScanInts | scanner.ScanRawStrings | scanner.ScanStrings

	accum := &accumHandler{h: h}
	lexer.Error = func(s *scanner.Scanner, msg string) {
		p.AddError(fmt.Errorf("scan error: %v", msg))
	}
	for tok := lexer.Scan(); tok != scanner.EOF; tok = lexer.Scan() {
		if p.Err != nil {
			return p.Err
		}
		// fmt.Println("TOK", tok, "name", scanner.TokenString(tok), "text", lexer.TokenText())
		switch tok {
		case scanner.Float:
			accum.HandleToken(token{tt: floatToken, text: lexer.TokenText()})
		case scanner.Int:
			accum.HandleToken(token{tt: intToken, text: lexer.TokenText()})
		case scanner.String:
			accum.HandleToken(token{tt: stringToken, text: lexer.TokenText()})
		case scanner.Ident:
			accum.HandleToken(token{tt: identToken, text: lexer.TokenText()})
		case ' ', '\r', '\t', '\n':
			accum.HandleToken(token{tt: whitespaceToken})
		default:
			accum.HandleToken(token{tt: charToken, text: lexer.TokenText()})
		}
	}
	accum.flush()
	return p.Err
}

// accumHandler provides additional tokenizing on top of the built-in scanner.
type accumHandler struct {
	h tokenHandler
	s string
}

func (h *accumHandler) HandleToken(t token) {
	switch t.tt {
	case whitespaceToken:
		h.flush()
	case charToken:
		h.s += t.text
	default:
		h.flush()
		h.h.HandleToken(t)
	}
}

func (h *accumHandler) flush() {
	if h.s != "" {
		h.h.HandleToken(token{tt: stringToken, text: h.s})
		h.s = ""
	}
}

// ------------------------------------------------------------
// AST

type astPipeline struct {
	nodes []*astNode
	pins  []*astPin
}

func (t astPipeline) print() string {
	eb := errors.FirstBlock{}
	w := ofstrings.GetWriter(&eb)
	defer ofstrings.PutWriter(w)

	w.WriteString("graph (")
	lastNode := ""
	for _, pin := range t.pins {
		if lastNode != pin.fromNode {
			w.WriteString(pin.fromNode)
		}
		w.WriteString(" -")
		w.WriteString(pin.pinName)
		w.WriteString("> ")
		w.WriteString(pin.toNode)
		lastNode = pin.toNode
	}
	w.WriteString(")")
	return ofstrings.String(w)
}

type astNode struct {
	nodeName string
}

type astPin struct {
	pinName          string
	fromNode, toNode string
}

// ------------------------------------------------------------
// TOKEN HANDLING

type tokenHandler interface {
	HandleToken(t token)
}

// baseHandler supplies the rules for turning tokens into AST nodes.
type baseHandler struct {
	astPipeline
	errors.FirstBlock
	stack []tokenHandler
}

func (h *baseHandler) AddError(e error) {
	h.FirstBlock.AddError(e)
}

func (h *baseHandler) HandleToken(t token) {
	size := len(h.stack)
	if size > 0 {
		h.stack[size-1].HandleToken(t)
		return
	}

	txt := strings.ToLower(t.text)
	switch txt {
	case "graph":
		h.stack = append(h.stack, &graphHandler{base: h})
	}
}

func (h *baseHandler) pop() {
	s := len(h.stack)
	if s > 0 {
		h.stack = h.stack[:s-1]
	}
}

func (h *baseHandler) flush() {
}

// graphHandler
type graphHandler struct {
	base *baseHandler

	currentNodeName string
	currentNode     *astNode
}

func (h *graphHandler) HandleToken(t token) {
	txt := strings.ToLower(t.text)
	switch txt {
	case "(":
		return
	case ")":
		h.flush()
		h.base.pop()
	case "-":
		h.flush()
		if h.currentNode == nil {
			h.base.AddError(fmt.Errorf("illegal syntax, missing node before pin"))
			return
		}
		h.base.stack = append(h.base.stack, &pinHandler{base: h.base, currentNode: h.currentNode})
		h.currentNode = nil
	case "->":
		h.flush()
		if h.currentNode == nil {
			h.base.AddError(fmt.Errorf("illegal syntax, missing node before pin"))
			return
		}
		h.base.pins = append(h.base.pins, &astPin{fromNode: h.currentNode.nodeName})
	default:
		h.currentNodeName += t.text
	}
}

func (h *graphHandler) flush() {
	if h.currentNodeName != "" {
		h.currentNode = &astNode{nodeName: h.currentNodeName}
		h.currentNodeName = ""
		h.base.nodes = append(h.base.nodes, h.currentNode)
		size := len(h.base.pins)
		if size > 0 {
			h.base.pins[size-1].toNode = h.currentNode.nodeName
		}
	}
}

// nodeHandler
type nodeHandler struct {
	base *baseHandler
}

func (h *nodeHandler) HandleToken(t token) {
}

// pinHandler
type pinHandler struct {
	base *baseHandler

	currentNode *astNode
	pinName     string
}

func (h *pinHandler) HandleToken(t token) {
	txt := strings.ToLower(t.text)
	switch txt {
	case ">":
		h.base.pins = append(h.base.pins, &astPin{pinName: h.pinName, fromNode: h.currentNode.nodeName})
		h.base.pop()
	default:
		h.pinName += t.text
	}
}
