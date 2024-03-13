package nxt

import (
	"github.com/hackborn/onefunc/cfg"
)

// NewHandlerFunc returns a Handler function
type NewHandlerFunc func(args NewHandlerArgs) (Handler, *NewHandlerOutput)

type NewHandlerArgs struct {
	// Acccess to the application settings.
	Settings cfg.Settings

	// Next is the next handler in the list. Never nil.
	// New handlers aren't required to send to a next, but
	// that does turn them into a destination that breaks the chain.
	Next Handler
}

type NewHandlerOutput struct {
	Factory Factory
}

// Handler handles an event.
type Handler func(args HandlerArgs, event any)

type HandlerArgs struct {
	Data any

	runnerId int64
}
