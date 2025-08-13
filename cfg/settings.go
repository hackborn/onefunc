package cfg

import (
	"encoding/json"
	"fmt"
	"iter"
	"os"
	"slices"
	"strconv"
	"strings"

	oferrors "github.com/hackborn/onefunc/errors"
)

type tree = map[string]any

// Settings stores a tree of settings values, accessed
// via a path syntax ("path/to/value").
// Settings are immutable. Modifications must be
// made via creating a new Settings.
type Settings struct {
	t tree
	// sliceKey is a special case where we make a subset to a
	// slice. For the client, the key that got us to the slice
	// has disappeared, but we are a map and need a key, so it gets saved here.
	sliceKey string
}

func NewSettings(opts ...Option) (Settings, error) {
	s := emptySettings()
	eb := &oferrors.FirstBlock{}
	builder := &_builder{t: s.t}
	for _, opt := range opts {
		if opt != nil {
			opt(builder, eb)
		}
	}
	return s, eb.Err
}

// SaveSettings saves the settings as JSON to the path.
// It will remove any private keys.
func SaveSettings(path string, s Settings) error {
	b, err := s.asJson()
	if err != nil {
		return err
	}
	err = os.WriteFile(path, b, 0644)
	return err
}

func WriteJson(s Settings) ([]byte, error) {
	return s.asJson()
}

// String answers the string value at the given path.
func (s Settings) String(path string) (string, bool) {
	return getType(s, path, leafString)
}

// MustString answers the string value at the given path or
// fallback if path is absent.
func (s Settings) MustString(path, fallback string) string {
	if str, ok := s.String(path); ok {
		return str
	}
	return fallback
}

// Strings answers the string slice at the given path.
// If path points to a string slice, it's returned,
// otherwise the current keys are returned.
func (s Settings) Strings(path string) []string {
	p := strings.Split(path, pathSeparator)
	p = slices.DeleteFunc(p, func(n string) bool {
		return len(n) < 1
	})
	sub := Settings{}
	switch len(p) {
	case 0:
		return s.treeKeys()
	case 1:
		sub = s.lockedSubset(p[0])
	default:
		newPath := strings.Join(p[0:len(p)-1], pathSeparator)
		sub = s.lockedSubset(newPath)
		//		return s.lockedSubset(newPath).flatString(p[len(p)-1])
	}
	if len(sub.sliceKey) < 1 {
		return sub.treeKeys()
	}
	if list, ok := sub.t[sub.sliceKey].([]any); ok {
		sl := make([]string, 0, len(list))
		for _, _item := range list {
			switch item := _item.(type) {
			case string:
				sl = append(sl, item)
			default:
				sl = append(sl, fmt.Sprintf("%s", item))
			}
		}
		return sl
	}
	return nil
}

func (s Settings) treeKeys() []string {
	if len(s.t) > 0 {
		keys := make([]string, 0, len(s.t))
		for k := range s.t {
			keys = append(keys, k)
		}
		return keys
	}
	return nil
}

// Bool answers the bool value at the given path. The value
// must be a bool, a string (with value "true" or "t") or an
// element of a slice (for example, if the Settings contains
// "fruits": ["apple", "orange"], then Bool("fruits/apple")
// will return true).
func (s Settings) Bool(path string) (bool, bool) {
	return getType(s, path, leafBool)
}

// MustBool answers the bool value at the given path or
// fallback if path is absent.
func (s Settings) MustBool(path string, fallback bool) bool {
	if b, ok := s.Bool(path); ok {
		return b
	}
	return fallback
}

// Float64 answers the float64 value at the given path.
func (s Settings) Float64(path string) (float64, bool) {
	return getType(s, path, leafFloat64)
}

// MustFloat64 answers the bool value at the given path or
// fallback if path is absent.
func (s Settings) MustFloat64(path string, fallback float64) float64 {
	if b, ok := s.Float64(path); ok {
		return b
	}
	return fallback
}

// Int64 answers the int64 value at the given path.
func (s Settings) Int64(path string) (int64, bool) {
	return getType(s, path, leafInt64)
}

// MustInt64 answers the bool value at the given path or
// fallback if path is absent.
func (s Settings) MustInt64(path string, fallback int64) int64 {
	if b, ok := s.Int64(path); ok {
		return b
	}
	return fallback
}

func (s Settings) flatBoolList(p string) (bool, bool) {
	if list, ok := s.t[s.sliceKey].([]any); ok {
		for _, _item := range list {
			switch item := _item.(type) {
			case string:
				if item == p {
					return true, true
				}
			}
		}
	}
	return false, false
}

// Subset answers a subset of the settings tree based on
// walking down the path. The path can have components
// separated with "/".
func (s Settings) Subset(path string) Settings {
	return s.lockedSubset(path)
}

// lockedSubset answers a subset of the settings tree based on
// walking down the path. The path can have components
// separated with "/".
func (s Settings) lockedSubset(path string) Settings {
	p := strings.Split(path, "/")
	if len(p) < 1 {
		return emptySettings()
	}
	t := s.t
	for i, n := range p {
		if sub, ok := t[n]; ok {
			if subv, ok2 := sub.(map[string]any); ok2 {
				// Recurse down the map.
				t = subv
			} else if st, ok3 := sub.([]any); ok3 {
				// Special case: We are indexing into a slice, and the
				// result of the index is a map.
				if index, ok4 := pathIndex(i+1, p); ok4 && index < len(st) {
					if stm, ok5 := st[index].(map[string]any); ok5 {
						return Settings{t: stm}
					}
				}
				// Slices are handled specially: The parent key
				// is maintained, and a new map with just the key
				// and the slice value is returned. The key is then
				// annotated in "key" so subsequent callers know
				// how to access the slice value.
				insert := map[string]any{n: sub}
				return Settings{t: insert, sliceKey: n}
			} else {
				// It would be nice to handle this case: Essentially,
				// a single value is being requested, with no children.
				// Seems handy, but it breaks everything and doesn't
				// make sense with the API.
				return emptySettings()
			}
		} else {
			return emptySettings()
		}
	}
	return Settings{t: t}
}

// Length answers the length of the slice at path, or 0 if
// path is not a slice.
func (s Settings) Length(path string) int {
	if v, ok := getType(s, path, leafSlice); ok {
		return len(v)
	}
	return 0
}

// AllKeys iterates all keys at the top level.
func (s Settings) AllKeys() iter.Seq[string] {
	return func(yield func(string) bool) {
		for k, _ := range s.t {
			if !yield(k) {
				return
			}
		}
	}
}

func (s Settings) asJson() ([]byte, error) {
	b, err := json.Marshal(s.t)
	return b, err
}

func (s Settings) Print() {
	if d, err := json.MarshalIndent(s.t, "", "  "); err == nil {
		fmt.Println(string(d))
	}
}

type getFlatTypeFunc[T any] func(s Settings, path string) (T, bool)

func getType[T any](s Settings, path string, getFn getFlatTypeFunc[T]) (T, bool) {
	p := strings.Split(path, pathSeparator)
	switch len(p) {
	case 0:
		var t T
		return t, false
	case 1:
		return getFn(s, p[0])
	default:
		newPath := strings.Join(p[0:len(p)-1], pathSeparator)
		return getFn(s.lockedSubset(newPath), p[len(p)-1])
	}
}

// leafBool takes a path with no seperator, i.e.
// assumes it is an index in this map and not a subset,
// and returns the value.
func leafBool(s Settings, p string) (bool, bool) {
	// Lists are a special case
	if s.sliceKey != "" {
		return s.flatBoolList(p)
	}
	if v1, ok := s.t[p]; ok {
		switch v2 := v1.(type) {
		//		case int:
		//			fmt.Printf("Twice %v is %v\n", v, v*2)
		case bool:
			return v2, true
		case string:
			lc := strings.ToLower(v2)
			if lc == "true" || lc == "t" {
				return true, true
			}
			return false, false
		default:
			return false, false
		}
	}
	return false, false
}

// leafFloat64 takes a path with no seperator, i.e.
// assumes it is an index in this map and not a subset,
// and returns the value.
func leafFloat64(s Settings, p string) (float64, bool) {
	// Lists are a special case
	if s.sliceKey != "" {
		// Floats don't have slice support
		return 0.0, false
	}
	if v1, ok := s.t[p]; ok {
		switch v2 := v1.(type) {
		case int:
			return float64(v2), true
		case int64:
			return float64(v2), true
		case float32:
			return float64(v2), true
		case float64:
			return v2, true
		default:
			return 0.0, false
		}
	}
	return 0.0, false
}

// leafInt64 takes a path with no seperator, i.e.
// assumes it is an index in this map and not a subset,
// and returns the value.
func leafInt64(s Settings, p string) (int64, bool) {
	// Lists are a special case
	if s.sliceKey != "" {
		// Ints don't have slice support
		return 0, false
	}
	if v1, ok := s.t[p]; ok {
		switch v2 := v1.(type) {
		case int:
			return int64(v2), true
		case int64:
			return v2, true
		case float32:
			return int64(v2), true
		case float64:
			return int64(v2), true
		default:
			return 0, false
		}
	}
	return 0, false
}

// leafSlice takes a path with no seperator, i.e.
// assumes it is an index in this map and not a subset,
// and returns the value.
func leafSlice(s Settings, p string) ([]any, bool) {
	if v1, ok := s.t[p]; ok {
		switch v2 := v1.(type) {
		case []any:
			return v2, true
		default:
			fmt.Printf("returning default on type %t\n", v2)
			return nil, false
		}
	}
	return nil, false
}

// leafString takes a path with no seperator, i.e.
// assumes it is an index in this map and not a subset,
// and returns the value.
func leafString(s Settings, p string) (string, bool) {
	if v1, ok := s.t[p]; ok {
		switch v2 := v1.(type) {
		case string:
			return v2, true
		default:
			return "", false
		}
	}
	return "", false
}

// pathIndex looks at an index in a path slice and returns it
// as an int, if it converts.
func pathIndex(index int, path []string) (int, bool) {
	if index >= len(path) {
		return 0, false
	}
	if i, err := strconv.Atoi(path[index]); err == nil {
		return i, true
	}
	return 0, false
}

func NewEmptySettings() Settings {
	return emptySettings()
}

func emptySettings() Settings {
	return Settings{t: make(map[string]any)}
}

const (
	pathSeparator = `/`
)
