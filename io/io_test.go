package io

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	oferrors "github.com/hackborn/onefunc/errors"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

// ---------------------------------------------------------
// TEST-ALL

func TestAll(t *testing.T) {
	f := func(path string, args ...any) {
		t.Helper()

		opts := []IterOption{}
		wants := []string{}
		for _, _arg := range args {
			switch arg := _arg.(type) {
			case IterOption:
				opts = append(opts, arg)
			case string:
				wants = append(wants, arg)
			default:
				oferrors.Panic(fmt.Errorf("Unknown argument %T %v", _arg, _arg))
			}
		}

		root := filepath.Join("testdata", path)
		wantIndex := 0
		for e, err := range All(root, opts...) {
			oferrors.Panic(err)
			p := strings.TrimPrefix(e.Path, root)
			// Due to the nature of testing I can just skip the root.
			if p == "" {
				continue
			}
			if wantIndex < len(wants) {
				if wants[wantIndex] != p {
					t.Fatalf("want \"%v\" at index %v but have \"%v\"", wants[wantIndex], wantIndex, p)
				}
				wantIndex++
			}
		}
		if len(wants) > wantIndex {
			t.Fatalf("want \"%v\" but nothing left", wants[wantIndex])
		}
	}

	f("iter", "/b.txt", "/c.txt", "/dir", "/dir/a.txt")
	f("iter", SkipDir("testdata/iter/dir"), "/b.txt", "/c.txt")
	f("iter", skipFileNamed("testdata/iter/b.txt"), "/c.txt", "/dir", "/dir/a.txt")
}

func skipFileNamed(n string) IterOption {
	return func(e IterEntry) error {
		if e.Path == n {
			return SkipFile
		}
		return nil
	}
}

// ---------------------------------------------------------
// TEST-EXPAND-PATH

func TestExpandPath(t *testing.T) {
	f := func(path, want string) {
		t.Helper()

		have, haveErr := ExpandPath(path)
		if haveErr != nil {
			t.Fatalf("TestExpandPath has err %v", haveErr)
		} else if have != want {
			t.Fatalf("TestExpandPath wants %v but has %v", want, have)
		}
	}

	f("termy", "termy")
}

// ---------------------------------------------------------
// TEST-XML-ALL-TOKENS

func TestXmlAllTokens(t *testing.T) {
	f := func(filename string, cmp XmlComparison) {
		t.Helper()

		filename = filepath.Join("testdata", fmt.Sprintf("%s.xml", filename))
		for tok, err := range XmlAllTokens(FileResource(filename)) {
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
				return
			}
			if err = cmp.Token(tok); err != nil {
				t.Fatalf("%v", err)
			}
		}
	}

	f("a", chardatas("hi", "bye"))
}

// ---------------------------------------------------------
// COMPARISON

type XmlComparison interface {
	Token(xml.Token) error
}

func chardatas(texts ...string) XmlComparison {
	return &chardataList{text: texts}
}

type chardataList struct {
	text []string
	pos  int
}

func (c *chardataList) Token(tok xml.Token) error {
	// It's OK for there to be more tokens then we want to check, just ignore.
	if c.pos >= len(c.text) {
		return nil
	}
	if cd, ok := tok.(xml.CharData); ok {
		have, want := string(cd), c.text[c.pos]
		if have != want {
			return fmt.Errorf("pos %v has \"%s\" but wants \"%s\"", c.pos, have, want)
		}
		c.pos++
	}
	return nil
}
