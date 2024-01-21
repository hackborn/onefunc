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

	//	accum := &accumHandler{h: h}
	accum := h
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
			accum.HandleToken(token{tt: stringToken, text: lexer.TokenText()})
		}
	}
	return p.Err
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
	// Hack to display nodes with no pins, but need to figute
	// a better way to reconstruct the input
	if len(t.pins) < 1 {
		for i, node := range t.nodes {
			if i > 0 {
				w.WriteString(" ")
			}
			w.WriteString(node.nodeName)
		}
	}
	w.WriteString(")")
	first := true
	// Vars
	for _, node := range t.nodes {
		if len(node.vars) < 1 && len(node.envVars) < 1 {
			continue
		}
		if first {
			w.WriteString(" vars (")
			first = false
		} else {
			w.WriteString(", ")
		}
		needsComma := false
		for k, v := range node.vars {
			if !needsComma {
				needsComma = true
			} else {
				w.WriteString(", ")
			}
			w.WriteString(node.nodeName + "/" + k)
			w.WriteString("=")
			w.WriteString(fmt.Sprintf("%v", v))
		}
		for k, v := range node.envVars {
			if !needsComma {
				needsComma = true
			} else {
				w.WriteString(", ")
			}
			w.WriteString(node.nodeName + "/" + k)
			w.WriteString("=")
			w.WriteString(fmt.Sprintf("%v", v))
		}
	}
	if !first {
		w.WriteString(")")
	}

	return ofstrings.String(w)
}

type astNode struct {
	nodeName string
	vars     map[string]any
	envVars  map[string]string
}

type astPin struct {
	pinName          string
	fromNode, toNode string
}

// ------------------------------------------------------------
// TOKEN HANDLING

type tokenHandler interface {
	HandleToken(t token)
	HandleVars(s varState)
	Pushed(base *baseHandler)
}

// baseHandler supplies the rules for turning tokens into AST nodes.
type baseHandler struct {
	astPipeline
	errors.FirstBlock
	stack []tokenHandler

	// Cached handlers. Push these onto the stack when needed.
	graph graphHandler
	vars  varsHandler
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
		h.push(&h.graph)
	}
}

func (h *baseHandler) HandleVars(s varState) {
}

func (h *baseHandler) Pushed(base *baseHandler) {
	h.AddError(fmt.Errorf("Illegal push on base handler"))
}

func (h *baseHandler) push(th tokenHandler) {
	h.stack = append(h.stack, th)
	th.Pushed(h)
}

func (h *baseHandler) pop(fn poppedFunc) {
	s := len(h.stack)
	if s > 0 {
		h.stack = h.stack[:s-1]
	}
	if fn != nil && len(h.stack) > 0 {
		fn(h.stack[len(h.stack)-1])
	}
}

func (h *baseHandler) flush() {
}

// graphHandler
type graphHandler struct {
	base *baseHandler

	handledOpen     bool
	currentNodeName string
	currentNode     *astNode
	currentAst      any // will be *astNode or *astPin (or nil)
}

func (h *graphHandler) HandleToken(t token) {
	txt := strings.ToLower(t.text)
	switch txt {
	case "(":
		if !h.handledOpen {
			h.handledOpen = true
			return
		}
		h.flush()
		h.base.push(&h.base.vars)
		return
	case ")":
		h.flush()
		h.base.pop(nil)
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
		pin := &astPin{fromNode: h.currentNode.nodeName}
		h.base.pins = append(h.base.pins, pin)
		h.currentAst = pin
	default:
		h.currentNodeName += t.text
	}
}

func (h *graphHandler) HandleVars(s varState) {
	switch t := h.currentAst.(type) {
	case *astNode:
		h.handleVarsOnNode(s, t)
	case *astPin:
		h.handleVarsOnPin(s, t)
	}
}

func (h *graphHandler) handleVarsOnNode(s varState, n *astNode) {
	n.vars = s.vars
	n.envVars = s.envVars
}

func (h *graphHandler) handleVarsOnPin(s varState, n *astPin) {
}

func (h *graphHandler) Pushed(base *baseHandler) {
	h.base = base
	h.handledOpen = false
	h.currentAst = nil
}

func (h *graphHandler) flush() {
	if h.currentNodeName != "" {
		h.currentNode = &astNode{nodeName: h.currentNodeName}
		h.currentAst = h.currentNode
		h.currentNodeName = ""
		h.base.nodes = append(h.base.nodes, h.currentNode)
		size := len(h.base.pins)
		if size > 0 {
			h.base.pins[size-1].toNode = h.currentNode.nodeName
		}
	}
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
		pin := &astPin{pinName: h.pinName, fromNode: h.currentNode.nodeName}
		h.base.pins = append(h.base.pins, pin)
		h.base.pop(func(t tokenHandler) {
			h.onPoppedFunc(pin, t)
		})
	default:
		h.pinName += t.text
	}
}

func (h *pinHandler) HandleVars(s varState) {
}

func (h *pinHandler) Pushed(base *baseHandler) {
}

func (h *pinHandler) onPoppedFunc(pin *astPin, t tokenHandler) {
	if g, ok := t.(*graphHandler); ok {
		g.currentAst = pin
	}
}

// varsHandler
type varsHandler struct {
	base *baseHandler

	current      string
	needsCurrent bool
	vars         map[string]any
	envVars      map[string]string
}

func (h *varsHandler) HandleToken(t token) {
	txt := strings.ToLower(t.text)
	switch txt {
	case "(":
		return
	case ")":
		h.flush()
		s := varState{current: h.current, vars: h.vars, envVars: h.envVars}
		h.base.pop(func(h tokenHandler) { h.HandleVars(s) })
	case "=":
		h.needsCurrent = false
	default:
		if h.needsCurrent {
			h.current = t.text
		} else {
			h.vars[h.current] = t.text
			h.needsCurrent = true
		}
	}
}

func (h *varsHandler) HandleVars(s varState) {
	h.base.AddError(fmt.Errorf("varsHandler can't handle vars"))
}

func (h *varsHandler) Pushed(base *baseHandler) {
	h.base = base
	h.current = ""
	h.needsCurrent = true
	h.vars = make(map[string]any)
	h.envVars = make(map[string]string)
}

func (h *varsHandler) flush() {
	if h.vars != nil {
		/*
			h.currentNode = &astNode{nodeName: h.currentNodeName}
			h.currentNodeName = ""
			h.base.nodes = append(h.base.nodes, h.currentNode)
			size := len(h.base.pins)
			if size > 0 {
				h.base.pins[size-1].toNode = h.currentNode.nodeName
			}
		*/
	}
}

// ------------------------------------------------------------
// TOKEN HANDLING TYPES

type varState struct {
	current string // The current, incomplete token, not-yet added to a map
	vars    map[string]any
	envVars map[string]string
}

// ------------------------------------------------------------
// TOKEN HANDLING FUNCS

// poppedFunc is called on a handler when the handler below it
// is popped from the stack.
type poppedFunc func(h tokenHandler)
