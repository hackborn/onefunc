package assign

import (
	"encoding/json"
	"fmt"
	"testing"
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
		panicErr(jsonErr)
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
	}
	for i, v := range table {
		haveErr := Values(v.req, v.dst)
		haveB, err := json.Marshal(v.dst)
		panicErr(err)
		have := string(haveB)

		if haveErr != v.wantErr {
			t.Fatalf("TestValues %v has err \"%v\" but wanted \"%v\"", i, haveErr, v.wantErr)
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
}

// ---------------------------------------------------------
// MACROS

func newAnyString(s string) any {
	n := new(any)
	*n = s
	return n
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
