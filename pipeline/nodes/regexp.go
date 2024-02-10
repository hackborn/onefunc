package nodes

import (
	"fmt"
	"regexp"
	"strings"

	oferrors "github.com/hackborn/onefunc/errors"
	"github.com/hackborn/onefunc/pipeline"
)

type RegexpNode struct {
	regexpData
}

type regexpData struct {
	// Expr is the expression used for the regex matching.
	// See godocs for a description.
	// A quick example, the following strips only the beginning
	// and end from a string if they match the pattern:
	// Given the phrase "allis is tall is"
	// and the expr "(^(all))|((is)$)"
	// the result is "is is tall "
	Expr string

	// Target is the type of data I operate on. Supported:
	// "content.name" -- a ContentData.Name
	Target string

	// The regex operation. Supported:
	// "" -- default to "ReplaceAllString"
	// "replace" -- "ReplaceAllString"
	Operation string

	// The replacement string when performing a "replace" operation.
	Replace string
}

func (n *RegexpNode) Start(state *pipeline.State) {
	data := n.regexpData
	state.NodeData = &data
}

func (n *RegexpNode) Run(state *pipeline.State, input pipeline.RunInput) (*pipeline.RunOutput, error) {
	data := state.NodeData.(*regexpData)
	re, runFn, setFn, err := n.prepare(data)
	if err != nil {
		return nil, err
	}
	output := pipeline.RunOutput{}
	output.Pins = make([]pipeline.Pin, 0, len(input.Pins))
	eb := &oferrors.FirstBlock{}
	for _, pin := range input.Pins {
		eb.AddError(setFn(pin.Payload, runFn, re, data))
		output.Pins = append(output.Pins, pin)
	}
	return &output, eb.Err
}

func (n *RegexpNode) prepare(data *regexpData) (*regexp.Regexp, regexpOperationFn, regexpTargetFn, error) {
	if data.Expr == "" {
		return nil, nil, nil, fmt.Errorf("regexp node: No expression")
	}
	re, err := regexp.Compile(data.Expr)
	if err != nil {
		return nil, nil, nil, err
	}
	runFn, ok := regexpOperations[strings.ToLower(data.Operation)]
	if !ok {
		return nil, nil, nil, fmt.Errorf("regexp node: No operation named \"%v\"", data.Operation)
	}
	setFn, ok := regexpTargets[strings.ToLower(data.Target)]
	if !ok {
		return nil, nil, nil, fmt.Errorf("regexp node: No target named \"%v\"", data.Target)
	}
	return re, runFn, setFn, nil
}
