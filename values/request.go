package values

import (
	"fmt"
)

type SetRequest struct {
	FieldNames []string
	NewValues  []any
	Assigns    []SetFunc // Optional: Provide assignment func for each value. Nil will use default.
	Flags      uint8
}

func SetRequestFrom(m map[string]any) SetRequest {
	r := SetRequest{}
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

func (r SetRequest) Validate() error {
	if len(r.FieldNames) != len(r.NewValues) {
		return fmt.Errorf("Size mismatch (%v names but %v values)", len(r.FieldNames), len(r.NewValues))
	}
	return nil
}
