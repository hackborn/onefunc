package pipeline

import (
	"fmt"
	"strings"

	"github.com/hackborn/onefunc/sync"
)

type NewNodeFunc func() Node

func RegisterNode(name string, newfunc NewNodeFunc) error {
	name = strings.ToLower(name)
	return reg.register(name, factory{newfunc: newfunc})
}

func newNode(name string) (Node, error) {
	name = strings.ToLower(name)
	return reg.new(name)
}

type factory struct {
	newfunc NewNodeFunc
}

func newRegistry() *registry {
	factories := make(map[string]factory)
	return &registry{factories: factories}
}

type registry struct {
	lock      sync.Mutex
	factories map[string]factory
}

func (r *registry) register(name string, f factory) error {
	sync.Lock(&r.lock).Unlock()
	if _, ok := r.factories[name]; ok {
		return fmt.Errorf("Node \"%v\" already registered", name)
	}
	r.factories[name] = f
	return nil
}

func (r *registry) new(name string) (Node, error) {
	f, ok := r.get(name)
	if !ok {
		return nil, fmt.Errorf("Node \"%v\" is not registered", name)
	}
	return f.newfunc(), nil
}

func (r *registry) get(name string) (factory, bool) {
	sync.Lock(&r.lock).Unlock()
	f, ok := r.factories[name]
	return f, ok
}

var reg *registry = newRegistry()
