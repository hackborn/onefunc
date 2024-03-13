package nxt

import (
	"github.com/hackborn/onefunc/cfg"
)

func Chain(args ChainArgs, newHandlers ...NewHandlerFunc) Pipeline {
	p := newPipeline()
	wd := &_wrapperData{}
	nhargs := NewHandlerArgs{Settings: args.Settings}
	nhargs.Next = NullHandler()
	for i := len(newHandlers) - 1; i >= 0; i-- {
		h, output := newHandlers[i](nhargs)
		// Chains are designed with optional components that might
		// be disabled, so skip them.
		if h != nil {
			// The output might trigger a housekeeping handler.
			if output != nil {
				if output.Factory != nil {
					fn := wrapperFunc(output.Factory, wd, h)
					nhargs.Next = fn
				}
			}
			nhargs.Next = h
		}
	}
	p.heads = append(p.heads, nhargs.Next)
	return p
}

type ChainArgs struct {
	Settings cfg.Settings
}
