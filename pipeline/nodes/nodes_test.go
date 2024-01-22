package nodes

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/hackborn/onefunc/pipeline"
)

// ---------------------------------------------------------
// TEST-LOAD-FILE
func TestLoadFile(t *testing.T) {
	table := []struct {
		pipeline string
		want     pinsCompare
		wantErr  error
	}{
		{`graph (loadfile(Glob="` + testDataShortGlob + `"))`, newPinsCompare("content", "data", "a", "content", "data", "b", "content", "data", "c"), nil},
	}
	for i, v := range table {
		p, err := pipeline.Compile(v.pipeline)
		if err != nil {
			t.Fatalf("TestLoadFile %v compile err %v", i, err)
		}
		output, haveErr := pipeline.Run(p, nil)
		cmpErr := v.want.Compare(output)
		if v.wantErr == nil && haveErr != nil {
			t.Fatalf("TestLoadFile %v expected no error but has %v", i, haveErr)
		} else if v.wantErr != nil && haveErr == nil {
			t.Fatalf("TestLoadFile %v has no error but exptected %v", i, v.wantErr)
		} else if cmpErr != nil {
			t.Fatalf("TestLoadFile %v comparison error: %v", i, cmpErr)
		}
	}
}

// ---------------------------------------------------------
// PINS COMPARE

// Create a new comparison object on the provided state.
// You need to supply a pin data type followed by whatever
// parameters you want to check.
// Types:
// "content" params "name" (string), "data" (string)
// example:
// "content", "name", "filename.txt"
func newPinsCompare(as ...any) pinsCompare {
	ans := pinsCompare{}
	var cmp pinDataCmp
	var key string
	for _, a := range as {
		newCmp := newCmpFrom(a)
		if newCmp != nil {
			cmp = newCmp
			key = ""
			ans.pins = append(ans.pins, cmp)
			continue
		}
		if cmp == nil {
			panic(fmt.Sprintf("setting value %v without a pinCmp", a))
		}
		if key == "" {
			s, ok := a.(string)
			if !ok {
				panic(fmt.Sprintf("missing key string, instead have %v", a))
			}
			key = s
		} else {
			err := cmp.Assign(key, a)
			if err != nil {
				panic(fmt.Sprintf("error in assign: %v", err))
			}
			key = ""
		}
	}
	return ans
}

func newCmpFrom(a any) pinDataCmp {
	if s, ok := a.(string); ok {
		switch s {
		case "content":
			return &contentCmp{}
		}
	}
	return nil
}

type pinsCompare struct {
	pins []pinDataCmp
}

func (c *pinsCompare) Compare(output *pipeline.RunOutput) error {
	outputCount := 0
	if output != nil {
		outputCount = len(output.Pins)
	}
	if len(c.pins) < 1 && outputCount < 1 {
		return nil
	}
	if len(c.pins) != outputCount {
		return fmt.Errorf("pin mismatch, have %v but want %v", outputCount, len(c.pins))
	}
	for i, pin := range c.pins {
		err := pin.Compare(output.Pins[i])
		if err != nil {
			return err
		}
	}
	return nil
}

type pinDataCmp interface {
	Assign(key string, value any) error
	Compare(pipeline.PinData) error
}

type contentCmp struct {
	name *string
	data *string
}

func (c *contentCmp) Assign(key string, value any) error {
	switch key {
	case "name":
		if vs, ok := value.(string); ok {
			c.name = &vs
		}
	case "data":
		if vs, ok := value.(string); ok {
			c.data = &vs
		}
	default:
		return fmt.Errorf("unknown key %v", key)
	}
	return nil
}

func (c *contentCmp) Compare(pin pipeline.PinData) error {
	cd, ok := pin.(*pipeline.ContentData)
	if !ok {
		return fmt.Errorf("mismatched pin types, have contentCmp but supplied %t", pin)
	}
	if c.name != nil && *c.name != cd.Name {
		return fmt.Errorf("mismatched names, have %v but want %v", cd.Name, *c.name)
	}
	if c.data != nil && *c.data != cd.Data {
		return fmt.Errorf("mismatched data, have %v but want %v", cd.Data, *c.data)
	}
	return nil
}

// ---------------------------------------------------------
// SUPPORT

// Globs
var (
	testDataDomainGlob = filepath.Join(".", "test_data", "domain_*")
	testDataShortGlob  = filepath.Join(".", "test_data", "short_*")
)
