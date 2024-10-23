package reflect

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	oferrors "github.com/hackborn/onefunc/errors"
	"github.com/hackborn/onefunc/jacl"
	"github.com/hackborn/onefunc/math/geo"
)

// ---------------------------------------------------------
// TEST-COPY
func TestCopy(t *testing.T) {
	f := func(dst, src any, wantErr error, want []string) {
		t.Helper()

		haveErr := Copy(dst, src)
		if err := jacl.RunErr(haveErr, wantErr); err != nil {
			t.Fatalf("Has err %v but wants %v", haveErr, wantErr)
		} else if err := jacl.Run(dst, want...); err != nil {
			dat, _ := json.Marshal(dst)
			t.Fatalf("Wants %v but has %v", want, string(dat))
		}
	}
	f(&Data1{}, Data1{A: "a"}, nil, []string{"A=a"})
	f(&Data2{}, Data1{A: "a"}, nil, []string{"A=a"})
	f(&Data1{}, Data2{A: "a", B: 10}, nil, []string{"A=a"})
	//
	// panic("n")
}

// ---------------------------------------------------------
// TEST-GET
func TestGet(t *testing.T) {
	table := []struct {
		src     any
		handler GetHandler
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
		Get(v.src, v.handler)
		have := v.handler.(Flattener).Flatten()

		if have != v.want {
			t.Fatalf("TestGet %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-GET-FLOAT-64
func TestGetFloat64(t *testing.T) {
	f := func(src any, wantErr error, want float64) {
		t.Helper()

		have, haveErr := GetFloat64(src)
		if err := jacl.RunErr(haveErr, wantErr); err != nil {
			t.Fatalf("Has err %v but wants %v", haveErr, wantErr)
		} else if geo.FloatsEqualTol(have, want, 0.0001) == false {
			t.Fatalf("Wants %v but has %v", want, have)
		}
	}
	f(10, nil, 10.)
	f(int(10), nil, 10.)
	f(int8(10), nil, 10.)
	f(int64(10), nil, 10.)
	f(uint8(10), nil, 10.)
	f(uint64(10), nil, 10.)
	f(10., nil, 10.)
	f(float32(10), nil, 10.)
	f(float64(10), nil, 10.)
}

// ---------------------------------------------------------
// TEST-GET-LAST
func TestGetLast(t *testing.T) {
	type ConvFunc func(GetHandler) (string, bool)
	convHandler1 := func(h GetHandler) (string, bool) {
		if t, ok := getLast[*Handler1](h); ok {
			return reflect.TypeOf(t).Elem().Name(), true
		}
		return "", false
	}
	convSlicer := func(h GetHandler) (string, bool) {
		if t, ok := getLast[Slicer](h); ok {
			return reflect.TypeOf(t).Elem().Name(), true
		}
		return "", false
	}

	table := []struct {
		h        GetHandler
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
// TEST-GET-AS-MAP
func TestGetAsMap(t *testing.T) {
	table := []struct {
		src     any
		handler GetHandler
		opts    *MapOpts
		want    string
	}{
		{Data1{A: "n"}, newHandler1(), nil, `{"A":"n"}`},
		{Data1{A: "n"}, NewChain(filterMap1, newHandler1()), nil, `{"name":"n"}`},
		{Data1{A: "n"}, NewChain(filterMap1, newHandler1()), &MapOpts{}, `{"name":"n"}`},
	}
	for i, v := range table {
		haveMap := GetAsMap(v.src, v.handler, v.opts)
		have := ""
		if len(haveMap) > 0 {
			have = flattenMap(haveMap)
		}

		if have != v.want {
			t.Fatalf("TestGetAsMap %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-GET-AS-SLICE
func TestGetAsSlice(t *testing.T) {
	table := []struct {
		src     any
		handler GetHandler
		opts    *SliceOpts
		want    string
	}{
		{Data1{A: "n"}, newHandler1(), nil, `["A","n"]`},
		{Data1{A: "n"}, NewChain(filterMap1, newHandler1()), nil, `["name","n"]`},
		{Data1{A: "n"}, NewChain(filterMap1, newHandler1()), &SliceOpts{Assign: "="}, `["name","=","n"]`},
	}
	for i, v := range table {
		haveSlice := GetAsSlice(v.src, v.handler, v.opts)
		have := ""
		if len(haveSlice) > 0 {
			have = flattenAny(haveSlice)
		}

		if have != v.want {
			t.Fatalf("TestGetAsSlice %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-SET
func TestSet(t *testing.T) {
	table := []struct {
		req     SetRequest
		dst     any
		want    string
		wantErr error
	}{
		{valuesReq1, &SetData1{}, `{"A":"ten"}`, nil},
		{valReqBool(true), &SetData1{}, `{"E":true}`, nil},
		{valReqBool(false), &SetData1{}, `{}`, nil},
		{valReqBool("true"), &SetData1{}, `{"E":true}`, nil},
		{valReqBool("t"), &SetData1{}, `{"E":true}`, nil},
		{valReqBool("false"), &SetData1{}, `{}`, nil},
		{valReqBool("f"), &SetData1{}, `{}`, nil},
		{valReqAny("D", 10.0, 0), &SetData1{}, `{"D":10}`, nil},
		{valReqAny("D", 10, 0), &SetData1{}, `{}`, fmt.Errorf("wrong type")},
		{valReqAny("D", 10, Fuzzy), &SetData1{}, `{"D":10}`, nil},
		{valReqJson("F", `[2, 4]`, 0), &SetData1{}, `{"F":[2,4]}`, nil},
	}
	for i, v := range table {
		haveErr := Set(v.req, v.dst)
		haveB, err := json.Marshal(v.dst)
		oferrors.Panic(err)
		have := string(haveB)

		if err := jacl.RunErr(haveErr, v.wantErr); err != nil {
			t.Fatalf("TestSet %v %v", i, err.Error())
		} else if have != v.want {
			t.Fatalf("TestSet %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-UNWRAP-VALUE-TO-ANY
func TestUnwrapValueToAny(t *testing.T) {
	table := []struct {
		v       any
		want    string
		wantErr error
	}{
		{tenStr, `string "ten"`, nil},
		{&tenStr, `string "ten"`, nil},
		{newAnyString("s"), `string "s"`, nil},
		{tenInt, `int 10`, nil},
		{&tenInt, `int 10`, nil},
		{tenFloat64, `float64 10`, nil},
		{&tenFloat64, `float64 10`, nil},
		{stringMap, `map[string]string {"a":"b"}`, nil},
		{&stringMap, `map[string]string {"a":"b"}`, nil},
	}
	for i, v := range table {
		haveV, haveErr := unwrapValueToAny(v.v)
		haveB, jsonErr := json.Marshal(haveV)
		oferrors.Panic(jsonErr)
		have := fmt.Sprintf("%T %v", haveV, string(haveB))

		if haveErr != v.wantErr {
			t.Fatalf("TestUnwrapValueToAny %v has err \"%v\" but wanted \"%v\"", i, haveErr, v.wantErr)
		} else if have != v.want {
			t.Fatalf("TestUnwrapValueToAny %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// CONST and VAR

var (
	filterMap1 = map[string]string{
		"A": "name",
	}
)

// Testing data. non-const because we use the address for some tests.
var (
	tenStr             = "ten"
	tenInt             = 10
	tenFloat64 float64 = 10.0
	stringMap          = map[string]string{"a": "b"}
)

var (
	valuesReq1 = SetRequest{
		FieldNames: []string{"A"},
		NewValues:  []any{&tenStr},
	}
)

// ---------------------------------------------------------
// SUPPORT interfaces

type Flattener interface {
	Flatten() string
}

// ---------------------------------------------------------
// SUPPORT types (Get)

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

func newHandler1() GetHandler {
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
	oferrors.Panic(err)
	sb.WriteString(string(b))
	sb.WriteString(",")
	b, err = json.Marshal(h.values)
	oferrors.Panic(err)
	sb.WriteString(string(b))
	return sb.String()
}

// ---------------------------------------------------------
// SUPPORT types (Set)

type SetData1 struct {
	A string  `json:",omitempty"`
	B int     `json:",omitempty"`
	C int64   `json:",omitempty"`
	D float64 `json:",omitempty"`
	E bool    `json:",omitempty"`
	F []int64 `json:",omitempty"`
}

// ---------------------------------------------------------
// MACROS (Set)

func newAnyString(s string) any {
	n := new(any)
	*n = s
	return n
}

func valReqAny(name string, v any, flags uint8) SetRequest {
	return SetRequest{
		FieldNames: []string{name},
		NewValues:  []any{v},
		Flags:      flags,
	}
}

func valReqBool(v any) SetRequest {
	return SetRequest{
		FieldNames: []string{"E"},
		NewValues:  []any{v},
	}
}

func valReqJson(name string, v any, flags uint8) SetRequest {
	return SetRequest{
		FieldNames: []string{name},
		NewValues:  []any{v},
		Assigns:    []SetFunc{SetJson},
		Flags:      flags,
	}
}

// ---------------------------------------------------------
// GET HANDLER SUPPORT

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
	oferrors.Panic(err)
	return string(b)
}

func flattenMap(m map[string]any) string {
	var sb strings.Builder
	_, err := sb.WriteString("{")
	oferrors.Panic(err)
	for i, t := range orderMap(m) {
		if i > 0 {
			_, err = sb.WriteString(",")
			oferrors.Panic(err)
		}
		_, err = sb.WriteString(fmt.Sprintf("\"%v\":", t.A))
		oferrors.Panic(err)
		b, err := json.Marshal(t.B)
		oferrors.Panic(err)
		_, err = sb.WriteString(string(b))
		oferrors.Panic(err)
	}
	_, err = sb.WriteString("}")
	oferrors.Panic(err)
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
