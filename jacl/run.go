package jacl

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/scanner"
	"unicode"

	"github.com/hackborn/onefunc/errors"
)

// Run compares a list of terms against a target. Target can be anything.
// Terms are described by the grammar
// {term} = {path}{comparison_operator}{value}
// Where {path} is a "/" separated list of identifiers
// {comparison_operator} is "="
// {value} is a string
// {path} identifiers can be either an integer to index slices, or
// a string for named fields.
// Example term:
// "0/Name=Ireland"
// where the target is a slice of structs that have a Name field.
func Run(target any, terms ...string) error {
	return RunOpts(Opts{}, target, terms...)
}

// RunOpts compares a list of terms against a target. Target can be anything.
// See Run docs for a decription of target and terms.
// Opts adds some configuration options, see Opts docs for a description.
func RunOpts(opts Opts, target any, terms ...string) error {
	for _, term := range terms {
		r := &runner{opts: opts, target: target}
		err := r.runTerm(term)
		if err != nil {
			return err
		}
	}
	return nil
}

type runner struct {
	first  errors.FirstBlock
	opts   Opts
	target any
}

func (r *runner) runTerm(term string) error {
	var scan scanner.Scanner
	scan.Init(strings.NewReader(term))
	scan.Whitespace = 0
	scan.Mode = scanner.ScanChars | scanner.ScanComments | scanner.ScanFloats | scanner.ScanIdents | scanner.ScanInts | scanner.ScanRawStrings | scanner.ScanStrings
	scan.IsIdentRune = r.isIdentRune
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
			return fmt.Errorf("expr \"%v\" contains tokens past the comparison (%v)", term, text)
		}
		if stage == runCompare {
			r.first.AddError(r.handleCompare(tok, text))
			stage = finishedCompare
			continue
		}
		if text == "=" {
			if stage != noCompare {
				return fmt.Errorf("expr \"%v\" has multiple comparisons", term)
			}
			stage = runCompare
			continue
		}

		r.first.AddError(r.handlePath(tok, text))
	}
	return r.first.Err
}

func (r *runner) isIdentRune(ch rune, i int) bool {
	// This is the standard text scanner ident rune, plus "{" and "}"
	// for keywords.
	ident := ch == '_' || ch == '{' || ch == '}' || unicode.IsLetter(ch) || (unicode.IsDigit(ch) && i > 0)
	return ident
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
	cmp, ok := r.target.(string)
	if !ok {
		return fmt.Errorf("Can't compare %v with %v", r.target, s)
	}
	if cmp != r.opts.processValue(s) {
		return fmt.Errorf("Have value \"%v\" but want \"%v\"", cmp, s)
	}
	return nil
}

func (r *runner) handlePath(tok rune, t string) error {
	switch tok {
	case scanner.Float:
		return fmt.Errorf("Can't navigate to float \"%v\"", t)
	case scanner.Int:
		i, err := strconv.Atoi(t)
		if err != nil {
			return err
		}
		return r.handlePathInt(i)
	case scanner.String:
		t = strings.Trim(t, `"`)
		return r.handlePathString(t)
	case scanner.Ident:
		return r.handlePathString(t)
	default:
		if t == "/" {
			// Path separator, continue
			return nil
		}
		return fmt.Errorf("Can't navigate to \"%v\"", t)
	}
}

func (r *runner) handlePathInt(i int) error {
	targetValue := reflect.ValueOf(r.target)
	switch targetValue.Kind() {
	case reflect.Slice:
		return r.handlePathIntOnSlice(i)
	default:
		return fmt.Errorf("Can't navigate to int \"%v\" on kind %v", i, targetValue.Kind())
	}
}

func (r *runner) handlePathIntOnSlice(i int) error {
	// We know r.dst is Kind slice
	sliceValue := reflect.ValueOf(r.target)
	if i >= sliceValue.Len() {
		return fmt.Errorf("Index %v is out of range on slice with len %v", i, sliceValue.Len())
	}
	v := sliceValue.Index(i)
	r.target = v.Interface()
	return nil
}

func (r *runner) handlePathString(s string) error {
	// Intercept keywords
	if s == keywordType {
		r.target = getTypeName(r.target)
		return nil
	}

	targetValue := reflect.ValueOf(r.target)
	switch targetValue.Kind() {
	case reflect.Struct:
		return r.handlePathStringOnStruct(s)
	case reflect.Ptr:
		elem := targetValue.Elem()
		r.target = elem.Interface()
		return r.handlePathString(s)
	case reflect.Map:
		return r.handlePathStringOnMap(s)
	default:
		return fmt.Errorf("Can't navigate to string \"%v\" on kind %v", s, targetValue.Kind())
	}
}

func (r *runner) handlePathStringOnStruct(fieldName string) error {
	// We know r.target is Kind struct
	structValue := reflect.ValueOf(r.target)
	field := structValue.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("no field for %v on struct %v", fieldName, r.target)
	}
	r.target = field.Interface()
	return nil
}

func (r *runner) handlePathStringOnMap(fieldName string) error {
	// We know r.target is Kind map
	mapValue := reflect.ValueOf(r.target)
	field := mapValue.MapIndex(reflect.ValueOf(fieldName))
	if !field.IsValid() {
		return fmt.Errorf("no field for %v on map %v", fieldName, r.target)
	}
	r.target = field.Interface()
	return nil
}

// ---------------------------------------------------------
// SUPPORT

// getTypeName answers the type of a, without the package name.
func getTypeName(a any) string {
	t := reflect.TypeOf(a)
	switch t.Kind() {
	case reflect.Ptr:
		return "*" + t.Elem().Name()
	default:
		return t.Name()
	}
}

// ---------------------------------------------------------
// CONST and VAR

type compareStage int

const (
	noCompare       compareStage = iota // No comparison token has been hit
	runCompare                          // The comparison token was processed; next token is the value
	finishedCompare                     // The comparison is finished
)

const (
	keywordType = `{type}`
)
