package nodes

import (
	"fmt"
	"regexp"
	"strings"

	oferrors "github.com/hackborn/onefunc/errors"
	"github.com/hackborn/onefunc/pipeline"
)

type RegexpNode struct {
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

func (n *RegexpNode) Run(s *pipeline.State, input pipeline.RunInput) (*pipeline.RunOutput, error) {
	re, runFn, setFn, err := n.prepare()
	if err != nil {
		return nil, err
	}
	output := pipeline.RunOutput{}
	output.Pins = make([]pipeline.Pin, 0, len(input.Pins))
	eb := &oferrors.FirstBlock{}
	for _, pin := range input.Pins {
		eb.AddError(setFn(pin.Payload, runFn, re, n))
		output.Pins = append(output.Pins, pin)
	}
	return &output, eb.Err
}

func (n *RegexpNode) prepare() (*regexp.Regexp, regexpOperationFn, regexpTargetFn, error) {
	if n.Expr == "" {
		return nil, nil, nil, fmt.Errorf("regexp node: No expression")
	}
	re, err := regexp.Compile(n.Expr)
	if err != nil {
		return nil, nil, nil, err
	}
	runFn, ok := regexpOperations[strings.ToLower(n.Operation)]
	if !ok {
		return nil, nil, nil, fmt.Errorf("regexp node: No operation named \"%v\"", n.Operation)
	}
	setFn, ok := regexpTargets[strings.ToLower(n.Target)]
	if !ok {
		return nil, nil, nil, fmt.Errorf("regexp node: No target named \"%v\"", n.Target)
	}
	return re, runFn, setFn, nil
}
