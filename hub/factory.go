package hub

import (
	"fmt"
	"sync/atomic"

	"github.com/hackborn/onefunc/cfg"
)

type NewServiceFunc func(BuildArgs, cfg.Settings) (any, error)

type Factory struct {
	// NewServiceFn lets clients return a new service.
	NewServiceFn NewServiceFunc
	// OnNewServiceFn (optional) gets called whenever
	// NewServiceFn returns a valid service. It can
	// be used to take additional actions.
	OnNewServiceFn OnNewServiceFunc
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
	return allFactories.register(name, factory)
}

func getFactory(name string) (Factory, error) {
	if f, ok := allFactories.all[name]; ok {
		return f, nil
	}
	return Factory{}, fmt.Errorf("No factory named %v", name)
}

// ---------------------------------------------------------
// ON NEW SERVICE

type OnNewServiceFunc func(OnNewServiceArgs)

type OnNewServiceArgs struct {
	// Name is the name of the service that was just added.
	Name string
	// Service is the added service.
	Service any

	// done is marked true when the framework is finished.
	// This is extra-protection against anyone capturing
	// the args and trying to use them past build time.
	done            *atomic.Bool
	services        *_services
	openingServices *_services
}

// SddService will set the service at name, replacing any existing.
func (a OnNewServiceArgs) SddService(name string, service any) {
	if a.services != nil && a.done != nil && a.done.Load() == false {
		a.services.all[name] = serviceEntry{service: service}
	}
}

// SetOpeningService will set a new opening service. Opeming services
// disappear after the Open stage.
func (a OnNewServiceArgs) SetOpeningService(name string, service any) {
	if a.openingServices != nil && a.done != nil && a.done.Load() == false {
		a.openingServices.all[name] = serviceEntry{service: service}
	}
}

// ---------------------------------------------------------
// FACTORIES

// factories stores the list of factories.
type factories struct {
	all map[string]Factory
}

func (r *factories) register(name string, factory Factory) error {
	if _, ok := r.all[name]; ok {
		return fmt.Errorf("Factory %v already registered", name)
	}
	r.all[name] = factory
	return nil
}

func newFactories() *factories {
	all := make(map[string]Factory)
	return &factories{all: all}
}

var allFactories = newFactories()
