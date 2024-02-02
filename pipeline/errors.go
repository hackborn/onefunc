package pipeline

import (
	"fmt"
)

func newSyntaxError(msg string) error {
	return fmt.Errorf("pipeline syntax error: %v", msg)
}
