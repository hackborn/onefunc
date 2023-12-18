package extract

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"
)

// ---------------------------------------------------------
// TEST-TAIL
func TestTail(t *testing.T) {
	table := []struct {
		h    Handler
		want string
	}{
		{nil, ""},
		{&Chainer1{A: "s"}, "s"},
		{newChainer1("s"), "s"},
		{newChainer1("s", "t"), "t"},
		{newChainer1("s", "t", "v"), "v"},
	}
	for i, v := range table {
		c := tail(v.h)
		have := ""
		if last, ok := c.(*Chainer1); ok {
			have = last.A
		}

		if v.h == nil && c != nil {
			t.Fatalf("TestTail %v expected nil but has %v", i, c)
		} else if v.h != nil && c == nil {
			t.Fatalf("TestTail %v expected value but has nil", i)
		} else if have != v.want {
			t.Fatalf("TestTail %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-VALUES
func TestValues(t *testing.T) {
	table := []struct {
		src     any
		handler Handler
		want    string
	}{
		{Data1{A: "n"}, newHandler1(), `{"A":"n"}`},
		{Data1{A: "n"}, WithMap(newHandler1(), remap1), `{"name":"n"}`},
		{Data1{A: "n"}, WithFilterMap(newHandler1(), remap1), `{"name":"n"}`},
		{Data2{A: "n", B: 10}, newHandler1(), `{"A":"n","B":10}`},
		{Data2{A: "n", B: 10}, WithMap(newHandler1(), remap1), `{"B":10,"name":"n"}`},
		{Data2{A: "n", B: 10}, WithFilterMap(newHandler1(), remap1), `{"name":"n"}`},
	}
	for i, v := range table {
		Values(v.src, v.handler)
		have := v.handler.(Flattener).Flatten()

		if have != v.want {
			t.Fatalf("TestValues %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// SUPPORT var and const

var (
	remap1 = map[string]string{
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

func newHandler1() Handler {
	results := make(map[string]any)
	return &Handler1{results: results}
}

type Handler1 struct {
	results map[string]any
}

func (h *Handler1) Handle(name string, value any) {
	h.results[name] = value
}

func (h *Handler1) Flatten() string {
	return flattenMap(h.results)
}

type Chainer1 struct {
	Chain
	A string
}

func (c *Chainer1) Handle(name string, value any) {
}

func newChainer1(all ...string) Handler {
	chain := make([]Chainer, 0, len(all))
	for _, a := range all {
		chain = append(chain, &Chainer1{A: a})
	}
	return NewChain(chain...)
}

// ---------------------------------------------------------
// SUPPORT

func (h *remapHandler) Flatten() string {
	return h.Next.(Flattener).Flatten()
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
