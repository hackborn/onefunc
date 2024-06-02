package assign

import (
	"encoding/json"
	"fmt"
	"testing"

	oferrors "github.com/hackborn/onefunc/errors"
	"github.com/hackborn/onefunc/jacl"
)

// ---------------------------------------------------------
// TESTS

// TestUnwrapValueToAny
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
			t.Fatalf("TestUnwrapValue %v has err \"%v\" but wanted \"%v\"", i, haveErr, v.wantErr)
		} else if have != v.want {
			t.Fatalf("TestUnwrapValue %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		}
	}
}

// TestValues
func TestValues(t *testing.T) {
	table := []struct {
		req     ValuesRequest
		dst     any
		want    string
		wantErr error
	}{
		{valuesReq1, &Data1{}, `{"A":"ten"}`, nil},
		{valReqBool(true), &Data1{}, `{"E":true}`, nil},
		{valReqBool(false), &Data1{}, `{}`, nil},
		{valReqBool("true"), &Data1{}, `{"E":true}`, nil},
		{valReqBool("t"), &Data1{}, `{"E":true}`, nil},
		{valReqBool("false"), &Data1{}, `{}`, nil},
		{valReqBool("f"), &Data1{}, `{}`, nil},
		{valReqAny("D", 10.0, 0), &Data1{}, `{"D":10}`, nil},
		{valReqAny("D", 10, 0), &Data1{}, `{}`, fmt.Errorf("wrong type")},
		{valReqAny("D", 10, Fuzzy), &Data1{}, `{"D":10}`, nil},
		{valReqJson("F", `[2, 4]`, 0), &Data1{}, `{"F":[2,4]}`, nil},
	}
	for i, v := range table {
		haveErr := Values(v.req, v.dst)
		haveB, err := json.Marshal(v.dst)
		oferrors.Panic(err)
		have := string(haveB)

		if err := jacl.RunErr(haveErr, v.wantErr); err != nil {
			t.Fatalf("TestValues %v %v", i, err.Error())
		} else if have != v.want {
			t.Fatalf("TestValues %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TYPES

type Data1 struct {
	A string  `json:",omitempty"`
	B int     `json:",omitempty"`
	C int64   `json:",omitempty"`
	D float64 `json:",omitempty"`
	E bool    `json:",omitempty"`
	F []int64 `json:",omitempty"`
}

// ---------------------------------------------------------
// MACROS

func newAnyString(s string) any {
	n := new(any)
	*n = s
	return n
}

func valReqAny(name string, v any, flags uint8) ValuesRequest {
	return ValuesRequest{
		FieldNames: []string{name},
		NewValues:  []any{v},
		Flags:      flags,
	}
}

func valReqBool(v any) ValuesRequest {
	return ValuesRequest{
		FieldNames: []string{"E"},
		NewValues:  []any{v},
	}
}

func valReqJson(name string, v any, flags uint8) ValuesRequest {
	return ValuesRequest{
		FieldNames: []string{name},
		NewValues:  []any{v},
		Assigns:    []AssignFunc{AssignJson},
		Flags:      flags,
	}
}

// ---------------------------------------------------------
// CONST and VAR

// Testing data. non-const because we use the address for some tests.
var (
	tenStr             = "ten"
	tenInt             = 10
	tenFloat64 float64 = 10.0
	stringMap          = map[string]string{"a": "b"}
)

var (
	valuesReq1 = ValuesRequest{
		FieldNames: []string{"A"},
		NewValues:  []any{&tenStr},
	}
)
