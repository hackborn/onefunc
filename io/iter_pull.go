package io

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Puller interface {
	// Answer the next item. IterEntry.Entry will be nil
	// when there are no more items.
	Pull() (IterEntry, error)
}

func PullAll(root string, options ...IterOption) (Puller, error) {
	info, err := os.Lstat(root)
	if err != nil {
		return nil, err
	}
	entry := fs.FileInfoToDirEntry(info)
	parent := strings.TrimSuffix(root, entry.Name())
	entries := []fs.DirEntry{entry}
	steps := []*_pullStep{&_pullStep{parent: parent, entries: entries}}
	return &_puller{options: options, steps: steps}, nil
}

// _puller walks a directory tree, given clients access
// with a pull pattern.
type _puller struct {
	options []IterOption
	steps   []*_pullStep
	err     error
}

func (p *_puller) Pull() (IterEntry, error) {
	if p.err != nil {
		return IterEntry{}, p.err
	}
	ie, valid, err := p.next()
	for !valid {
		ie, valid, err = p.next()
	}
	p.err = err
	return ie, err
}

func (p *_puller) next() (IterEntry, bool, error) {
	if pd := p.tail(); pd == nil {
		return IterEntry{}, true, nil
	} else {
		if pd.pos >= len(pd.entries) {
			p.steps = p.steps[:len(p.steps)-1]
			return IterEntry{}, false, nil
		}
		entry := pd.entries[pd.pos]
		ie := IterEntry{Path: filepath.Join(pd.parent, entry.Name()),
			Entry: entry}
		// Run the options
		skipHandling := false
		if err := p.runOptions(ie); err != nil {
			switch err {
			case filepath.SkipAll:
				// Cancel
				p.steps = nil
				return IterEntry{}, true, nil
			case filepath.SkipDir:
				// Advance
				pd.pos++
				return IterEntry{}, false, nil
			case SkipFile:
				// This is very counterintuitive, but All() handles
				// SkipFile on directories by not passing on the directory
				// to the client, but still diving into it.
				if !ie.Entry.IsDir() {
					// Advance
					pd.pos++
					return IterEntry{}, false, nil
				} else {
					skipHandling = true
				}
			default:
				return IterEntry{}, true, err
			}
		}
		// Always advance the current item
		pd.pos++
		if ie.Entry.IsDir() {
			// Process directories.
			if dirs, err := os.ReadDir(ie.Path); err != nil {
				return IterEntry{}, true, err
			} else {
				step := &_pullStep{parent: ie.Path, entries: dirs}
				p.steps = append(p.steps, step)
				// With my children in place, let clients process me.
				if skipHandling {
					return IterEntry{}, false, nil
				} else {
					return ie, true, nil
				}
			}
		} else {
			return ie, true, nil
		}
	}
}

func (p *_puller) runOptions(entry IterEntry) error {
	for _, opt := range p.options {
		if err := opt(entry); err != nil {
			return err
		}
	}
	return nil
}

func (p *_puller) tail() *_pullStep {
	if len(p.steps) < 1 {
		return nil
	}
	return p.steps[len(p.steps)-1]
}

type _pullStep struct {
	parent  string
	entries []fs.DirEntry
	pos     int
}
