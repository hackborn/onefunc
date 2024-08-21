package nxt

import (
	"reflect"
	"testing"
)

// ---------------------------------------------------------
// TEST-CHAIN-2
func TestChain2(t *testing.T) {
	f := func(nhfn ...NewHandler2Func) {
		t.Helper()

		c2a := Chain2Args{}
		Chain2(c2a, nhfn...)
	}
	f(newTextFunc("hi"), addTextFunc("", "!"))
	panic("no")
}

// ---------------------------------------------------------
// TEST-CHAIN
func _TestChain(t *testing.T) {
	table := []struct {
		args ChainArgs
		src  any
		want any
	}{
		{ChainArgs{}, nil, nil},
	}
	for i, v := range table {
		c := &_capture{}
		p := Chain(v.args, newCaptureFunc(c))
		p.Run(v.src)
		have := c.captured
		if reflect.DeepEqual(v.want, have) != true {
			t.Fatalf("TestChain %v has \"%v\" but wants \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// HANDLERS

func newCaptureFunc(c *_capture) NewHandlerFunc {
	return func(args NewHandlerArgs) (Handler, *NewHandlerOutput) {
		return c.Handle, nil
	}
}

type _capture struct {
	captured []any
}

func (c *_capture) Handle(args HandlerArgs, event any) {
	c.captured = append(c.captured, event)
}

// ---------------------------------------------------------
// CONST and VAR
