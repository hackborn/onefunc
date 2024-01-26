package jacl

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/hackborn/onefunc/errors"
)

func Run(dst any, exprs ...string) error {
	for _, expr := range exprs {
		r := &runner{dst: dst}
		err := r.runExpr(expr)
		if err != nil {
			return err
		}
	}
	return nil
}

type runner struct {
	first errors.FirstBlock
	dst   any
}

func (r *runner) runExpr(expr string) error {
	var scan scanner.Scanner
	scan.Init(strings.NewReader(expr))
	scan.Whitespace = 0
	scan.Mode = scanner.ScanChars | scanner.ScanComments | scanner.ScanFloats | scanner.ScanIdents | scanner.ScanInts | scanner.ScanRawStrings | scanner.ScanStrings
	scan.Error = func(s *scanner.Scanner, msg string) {
		r.first.AddError(fmt.Errorf("run error: %v", msg))
	}
	stage := noCompare

	for tok := scan.Scan(); tok != scanner.EOF; tok = scan.Scan() {
		if r.first.Err != nil {
			return r.first.Err
		}
		// fmt.Println("TOK", tok, "name", scanner.TokenString(tok), "text", lexer.TokenText())

		// Always skip whitespace, whatever the stage
		switch tok {
		case ' ', '\r', '\t', '\n':
			continue
		}

		text := scan.TokenText()
		// If there are any tokens passed the comparison, error
		if stage == finishedCompare {
			return fmt.Errorf("expr \"%v\" contains tokens past the comparison (%v)", expr, text)
		}
		if stage == runCompare {
			r.first.AddError(r.handleCompare(tok, text))
			stage = finishedCompare
			continue
		}
		if text == "=" {
			if stage != noCompare {
				return fmt.Errorf("expr \"%v\" has multiple comparisons", expr)
			}
			stage = runCompare
			continue
		}

		r.first.AddError(r.handleNavigate(tok, text))
	}
	return r.first.Err
}

func (r *runner) handleCompare(tok rune, t string) error {
	switch tok {
	/*
		case scanner.Float:
			tt.tt = floatToken
		case scanner.Int:
			tt.tt = intToken
	*/
	case scanner.String:
		t = strings.Trim(t, `"`)
		return r.handleStringCompare(t)
	case scanner.Ident:
		return r.handleStringCompare(t)
	default:
		return r.handleStringCompare(t)
	}
}

func (r *runner) handleStringCompare(s string) error {
	cmp, ok := r.dst.(string)
	if !ok {
		return fmt.Errorf("Can't compare %v with %v", r.dst, s)
	}
	if cmp != s {
		return fmt.Errorf("Have value \"%v\" but want \"%v\"", cmp, s)
	}
	return nil
}

func (r *runner) handleNavigate(tok rune, t string) error {
	switch tok {
	case scanner.Float:
		return fmt.Errorf("Can't navigate to float \"%v\"", t)
	case scanner.Int:
		i, err := strconv.Atoi(t)
		if err != nil {
			return err
		}
		return r.handleNavigateInt(i)
	case scanner.String:
		t = strings.Trim(t, `"`)
		return r.handleNavigateString(t)
	case scanner.Ident:
		return r.handleNavigateString(t)
	default:
		if t == "/" {
			// Path separator, continue
			return nil
		}
		return fmt.Errorf("Can't navigate to \"%v\"", t)
	}
}

func (r *runner) handleNavigateInt(i int) error {
	dstValue := reflect.ValueOf(r.dst)
	switch dstValue.Kind() {
	case reflect.Slice:
		return r.handleNavigateIntOnSlice(i)
	default:
		return fmt.Errorf("Can't navigate to int \"%v\" on kind %v", i, dstValue.Kind())
	}
}

func (r *runner) handleNavigateIntOnSlice(i int) error {
	// We know r.dst is Kind slice
	sliceValue := reflect.ValueOf(r.dst)
	if i >= sliceValue.Len() {
		return fmt.Errorf("Index %v is out of range on slice with len %v", i, sliceValue.Len())
	}
	v := sliceValue.Index(i)
	r.dst = v.Interface()
	return nil
}

func (r *runner) handleNavigateString(s string) error {
	dstValue := reflect.ValueOf(r.dst)
	switch dstValue.Kind() {
	case reflect.Struct:
		return r.handleNavigateStringOnStruct(s)
	case reflect.Ptr:
		elem := dstValue.Elem()
		r.dst = elem.Interface()
		return r.handleNavigateString(s)
	default:
		return fmt.Errorf("Can't navigate to string \"%v\" on kind %v", s, dstValue.Kind())
	}
}

func (r *runner) handleNavigateStringOnStruct(fieldName string) error {
	// We know r.dst is Kind struct
	structValue := reflect.ValueOf(r.dst)
	field := structValue.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("no field for %v on %v", fieldName, r.dst)
	}
	r.dst = field.Interface()
	return nil
}

// ---------------------------------------------------------
// CONST and VAR

type compareStage int

const (
	noCompare       compareStage = iota // No comparison token has been hit
	runCompare                          // The comparison token was processed; next token is the value
	finishedCompare                     // The comparison is finished
)
