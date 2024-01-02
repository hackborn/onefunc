package strings

import (
	"fmt"

	"github.com/hackborn/onefunc/errors"
)

type CompileArgs struct {
	Quote     string
	Separator string
	Eb        oferrors.Block
}

// Compile all values into a single comma-separated string.
// Any string values are quoted with quote.
func Compile(args CompileArgs, singles ...any) string {
	w := GetWriter(args.Eb)
	defer PutWriter(w)

	for i, s := range singles {
		if i > 0 {
			w.WriteString(args.Separator)
		}
		switch v := s.(type) {
		case string:
			w.WriteString(args.Quote)
			w.WriteString(v)
			w.WriteString(args.Quote)
		default:
			w.WriteString(fmt.Sprintf("%v", v))
		}
	}
	return String(w)
}

// CompileStrings compiles all values into a single comma-separated string.
// Any string values are quoted with quote.
func CompileStrings(args CompileArgs, singles ...string) string {
	w := GetWriter(args.Eb)
	defer PutWriter(w)

	for i, s := range singles {
		if i > 0 {
			w.WriteString(args.Separator)
		}
		w.WriteString(args.Quote)
		w.WriteString(s)
		w.WriteString(args.Quote)
	}
	return String(w)
}
