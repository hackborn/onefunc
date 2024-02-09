package cfg

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hackborn/onefunc/lock"
)

type tree = map[string]any

type Settings struct {
	rw *sync.RWMutex
	t  tree
}

func LoadSettings(fsys fs.FS, fn string) (Settings, error) {
	s := emptySettings(&sync.RWMutex{})
	dat, err := readFileFS(fsys, fn)
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(dat, &s.t)
	return s, err
}

func LoadSettingsFolder(fsys fs.FS, fn string) (Settings, error) {
	s := emptySettings(&sync.RWMutex{})
	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		ext := strings.ToLower(filepath.Ext(p))
		if d.IsDir() || ext != ".json" {
			return nil
		}
		dat, err := readFileFS(fsys, p)
		if err != nil {
			return err
		}
		s2 := emptySettings(s.rw)
		err = json.Unmarshal(dat, &s2.t)
		mergeKeys(s.t, s2.t)
		return err
	})
	return s, err
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

// String answers the string value at the given path.
func (s Settings) String(path string) (string, bool) {
	defer lock.Read(s.rw).Unlock()
	p := strings.Split(path, pathSeparator)
	switch len(p) {
	case 0:
		return "", false
	case 1:
		return s.flatString(p[0])
	default:
		newPath := strings.Join(p[0:len(p)-1], pathSeparator)
		return s.lockedSubset(newPath).flatString(p[len(p)-1])
	}
}

// MustString answers the string value at the given path or
// fallback if path is absent.
func (s Settings) MustString(path, fallback string) string {
	if str, ok := s.String(path); ok {
		return str
	}
	return fallback
}

// Bool answers the bool value at the given path.
func (s Settings) Bool(path string) (bool, bool) {
	defer lock.Read(s.rw).Unlock()
	p := strings.Split(path, pathSeparator)
	switch len(p) {
	case 0:
		return false, false
	case 1:
		return s.flatBool(p[0])
	default:
		newPath := strings.Join(p[0:len(p)-1], pathSeparator)
		return s.lockedSubset(newPath).flatBool(p[len(p)-1])
	}
}

// MustBool answers the bool value at the given path or
// fallback if path is absent.
func (s Settings) MustBool(path string, fallback bool) bool {
	if b, ok := s.Bool(path); ok {
		return b
	}
	return fallback
}

// flatBool takes a path with no seperator, i.e.
// assumes it is an index in this map and not a subset,
// and returns the value.
func (s Settings) flatBool(p string) (bool, bool) {
	if v1, ok := s.t[p]; ok {
		switch v2 := v1.(type) {
		//		case int:
		//			fmt.Printf("Twice %v is %v\n", v, v*2)
		case bool:
			return v2, true
		default:
			return false, false
		}
	}
	return false, false
}

// latString takes a path with no seperator, i.e.
// assumes it is an index in this map and not a subset,
// and returns the value.
func (s Settings) flatString(p string) (string, bool) {
	if v1, ok := s.t[p]; ok {
		switch v2 := v1.(type) {
		//		case int:
		//			fmt.Printf("Twice %v is %v\n", v, v*2)
		case string:
			return v2, true
		default:
			return "", false
		}
	}
	return "", false
}

// Subset answers a subset of the settings tree based on
// walking down the path. The path can have components
// separated with "/".
func (s Settings) Subset(path string) Settings {
	defer lock.Read(s.rw).Unlock()
	return s.lockedSubset(path)
}

// lockedSubset answers a subset of the settings tree based on
// walking down the path. The path can have components
// separated with "/".
func (s Settings) lockedSubset(path string) Settings {
	p := strings.Split(path, "/")
	if len(p) < 1 {
		return emptySettings(s.rw)
	}
	t := s.t
	for _, n := range p {
		if sub, ok := t[n]; ok {
			if subv, ok2 := sub.(map[string]any); ok2 {
				t = subv
			} else {
				// It would be nice to handle this case: Essentially,
				// a single value is being requested, with no children.
				// Seems handy, but it breaks everything and doesn't
				// make sense with the API.
				return emptySettings(s.rw)
			}
		} else {
			return emptySettings(s.rw)
		}
	}
	return Settings{rw: s.rw, t: t}
}

// SetValue sets the given key to the given value.
// `value` can be nil, and an empty map will be created.
// Currently the key can not contain a path element; if you
// want to se a path, find the subset first.
func (s Settings) SetValue(key string, value interface{}) error {
	if s.rw == nil || s.t == nil {
		return fmt.Errorf("invalid state")
	}
	if strings.Contains(key, pathSeparator) {
		// If this become annoying I will build in finding the subset.
		return fmt.Errorf(`key can not contain path character "` + pathSeparator + `"`)
	}
	if value == nil {
		value = make(map[string]any)
	}

	defer lock.Write(s.rw).Unlock()
	s.t[key] = value
	s.t[changedKey] = true
	return nil
}

func (s Settings) IsChanged() bool {
	v, _ := s.Bool(changedKey)
	return v
}

func (s Settings) WalkKeys(fn WalkKeysFunc) error {
	for k, _ := range s.t {
		err := fn(k)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Settings) asJson() ([]byte, error) {
	defer lock.Read(s.rw).Unlock()
	s.lockedRemotePrivateKeys()
	b, err := json.Marshal(s.t)
	return b, err
}

func (s Settings) lockedRemotePrivateKeys() {
	for k, _ := range s.t {
		if strings.HasPrefix(k, privateKeyPrefix) {
			delete(s.t, k)
		}
	}
}

func (s Settings) Print() {
	fmt.Println(s.t)
}

type WalkKeysFunc func(key string) error

func emptySettings(rw *sync.RWMutex) Settings {
	return Settings{rw: rw, t: make(map[string]any)}
}

const (
	pathSeparator    = `/`
	privateKeyPrefix = `_$cfg_`
	changedKey       = privateKeyPrefix + `changed`
)
