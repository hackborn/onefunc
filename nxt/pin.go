package nxt

import (
	"fmt"

	"github.com/hackborn/onefunc/cfg"
)

type Pin[T any] struct {
}

// What's the thinking here? We want types pins but also
// function handlers. So a graph like

// src -> go -> folder

// would become what?
// src needs to be given a handler like
// func[T Content](c Content)
// and then go needs to handle that. So I guess one way
// to handle this is if we have different ary functions i.e.
// type Handle1Func[T1 any](t1 T1)
// type Handle2Func[T1 any, T2 any](t1 T1, t2 T2)
// So then go, presuming it has two pins, would actually
// be a handler func like
// gonxt[T1 Content, T2 Content](t1 Content, t2 Content)
// and when we're composing, how do we know what the nxt
// is and where we get assigned?
// Also this makes no sense -- we wouldn't know what the other
// args are. So it really seems like functions need to return
// single-argument functions
// So the newhandler func can return I guess a slice of pins,
// and clients then know which one they are and can cast.

// Ok, the following implements a system where every newhandler
// needs to know the pin number it's connecting to and the method
// sig, but all casting happens in the construction and no
// interfaces are used during performance. It's also pretty confusing,
// with a lot of boilerplate (possibly reducable via generics, haven't
// thought about it much).
//
// So... what is missing from the more complex pipeline?
// * Waiting. Any events entering the system run through the entire
// thing. A given handler could take responsibility for waiting until
// events have arrived at every pin, but it would be stuck forever if
// they don't come. Certainly, there is nothing outside the system
// that is shuttling data around, and I'm not sure how there could be.
// * Per-thread data. The Pipeline essentially instantiates read-only
// nodes which are then performed in the context of a runner. Adding
// something similar would be through the handling args, but it couldn't
// be multi-stage, where you have a Start() and a Run().

// NewHandlerFunc returns a Handler function
type NewHandler2Func func(args NewHandler2Args) ([]any, *NewHandler2Output)

type NewHandler2Args struct {
	// Acccess to the application settings.
	Settings cfg.Settings

	// Next is the next handler in the list. Never nil.
	// New handlers aren't required to send to a next, but
	// that does turn them into a destination that breaks the chain.
	Next []any
}

type NewHandler2Output struct {
	Factory Factory
}

// Handler handles an event.
type Handler2[T any] func(args Handler2Args, event T)

type Handler2Args struct {
	Data any

	runnerId int64
}

func Chain2(args Chain2Args, newHandlers ...NewHandler2Func) Pipeline {
	p := newPipeline()
	//	wd := &_wrapperData{}
	nhargs := NewHandler2Args{Settings: args.Settings}
	//	nhargs.Next = NullHandler()
	output := func(a Handler2Args, s string) {
		fmt.Println("output", s)
	}
	nhargs.Next = []any{output}
	for i := len(newHandlers) - 1; i >= 0; i-- {
		nh := newHandlers[i]
		if nh == nil {
			continue
		}
		h, output := nh(nhargs)
		// Chains are designed with optional components that might
		// be disabled, so skip them.
		if h != nil {
			// The output might trigger a housekeeping handler.
			if output != nil {
				if output.Factory != nil {
					//					fn := wrapperFunc(output.Factory, wd, h)
					//					nhargs.Next = fn
				}
			}
			nhargs.Next = h
		}
	}
	h2a := Handler2Args{}
	for _, nxt := range nhargs.Next {
		if h, ok := nxt.(func(args2 Handler2Args)); ok {
			h(h2a)
		}
	}
	//	p.heads = append(p.heads, nhargs.Next)
	return p
}

type Chain2Args struct {
	Settings cfg.Settings
}

func newTextFunc(s string) NewHandler2Func {
	return func(args NewHandler2Args) ([]any, *NewHandler2Output) {
		nxt := func(a Handler2Args, s string) {
		}
		if len(args.Next) > 0 {
			if h, ok := args.Next[0].(func(a Handler2Args, s string)); ok {
				nxt = h
			}
		}
		fn := func(args2 Handler2Args) {
			nxt(args2, s)
		}
		return []any{fn}, nil
	}
}

func addTextFunc(pre, post string) NewHandler2Func {
	return func(args NewHandler2Args) ([]any, *NewHandler2Output) {
		nxt := func(a Handler2Args, s string) {
		}
		if len(args.Next) > 0 {
			if h, ok := args.Next[0].(func(a Handler2Args, s string)); ok {
				nxt = h
			}
		}
		fn := func(args2 Handler2Args, s string) {
			s = pre + s + post
			nxt(args2, s)
		}
		return []any{fn}, nil
	}
}
