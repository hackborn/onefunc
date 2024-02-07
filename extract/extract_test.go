package extract

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// ---------------------------------------------------------
// TEST-FROM
func TestFrom(t *testing.T) {
	table := []struct {
		src     any
		handler Handler
		want    string
	}{
		{Data1{A: "n"}, newHandler1(), `{"A":"n"}`},
		{Data1{A: "n"}, NewChain(filterMap1, newHandler1()), `{"name":"n"}`},
		{Data2{A: "n", B: 10}, newHandler1(), `{"A":"n","B":10}`},
		{Data2{A: "n", B: 10}, NewChain(filterMap1, newHandler1()), `{"name":"n"}`},
		{Data2{A: "n", B: 10}, NewChain(FilterMapOpts{F: filterMap1, Passthrough: true}, newHandler1()), `{"B":10,"name":"n"}`},
		{Data1{A: "n"}, &pairHandler{}, `["A"],["n"]`},
		{Data1{A: "n"}, NewChain(filterMap1, &pairHandler{}), `["name"],["n"]`},
		{Data3{A: "n"}, newHandler1(), `{"A":"n","B":0}`},
	}
	for i, v := range table {
		From(v.src, v.handler)
		have := v.handler.(Flattener).Flatten()

		if have != v.want {
			t.Fatalf("TestValues %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-GET-LAST
func TestGetLast(t *testing.T) {
	type ConvFunc func(Handler) (string, bool)
	convHandler1 := func(h Handler) (string, bool) {
		if t, ok := getLast[*Handler1](h); ok {
			return reflect.TypeOf(t).Elem().Name(), true
		}
		return "", false
	}
	convSlicer := func(h Handler) (string, bool) {
		if t, ok := getLast[Slicer](h); ok {
			return reflect.TypeOf(t).Elem().Name(), true
		}
		return "", false
	}

	table := []struct {
		h        Handler
		fn       ConvFunc
		wantName string // a type name
		wantOk   bool
	}{
		{nil, convHandler1, "", false},
		{&Handler1{}, convHandler1, "Handler1", true},
		{&sliceHandler{}, convSlicer, "sliceHandler", true},
		{NewChain(&Handler1{}, &sliceHandler{}), convSlicer, "sliceHandler", true},
	}
	for i, v := range table {
		have, haveOk := v.fn(v.h)

		if haveOk != v.wantOk {
			t.Fatalf("TestGetLast %v has ok \"%v\" but wanted \"%v\"", i, haveOk, v.wantOk)
		} else if have != v.wantName {
			t.Fatalf("TestGetLast %v has \"%v\" but wanted \"%v\"", i, have, v.wantName)
		}
	}
}

// ---------------------------------------------------------
// TEST-AS-MAP
func TestAsMap(t *testing.T) {
	table := []struct {
		src     any
		handler Handler
		opts    *MapOpts
		want    string
	}{
		{Data1{A: "n"}, newHandler1(), nil, `{"A":"n"}`},
		{Data1{A: "n"}, NewChain(filterMap1, newHandler1()), nil, `{"name":"n"}`},
		{Data1{A: "n"}, NewChain(filterMap1, newHandler1()), &MapOpts{}, `{"name":"n"}`},
	}
	for i, v := range table {
		haveMap := AsMap(v.src, v.handler, v.opts)
		have := ""
		if len(haveMap) > 0 {
			have = flattenMap(haveMap)
		}

		if have != v.want {
			t.Fatalf("TestAsMap %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-AS-SLICE
func TestAsSlice(t *testing.T) {
	table := []struct {
		src     any
		handler Handler
		opts    *SliceOpts
		want    string
	}{
		{Data1{A: "n"}, newHandler1(), nil, `["A","n"]`},
		{Data1{A: "n"}, NewChain(filterMap1, newHandler1()), nil, `["name","n"]`},
		{Data1{A: "n"}, NewChain(filterMap1, newHandler1()), &SliceOpts{Assign: "="}, `["name","=","n"]`},
	}
	for i, v := range table {
		haveSlice := AsSlice(v.src, v.handler, v.opts)
		have := ""
		if len(haveSlice) > 0 {
			have = flattenAny(haveSlice)
		}

		if have != v.want {
			t.Fatalf("TestAsSlice %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// SUPPORT var and const

var (
	filterMap1 = map[string]string{
		"A": "name",
	}
)

// ---------------------------------------------------------
// SUPPORT interfaces

type Flattener interface {
	Flatten() string
}

// ---------------------------------------------------------
// SUPPORT types

type Tuple[A any, B any] struct {
	A A
	B B
}

type Data1 struct {
	A string
}

type Data2 struct {
	A string
	B int
}

type Data3 struct {
	A     string
	B     int
	_priv string
}

func newHandler1() Handler {
	results := make(map[string]any)
	return &Handler1{results: results}
}

type Handler1 struct {
	results map[string]any
}

func (h *Handler1) Handle(name string, value any) (string, any) {
	h.results[name] = value
	return name, value
}

func (h *Handler1) Flatten() string {
	return flattenMap(h.results)
}

type pairHandler struct {
	fields []any
	values []any
}

func (h *pairHandler) Handle(name string, value any) (string, any) {
	h.fields = append(h.fields, name)
	h.values = append(h.values, value)
	return name, value
}

func (h *pairHandler) Flatten() string {
	var sb strings.Builder
	b, err := json.Marshal(h.fields)
	panicErr(err)
	sb.WriteString(string(b))
	sb.WriteString(",")
	b, err = json.Marshal(h.values)
	panicErr(err)
	sb.WriteString(string(b))
	return sb.String()
}

// ---------------------------------------------------------
// HANDLER SUPPORT

func (c Chain) Flatten() string {
	// Assume the last element has the data of interest
	last := len(c) - 1
	if last >= 0 && c[last] != nil {
		return c[last].(Flattener).Flatten()
	}
	return ""
}

// ---------------------------------------------------------
// FLATTEN SUPPORT

func flattenAny(a any) string {
	b, err := json.Marshal(a)
	panicErr(err)
	return string(b)
}

func flattenMap(m map[string]any) string {
	var sb strings.Builder
	_, err := sb.WriteString("{")
	panicErr(err)
	for i, t := range orderMap(m) {
		if i > 0 {
			_, err = sb.WriteString(",")
			panicErr(err)
		}
		_, err = sb.WriteString(fmt.Sprintf("\"%v\":", t.A))
		panicErr(err)
		b, err := json.Marshal(t.B)
		panicErr(err)
		_, err = sb.WriteString(string(b))
		panicErr(err)
	}
	_, err = sb.WriteString("}")
	panicErr(err)
	return sb.String()
}

func orderMap(m map[string]any) []Tuple[string, any] {
	keys := keys(m)
	slices.Sort(keys)
	results := make([]Tuple[string, any], 0, len(keys))
	for _, s := range keys {
		v, _ := m[s]
		results = append(results, Tuple[string, any]{A: s, B: v})
	}
	return results
}

func keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
