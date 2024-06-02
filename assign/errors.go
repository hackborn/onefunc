package assign

import (
	"fmt"
)

var (
	mustBePointerErr      = fmt.Errorf("must have pointer destination")
	unhandledValueTypeErr = fmt.Errorf("unhandled value type")
)
