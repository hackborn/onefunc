package background

import (
	"fmt"
	"sync/atomic"
)

// This file provides a menchanism for clients to add their own
// background routines to the service, managed by the service.
// This can only be done during the Open() stage of the
// service lifecycle.

// OpenerService is only available during Open(). It allows clients
// to add managed background tasks.
type OpenerService interface {
	// AddChannel adds a new channel to the background service.
	// Access it by calling RunOn with the returned ID.
	// The remoteArg and localArg will be available in the
	// remote and local args.
	AddChannel(AddChannelArgs) int

	// Go will run the func on the background thread. Closer
	// will be closed when the service ends, and the func
	// should check the closed state and return when closed.
	// TODO: I think this can be deprecated, and just use AddChannel.
	Go(func(), Closer)

	// SendLocalFunc answers the function used to send the
	// task to the local thread. Some apps might need to replace
	// this with something appropriate to the UI system.
	SendLocalFunc() SendLocalFunc
}

type AddChannelArgs struct {
	// RemoteData will be available from RemoteArgs.Get.
	RemoteData any
	// LocalData will be available from LocalArgs.Get.
	LocalData any
	// Max will cause the worker threads to dynamically
	// scale, from 0 to max.
	// TODO: Unimplemented.
	//Max int
}

// ConfigService is only available during the build stage. It lets
// clients pass in data to the service.
type ConfigService interface {
	SetRunLocalFunc(RunLocalFunc)
	SetLocalData(any)
}

type Closer interface {
	Close()
	IsClosed() bool
}

// ChannelCloser is a convenience for the most common case:
// Clients have a channel and ends when it closes.
type ChannelCloser[T any] struct {
	closed atomic.Bool
	c      chan T
}

func NewChannelCloser[T any](c chan T) *ChannelCloser[T] {
	return &ChannelCloser[T]{c: c}
}

func (c *ChannelCloser[T]) Close() {
	c.closed.Store(true)
	close(c.c)
}

func (c *ChannelCloser[T]) IsClosed() bool {
	// The closed state is used so clients don't have to
	// wait for the channel to drain before they know to
	// return. They can if they want, but most won't.
	return c.c == nil || c.closed.Load() == true
}

type openerService struct {
	s *service
}

func (s *openerService) AddChannel(args AddChannelArgs) int {
	panicService(s.s, "opener")
	return s.s.doAddChannel(args)
}

func (s *openerService) Go(fn func(), closer Closer) {
	panicService(s.s, "opener")
	s.s.doGo(fn, closer)
}

func (s *openerService) SendLocalFunc() SendLocalFunc {
	panicService(s.s, "opener")
	return s.s.doSendLocalFunc()
}

func (s *openerService) BuildFinished() {
	s.s = nil
}

type configService struct {
	s *service
}

func (s *configService) BuildFinished() {
	s.s = nil
}

func (s *configService) SetRunLocalFunc(fn RunLocalFunc) {
	panicService(s.s, "config")
	s.s.runLocalFn = fn
}

func (s *configService) SetLocalData(d any) {
	panicService(s.s, "config")
	s.s.localArgs.data = append(s.s.localArgs.data, d)
}

func panicService(s *service, name string) {
	if s == nil {
		panic(fmt.Sprintf("background %s after the build has finished", name))
	}
}
