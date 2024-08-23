package reflect

import (
	"fmt"
)

// ---------------------------------------------------------
// HANDLING

// GetHandler is used to get name/value pairs from a struct.
type GetHandler interface {
	// Handle the pair in some fashion, returning the (potentially
	// changed, filtered etc. data). The basic system
	// is designed to work without errors; if a client needs
	// any error handling, it should track that in an internal state.
	Handle(name string, value any) (string, any)
}

// Slicer can return a slice.
type Slicer interface {
	Slice() []any
}

// Mapper can return a map.
type Mapper interface {
	Map() map[string]any
}

// ---------------------------------------------------------
// CHAINING

type Chain []GetHandler

func (c Chain) Handle(name string, value any) (string, any) {
	for _, h := range c {
		name, value = h.Handle(name, value)
		if name == "" {
			return "", nil
		}
	}
	return name, value
}

// NewChain takes a list of items and converts them into handlers,
// answering the new handler Chain. Items can be:
// * A Handler. It is added directly to the chain.
// * An *Opts (SlicerOpts, etc.). It is converted to the associated
// Handler and added to the chain.
// * A map[string]string. It is converted to a FilterMap Handler.
func NewChain(items ...any) Chain {
	chain := make(Chain, 0, len(items))
	for _, item := range items {
		var h GetHandler
		switch t := item.(type) {
		case GetHandler:
			h = t
		case FilterMapOpts:
			h = FilterMap(t)
		case *FilterMapOpts:
			h = FilterMap(*t)
		case MapOpts:
			h = Map(t)
		case *MapOpts:
			h = Map(*t)
		case SliceOpts:
			h = Slice(t)
		case *SliceOpts:
			h = Slice(*t)
		case map[string]string:
			h = FilterMap(FilterMapOpts{F: t})
		default:
			fmt.Printf("new chain unknown type %T\n", t)
		}
		if h != nil {
			chain = append(chain, h)
		}
	}
	return chain
}

// ---------------------------------------------------------
// HANDLER TYPES

type FilterMapOpts struct {
	F map[string]string
	// When passthrough is true, anything not filtered is passed on unchanged.
	Passthrough bool
}

func FilterMap(opts FilterMapOpts) GetHandler {
	return &filterMapHandler{opts: opts}
}

type SliceOpts struct {
	Assign  string
	Combine string
}

func Slice(opts SliceOpts) GetHandler {
	return &sliceHandler{opts: opts}
}

type MapOpts struct {
}

func Map(opts MapOpts) GetHandler {
	result := make(map[string]any)
	return &mapHandler{result: result, opts: opts}
}

// filterMapHandler
type filterMapHandler struct {
	opts FilterMapOpts
}

func (h *filterMapHandler) Handle(name string, value any) (string, any) {
	if n, ok := h.opts.F[name]; ok {
		return n, value
	} else if h.opts.Passthrough {
		return name, value
	} else {
		return "", nil
	}
}

// sliceHandler
type sliceHandler struct {
	result []any
	opts   SliceOpts
}

func (h *sliceHandler) Handle(name string, value any) (string, any) {
	if h.opts.Combine != "" && len(h.result) > 0 {
		h.result = append(h.result, h.opts.Combine)
	}
	h.result = append(h.result, name)
	if h.opts.Assign != "" {
		h.result = append(h.result, h.opts.Assign)
	}
	h.result = append(h.result, value)

	return name, value
}

func (h *sliceHandler) Slice() []any {
	return h.result
}

// mapHandler
type mapHandler struct {
	result map[string]any
	opts   MapOpts
}

func (h *mapHandler) Handle(name string, value any) (string, any) {
	if name == "" {
		return "", nil
	}
	h.result[name] = value
	return name, value
}

func (h *mapHandler) Map() map[string]any {
	return h.result
}

// ---------------------------------------------------------
// SUPPORT

// getLast answers the last handler in the (potential) chain
// that matches type T.
func getLast[T any](h GetHandler) (T, bool) {
	if t, ok := h.(T); ok {
		return t, true
	}
	switch t := h.(type) {
	case Chain:
		for i := len(t) - 1; i >= 0; i-- {
			if ans, ok := getLast[T](t[i]); ok {
				return ans, ok
			}
		}
	}
	var t T
	return t, false
}
