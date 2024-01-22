package pipeline

import (
	"fmt"
	"strings"

	"github.com/hackborn/onefunc/assign"
)

func Compile(expr string) (*Pipeline, error) {
	ast, err := parse(expr)
	if err != nil {
		return nil, err
	}
	pipeline := &Pipeline{}
	nodes := make(map[string]*runningNode)
	roots := make(map[string]*runningNode)
	for _, nn := range ast.nodes {
		splitName := strings.Split(nn.nodeName, "/")
		node, err := newNode(splitName[0])
		if err != nil {
			return nil, err
		}
		rn := &runningNode{node: node}
		nodes[nn.nodeName] = rn
		roots[nn.nodeName] = rn
		pipeline.nodes = append(pipeline.nodes, rn)
		// apply fixed vars
		if len(nn.vars) > 0 {
			req := assign.ValuesRequestFrom(nn.vars)
			err = assign.Values(req, node)
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
	if len(roots) < 1 {
		return nil, fmt.Errorf("No roots")
	}
	for _, v := range roots {
		pipeline.roots = append(pipeline.roots, v)
	}
	return pipeline, nil
}
