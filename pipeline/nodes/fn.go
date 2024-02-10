package nodes

import (
	"regexp"
)

// regexpOperationFn abstracts performing the RegexpNode.Operation.
type regexpOperationFn func(re *regexp.Regexp, s string, data *regexpData) (string, error)

// regexpTargetFn abstracts performing the RegexpNode.Target.
type regexpTargetFn func(pin any, fn regexpOperationFn, re *regexp.Regexp, data *regexpData) error
