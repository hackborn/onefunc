package io

import (
	"io/fs"
	"iter"
	"path/filepath"
	"slices"
)

type IterEntry struct {
	Path  string
	Entry fs.DirEntry
}

// All iterates over a directory. It is a replacement for
// filepath.Walk in a cleaner API.
// options allows clients to filter the results.
//
// Not a lot of value, but I imagine this is how things
// would have been designed had iterators been there at the start.
func All(root string, options ...IterOption) iter.Seq2[IterEntry, error] {
	return func(yield func(IterEntry, error) bool) {
		filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			entry := IterEntry{Path: path, Entry: d}
			if err != nil {
				if !yield(entry, err) {
					return filepath.SkipAll
				}
				return nil
			}

			for _, opt := range options {
				if err = opt(entry); err != nil {
					if err == SkipFile {
						return nil
					}
					return err
				}
			}
			if !yield(entry, nil) {
				return filepath.SkipAll
			}
			return nil
		})
	}
}

// ---------------------------------------------------------
// ITER-OPTIONS

// IterOption is used to determine if the IterEntry is
// passed to the client. Return nil if you want it sent on,
// otherwise a Walk error (or fs.SkipDir, etc.) Also supports
// this package's SkipFile, to not pass this single entry
// to the client but otherwise continue iterating.
type IterOption func(IterEntry) error

// SkipDir will skip any matching directories.
func SkipDir(dirs ...string) IterOption {
	return func(e IterEntry) error {
		if e.Entry.IsDir() && slices.Contains(dirs, e.Path) {
			return filepath.SkipDir
		}
		return nil
	}
}
