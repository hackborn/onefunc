package cfg

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	oferrors "github.com/hackborn/onefunc/errors"
)

// Option used during construction of a Settings.
// Clients can add to the settings with the Builder.
type Option func(Builder, oferrors.Block)

type Process func(map[string]any) map[string]any

// WithFS loads all files that match the pattern
// into the Settings. All matched files must be in JSON format.
// See path.Match() for match rules.
func WithFS(fsys fs.FS, pattern string, processors ...Process) Option {
	return func(b Builder, eb oferrors.Block) {
		matches, err := fs.Glob(fsys, pattern)
		eb.AddError(err)
		for _, match := range matches {
			dat, err := fs.ReadFile(fsys, match)
			if err != nil {
				err = fmt.Errorf("%v: %w", match, err)
				eb.AddError(err)
			}

			s := b.NewSettings()
			err = json.Unmarshal(dat, &s)
			if err != nil {
				err = fmt.Errorf("%v: %w", match, err)
				eb.AddError(err)
			}
			for _, process := range processors {
				s = process(s)
			}
			b.AddSettings(s)
		}
	}
}

// WithEnv loads all matching env vars.
func WithEnv(match EnvMatcher) Option {
	return func(b Builder, eb oferrors.Block) {
		envs := os.Environ()
		s := b.NewSettings()
		for _, env := range envs {
			pos := strings.Index(env, "=")
			if pos > 0 && pos < len(env)-1 {
				left := env[0:pos]
				right := env[pos+1:]
				left, err := match.Match(left)
				eb.AddError(err)
				if left != "" {
					s[left] = right
				}
			}
		}
		b.AddSettings(s)
	}
}

// EnvMatcher is used by WithEnv to accept env vars and
// potentially modify the name.
type EnvMatcher interface {
	// Answer a key if a match, or else "". The key
	// might be the same key in the args, or it might be modified.
	Match(key string) (string, error)
}

// EnvPattern loads all env vars that match the pattern
// into the Settings. See path.Match() for match rules.
func EnvPattern(pattern string) EnvMatcher {
	return &envPatternMatcher{pattern: pattern}
}

// EnvPrefix loads all env vars that match the prefix
// into the Settings. The prefix will be stripped.
func EnvPrefix(prefix string) EnvMatcher {
	return &envPrefixMatcher{prefix: prefix}
}

type envPatternMatcher struct {
	pattern string
}

func (m *envPatternMatcher) Match(key string) (string, error) {
	ok, err := path.Match(m.pattern, key)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", nil
	}
	return key, nil
}

type envPrefixMatcher struct {
	prefix string
}

func (m *envPrefixMatcher) Match(key string) (string, error) {
	if !strings.HasPrefix(key, m.prefix) {
		return "", nil
	}
	return key[len(m.prefix):], nil
}

// WithKeys adds all keys in the supplied settings.
func WithKeys(src Settings, keys []string) Option {
	return func(b Builder, eb oferrors.Block) {
		s := b.NewSettings()
		for _, key := range keys {
			if v, ok := src.t[key]; ok {
				s[key] = v
			}
		}
		b.AddSettings(s)
	}
}

// WithMap adds all keys in the supplied map.
// Note this is only a shallow copy of the top layer,
// anything below that is the original reference.
func WithMap(m map[string]any) Option {
	return func(b Builder, eb oferrors.Block) {
		s := b.NewSettings()
		for k, v := range m {
			s[k] = v
		}
		b.AddSettings(s)
	}
}

// WithSettings acts as a deep copy on src.
func WithSettings(src Settings) Option {
	return func(b Builder, eb oferrors.Block) {
		t := b.NewSettings()
		dat, err := src.asJson()
		err = cmp.Or(err, json.Unmarshal(dat, &t))
		if err != nil {
			eb.AddError(err)
		} else {
			removePrivateKeys(t)
			b.AddSettings(t)
		}
	}
}
