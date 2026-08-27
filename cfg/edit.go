package cfg

import (
	"cmp"
	"encoding/json"
	"strings"
)

// Add returns a new settings that is a copy of s but
// adds the desired value at the given key.
func Add(s Settings, key string, value any) (Settings, error) {
	if s.t == nil {
		s.t = make(map[string]any)
	} else {
		t := make(map[string]any)
		dat, err := s.asJson()
		if err = cmp.Or(err, json.Unmarshal(dat, &t)); err != nil {
			return s, err
		}
		s.t = t
		cur := t
		for part := range strings.SplitSeq(key, pathSeparator) {
			if strings.HasSuffix(key, part) {
				cur[part] = value
			} else {
				next := make(map[string]any)
				cur[part] = next
				cur = next
			}
		}
	}
	return s, nil
}
