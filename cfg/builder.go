package cfg

import (
	"fmt"
	"reflect"
	"strings"
)

// Builder is used to add new settings during a Settings construction.
type Builder interface {
	// Settings answers the raw settings. Generally treat this
	// as read-only, and use AddSettings() to make any changes.
	Settings() map[string]any

	// NewSettings answers a new empty settings.
	NewSettings() map[string]any

	// AddSettings adds the new settings to the settings being built.
	AddSettings(map[string]any)
}

type _builder struct {
	t tree
}

func (b *_builder) Settings() map[string]any {
	return b.t
}

func (b *_builder) NewSettings() map[string]any {
	return make(map[string]any)
}

func (b *_builder) AddSettings(m map[string]any) {
	if len(m) > 0 {
		mergeKeys(b.t, m)
	}
}

// Given two maps, recursively merge right into left. Adapted from
// https://stackoverflow.com/questions/22621754/how-can-i-merge-two-maps-in-go
func mergeKeys(left, right tree) tree {
	for key, rightVal := range right {
		if leftVal, present := left[key]; present {
			// Key exists. If my new value is a map, recurse.
			// If it's not, replace.
			if reflect.TypeOf(rightVal).Kind() == reflect.Map {
				left[key] = mergeKeys(leftVal.(tree), rightVal.(tree))
			} else {
				left[key] = rightVal
			}
		} else {
			// key not in left so we can just shove it in
			left[key] = rightVal
		}
	}
	return left
}

// Setter is a simplified version of Builder. Builder is
// designed to give raw access to the settings, but idelly
// that is unnecessary. However I'm sure there are things
// relying on it by this point.
type Setter struct {
	all tree
}

func (s Setter) Settings() Settings {
	return Settings{t: s.all}
}

func (s Setter) Set(path string, value any) error {
	// The framework should guarantee this is never true.
	addTo := s.all
	for part := range strings.SplitSeq(path, pathSeparator) {
		if addTo == nil {
			return fmt.Errorf("Setter.Set is missing settings")
		}
		if strings.HasSuffix(path, part) {
			addTo[part] = value
		} else {
			if v := addTo[part]; v == nil {
				newAddTo := make(map[string]any)
				addTo[part] = newAddTo
				addTo = newAddTo
			} else if found, ok := v.(map[string]any); ok {
				addTo = found
			} else {
				return fmt.Errorf("Trying to replace item %s (%T) with map", part, v)
			}
		}
	}
	// Should this return error if can't add but nothing identifiably went wrong?
	return nil
}
