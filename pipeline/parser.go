package pipeline

import (
	"fmt"
	"strings"
	"text/scanner"
	"unicode"

	"github.com/hackborn/onefunc/errors"
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
	lexer.IsIdentRune = p.isIdentRune
	lexer.Error = func(s *scanner.Scanner, msg string) {
		p.AddError(fmt.Errorf("scan error: %v", msg))
	}
	tt := token{}
	for tok := lexer.Scan(); tok != scanner.EOF; tok = lexer.Scan() {
		if p.Err != nil {
			return p.Err
		}
		// fmt.Println("TOK", tok, "name", scanner.TokenString(tok), "text", lexer.TokenText())
		tt.text = lexer.TokenText()
		switch tok {
		case scanner.Float:
			tt.tt = floatToken
		case scanner.Int:
			tt.tt = intToken
		case scanner.String:
			tt.tt = stringToken
			tt.text = strings.Trim(tt.text, `"`)
		case scanner.Ident:
			tt.tt = identToken
		case ' ', '\r', '\t', '\n':
			tt.tt = whitespaceToken
			tt.text = ""
		default:
			tt.tt = stringToken
		}
		h.HandleToken(tt)
	}
	return p.Err
}

func (p *parser) isIdentRune(ch rune, i int) bool {
	// This is the standard text scanner ident rune, plus "$" at the start for env vars.
	ident := ch == '_' || unicode.IsLetter(ch) || (unicode.IsDigit(ch) && i > 0) || (ch == '$' && i == 0)
	return ident
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

// graphHandler handles creating nodes, pins, and connecting them.
type graphHandler struct {
	base *baseHandler

	handledOpen bool
	current     currentObj
	currentName string

	nodePushed nodePushedFunc

	// Store existing node names so I don't reuse.
	nodeNames map[string]*astNode
}

func (h *graphHandler) HandleToken(t token) {
	if t.tt == whitespaceToken {
		h.flush()
		return
	}
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
		if h.current.node == nil {
			h.base.AddError(fmt.Errorf("illegal syntax, missing node before pin"))
			return
		}
		h.base.stack = append(h.base.stack, &pinHandler{base: h.base, currentNode: h.current.node})
		h.current = currentObj{}
	case "->":
		h.flush()
		if h.current.node == nil {
			h.base.AddError(fmt.Errorf("illegal syntax, missing node before pin"))
			return
		}
		h.pushNewPin(&astPin{fromNode: h.current.node.nodeName})
	default:
		h.currentName += t.text
	}
}

func (h *graphHandler) HandleVars(s varState) {
	if h.current.node != nil {
		h.handleVarsOnNode(s, h.current.node)
	} else if h.current.pin != nil {
		h.handleVarsOnPin(s, h.current.pin)
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
	h.current = currentObj{}
	h.currentName = ""
	h.nodePushed = h.nullNodePushed
	h.nodeNames = make(map[string]*astNode)
}

func (h *graphHandler) flush() {
	if h.currentName != "" {
		h.pushNewNode(&astNode{nodeName: h.currentName})
		h.currentName = ""
	}
}

func (h *graphHandler) pushNewNode(node *astNode) {
	if found, ok := h.nodeNames[node.nodeName]; ok {
		node = found
	} else {
		h.nodeNames[node.nodeName] = node
		h.base.nodes = append(h.base.nodes, node)
	}

	h.current = currentObj{node: node}
	h.nodePushed(node)
	h.nodePushed = h.nullNodePushed
}

func (h *graphHandler) pushNewPin(pin *astPin) {
	h.current = currentObj{pin: pin}
	h.base.pins = append(h.base.pins, pin)
	h.nodePushed = func(n *astNode) {
		// This happens when a node is pushed that the current
		// pin is connected to. Hook up the connection.
		pin.toNode = n.nodeName
	}
}

func (h *graphHandler) nullNodePushed(n *astNode) {
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
		g.pushNewPin(pin)
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
		s := varState{current: h.current, vars: h.vars, envVars: h.envVars}
		h.base.pop(func(h tokenHandler) { h.HandleVars(s) })
	case "=":
		h.needsCurrent = false
	default:
		if h.needsCurrent {
			h.current = t.text
		} else {
			if strings.HasPrefix(t.text, "$") {
				h.envVars[h.current] = t.text
			} else {
				h.vars[h.current] = t.text
			}
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

// ------------------------------------------------------------
// TOKEN HANDLING TYPES

type varState struct {
	current string // The current, incomplete token, not-yet added to a map
	vars    map[string]any
	envVars map[string]string
}

// currentObj tracks the most recent object on the graph handler
// The node and pin are mutually exclusive.
type currentObj struct {
	node *astNode
	pin  *astPin
}

// ------------------------------------------------------------
// TOKEN HANDLING FUNCS

// poppedFunc is called on a handler when the handler below it
// is popped from the stack.
type poppedFunc func(h tokenHandler)

// nodePushedFunc is called when a new node is pushed on the graph.
type nodePushedFunc func(n *astNode)
