package pipeline

import (
	"fmt"
	"sync"

	"github.com/hackborn/onefunc/lock"
)

type NewNodeFunc func() Node

func Register(name string, newfunc NewNodeFunc) error {
	return reg.register(name, factory{newfunc: newfunc})
}

type factory struct {
	newfunc NewNodeFunc
}

type registry struct {
	mu        sync.Mutex
	factories map[string]factory
}

func (r *registry) register(name string, f factory) error {
	lock.Locker(&r.mu).Unlock()
	if _, ok := r.factories[name]; ok {
		return fmt.Errorf("node \"%v\" already registered", name)
	}
	r.factories[name] = f
	return nil
}

var reg *registry = &registry{}
