package strings

import (
	"iter"
	"strings"
	"unicode/utf8"
)

// Coming in 1.2.4, not quite here.
func SplitSeq(s, sep string) iter.Seq[string] {
	return splitSeq(s, sep, 0)
}

// splitSeq is SplitSeq or SplitAfterSeq, configured by how many
// bytes of sep to include in the results (none or all).
func splitSeq(s, sep string, sepSave int) iter.Seq[string] {
	if len(sep) == 0 {
		return explodeSeq(s)
	}
	return func(yield func(string) bool) {
		for {
			i := strings.Index(s, sep)
			if i < 0 {
				break
			}
			frag := s[:i+sepSave]
			if !yield(frag) {
				return
			}
			s = s[i+len(sep):]
		}
		yield(s)
	}
}

func explodeSeq(s string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for len(s) > 0 {
			_, size := utf8.DecodeRuneInString(s)
			if !yield(s[:size]) {
				return
			}
			s = s[size:]
		}
	}
}
