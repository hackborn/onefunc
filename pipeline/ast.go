package pipeline

import (
	"fmt"

	"github.com/hackborn/onefunc/errors"
	ofstrings "github.com/hackborn/onefunc/strings"
)

// astNode stores an abstract node from parse.
type astNode struct {
	nodeName string
	vars     map[string]any
	envVars  map[string]string
}

// astPin stores an abstract pin from a parse.
type astPin struct {
	pinName          string
	fromNode, toNode string
}

// astPipeline stores an abstract pipeline from a parse.
type astPipeline struct {
	nodes []*astNode
	pins  []*astPin
	env   map[string]any
}

func (t astPipeline) print() string {
	if len(t.nodes) < 1 && len(t.pins) < 1 {
		return ""
	}
	eb := errors.FirstBlock{}
	w := ofstrings.GetWriter(&eb)
	defer ofstrings.PutWriter(w)

	w.WriteString("graph (")
	// Display unattached nodes
	needsSpace := false
	unattached := t.unattachedNodes()
	for i, node := range unattached {
		if i > 0 {
			w.WriteString(" ")
		}
		w.WriteString(node.nodeName)
		needsSpace = true
	}
	// Display connected nodes
	lastNode := ""
	for _, pin := range t.pins {
		if needsSpace {
			needsSpace = false
			w.WriteString(" ")
		}
		if lastNode != pin.fromNode {
			if lastNode != "" {
				w.WriteString(" ")
			}
			w.WriteString(pin.fromNode)
		}
		w.WriteString(" -")
		w.WriteString(pin.pinName)
		w.WriteString("> ")
		w.WriteString(pin.toNode)
		lastNode = pin.toNode
	}
	w.WriteString(")")
	// Vars
	first := true
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
	// Env
	first = true
	for _, node := range t.nodes {
		if len(node.envVars) < 1 {
			continue
		}
		if first {
			w.WriteString(" env (")
			first = false
		} else {
			w.WriteString(", ")
		}
		needsComma := false
		for _, v := range node.envVars {
			if !needsComma {
				needsComma = true
			} else {
				w.WriteString(", ")
			}
			w.WriteString(fmt.Sprintf("%v", v))
		}
	}
	if !first {
		w.WriteString(")")
	}

	// Env
	if len(t.env) > 0 {
		if ofstrings.StringLen(w) > 0 {
			w.WriteString(" ")
		}
		w.WriteString("env (")
		first := true
		for k, v := range t.env {
			if !first {
				w.WriteString(", ")
				first = true
			}
			w.WriteString(k)
			w.WriteString("=")
			w.WriteString(fmt.Sprintf("%v", v))
		}
		w.WriteString(")")
	}

	return ofstrings.String(w)
}

func (t astPipeline) unattachedNodes() []*astNode {
	m := make(map[*astNode]struct{})
	for _, n := range t.nodes {
		m[n] = struct{}{}
	}
	for _, pin := range t.pins {
		for k, _ := range m {
			if pin.fromNode == k.nodeName || pin.toNode == k.nodeName {
				delete(m, k)
			}
		}
	}
	// Keep initial order preserved
	var ans []*astNode
	for _, n := range t.nodes {
		if _, ok := m[n]; ok {
			ans = append(ans, n)
		}
	}
	return ans
}
