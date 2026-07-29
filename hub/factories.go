package hub

import (
	"fmt"

	"github.com/hackborn/onefunc/cfg"
)

type NewServiceFunc func(BuildArgs, cfg.Settings) (any, error)

type Factory struct {
	NewServiceFn NewServiceFunc
	// Dependencies is a list of services I am dependent on.
	// This will determine opening and closing order.
	// TODO: Opening not implemented yet, clients should still
	// return error in open when a dependency isn't ready.
	Dependencies []string
}

func (f Factory) newService(args BuildArgs, settings cfg.Settings) (any, error) {
	if f.NewServiceFn == nil {
		return nil, fmt.Errorf("no service func")
	}
	return f.NewServiceFn(args, settings)
}

func RegisterFactory(name string, factory Factory) error {
	return factories.Register(name, factory)
}

func getFactory(name string) (Factory, error) {
	if f, ok := factories.all[name]; ok {
		return f, nil
	}
	return Factory{}, fmt.Errorf("No factory named %v", name)
}

type Factories struct {
	all map[string]Factory
}

func (r *Factories) Register(name string, factory Factory) error {
	if _, ok := r.all[name]; ok {
		return fmt.Errorf("Factory %v already registered", name)
	}
	r.all[name] = factory
	return nil
}

func newFactories() *Factories {
	all := make(map[string]Factory)
	return &Factories{all: all}
}

var factories = newFactories()
