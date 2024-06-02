package values

import (
	"fmt"
)

const (
	FuzzyFloats = 1 << iota
	FuzzyInts
	FuzzyStrings
	Fuzzy = FuzzyFloats | FuzzyInts | FuzzyStrings
)

var (
	mustBePointerErr = fmt.Errorf("must have pointer destination")
)
