package nodes

import (
	"regexp"
)

// regexpOperationFn abstracts performing the RegexpNode.Operation.
type regexpOperationFn func(re *regexp.Regexp, s string, n *RegexpNode) (string, error)

// regexpTargetFn abstracts performing the RegexpNode.Target.
type regexpTargetFn func(pin any, fn regexpOperationFn, re *regexp.Regexp, n *RegexpNode) error
