package jacl

import (
	"strings"
)

type Opts struct {
	// If RawValues is true then the value tokens are left
	// completely unprocessed.
	// If it is false, then processing is applied:
	// * Two single quotes ('') are replaced with a double quote (").
	// Default is false.
	RawValues bool
}

func (o Opts) processValue(s string) string {
	if o.RawValues {
		return s
	}
	s = strings.ReplaceAll(s, `''`, `"`)
	return s
}
