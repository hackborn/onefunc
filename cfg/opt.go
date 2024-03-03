package cfg

import (
	"encoding/json"
	"io/fs"
	"os"
	"path"
	"strings"

	oferrors "github.com/hackborn/onefunc/errors"
)

type Option func(Settings, oferrors.Block)

// WithFS loads all files that match the pattern
// into the Settings. All matched files must be in JSON format.
// See path.Match() for match rules.
func WithFS(fsys fs.FS, pattern string) Option {
	return func(s Settings, eb oferrors.Block) {
		matches, err := fs.Glob(fsys, pattern)
		eb.AddError(err)
		for _, match := range matches {
			dat, err := fs.ReadFile(fsys, match)
			eb.AddError(err)

			s2 := emptySettings(s.rw)
			err = json.Unmarshal(dat, &s2.t)
			eb.AddError(err)

			mergeKeys(s.t, s2.t)
		}
	}
}

// WithEnv loads all matching env vars.
func WithEnv(match EnvMatcher) Option {
	return func(s Settings, eb oferrors.Block) {
		envs := os.Environ()
		s2 := emptySettings(s.rw)
		for _, env := range envs {
			pos := strings.Index(env, "=")
			if pos > 0 && pos < len(env)-1 {
				left := env[0:pos]
				right := env[pos+1:]
				left, err := match.Match(left)
				eb.AddError(err)
				if left != "" {
					s2.t[left] = right
				}
			}
		}
		mergeKeys(s.t, s2.t)
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
