package assign

import (
	"fmt"
)

type ValuesRequest struct {
	FieldNames []string
	NewValues  []any
}

func ValuesRequestFrom(m map[string]any) ValuesRequest {
	r := ValuesRequest{}
	if len(m) > 0 {
		r.FieldNames = make([]string, 0, len(m))
		r.NewValues = make([]any, 0, len(m))
		for k, v := range m {
			r.FieldNames = append(r.FieldNames, k)
			r.NewValues = append(r.NewValues, v)
		}
	}
	return r
}

func (r ValuesRequest) Validate() error {
	if len(r.FieldNames) != len(r.NewValues) {
		return fmt.Errorf("Size mismatch (%v names but %v values)", len(r.FieldNames), len(r.NewValues))
	}
	return nil
}
