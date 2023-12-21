package assign

import (
	"fmt"
)

type ValuesRequest struct {
	FieldNames []string
	NewValues  []any
}

func (r ValuesRequest) Validate() error {
	if len(r.FieldNames) != len(r.NewValues) {
		return fmt.Errorf("Size mismatch (%v names but %v values)", len(r.FieldNames), len(r.NewValues))
	}
	return nil
}
