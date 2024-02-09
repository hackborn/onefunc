package cfg

import (
	"encoding/json"
	"io/fs"

	oferrors "github.com/hackborn/onefunc/errors"
)

type Option func(Settings, oferrors.Block)

// WithFS loads all files that match the glob pattern
// into the Settings. All matched files must be in JSON format.
// See fs.Glob() for glob rules.
func WithFS(fsys fs.FS, glob string) Option {
	return func(s Settings, eb oferrors.Block) {
		matches, err := fs.Glob(fsys, glob)
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
