package pipeline

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hackborn/onefunc/assign"
)

// Compile converts an expression into a Pipeline. It also supplies
// the results of the expressions env() term, if anything -- this
// has no functional impact, but serves as a form of documentation
// to let clients discover what env vars the pipeline supports.
func Compile(expr string) (*Pipeline, error) {
	ast, err := parse(expr)
	if err != nil {
		return nil, err
	}
	pipeline := &Pipeline{env: ast.env}
	nodes := make(map[string]*runningNode)
	roots := make(map[string]compileRoot)
	for i, nn := range ast.nodes {
		splitName := strings.Split(nn.nodeName, "/")
		node, err := newNode(splitName[0])
		if err != nil {
			return nil, err
		}
		rn := &runningNode{node: node, envVars: nn.envVars}
		if starter, ok := node.(Starter); ok {
			if rn.nodeState = starter.StartNodeState(); rn.nodeState != nil {
				rn.hasStartNodeState = true
			}
		}
		if rn.nodeState == nil {
			rn.nodeState = node
		}
		nodes[nn.nodeName] = rn
		roots[nn.nodeName] = compileRoot{index: i, node: rn}
		pipeline.nodes = append(pipeline.nodes, rn)
		// apply fixed vars
		if len(nn.vars) > 0 {
			req := assign.ValuesRequestFrom(nn.vars)
			err = assign.Values(req, rn.nodeState)
			if err != nil {
				return nil, err
			}
		}
	}
	for _, pin := range ast.pins {
		fromNode, _ := nodes[pin.fromNode]
		toNode, _ := nodes[pin.toNode]
		if fromNode == nil {
			return nil, fmt.Errorf("Missing node %v", pin.fromNode)
		}
		if toNode == nil {
			return nil, fmt.Errorf("Missing node %v", pin.toNode)
		}
		delete(roots, pin.toNode)
		rp := &runningPin{toNode: toNode}
		toNode.maxInputCount += 1
		fromNode.output = append(fromNode.output, rp)
	}
	pipeline.roots = compileRoots(roots)
	if len(pipeline.roots) < 1 {
		return nil, fmt.Errorf("No roots")
	}
	return pipeline, nil
}

func compileRoots(mapRoots map[string]compileRoot) []*runningNode {
	// Keep the roots in the same order as the AST nodes. Not
	// strictly necessary, but it does make tests predictable.
	sortedRoots := make([]compileRoot, 0, len(mapRoots))
	for _, r := range mapRoots {
		sortedRoots = append(sortedRoots, r)
	}
	slices.SortFunc(sortedRoots, func(a, b compileRoot) int {
		if a.index < b.index {
			return -1
		} else if a.index > b.index {
			return 1
		} else {
			return 0
		}
	})
	roots := make([]*runningNode, 0, len(sortedRoots))
	for _, r := range sortedRoots {
		roots = append(roots, r.node)
	}
	return roots
}

type compileRoot struct {
	index int
	node  *runningNode
}
