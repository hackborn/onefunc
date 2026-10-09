package hub

import (
	"fmt"
	"sync/atomic"

	"github.com/hackborn/onefunc/cfg"
)

type BuildArgs struct {
	GlobalSettings cfg.Settings
	// Steps are arbitrary functions run after services have
	// been created from the config files. Clients can use
	// a step to add or modify services.
	Steps []BuildStepFunc
}

// Build read the services config from the settings,
// opens them and returns the services.
//
// Example supplied settings source file:
//
//	{
//	 "background": {
//	   "workers": [4, 1]
//	 }
//	}
func Build(args BuildArgs, settings cfg.Settings) (Services, error) {
	waiting := newServices()
	opening := newServices()
	onargs := OnNewServiceArgs{done: &atomic.Bool{},
		openingServices: opening}

	// Install configured services
	for name := range settings.AllKeys() {
		fac, err := getFactory(name)
		if err != nil {
			return nil, err
		}
		service, err := fac.newService(args, settings.Subset(name))
		if err != nil {
			return nil, fmt.Errorf("factory %v: %w", name, err)
		}
		waiting.all[name] = serviceEntry{service: service,
			depedencies: fac.Dependencies}
		if fac.OnNewServiceFn != nil {
			onargs.Name, onargs.Service = name, service
			fac.OnNewServiceFn(onargs)

		}
	}
	onargs.done.Store(true)

	// Run steps
	stepArgs := StepArgs{GlobalSettings: args.GlobalSettings,
		Services:  waiting,
		_services: waiting,
	}
	for _, fn := range args.Steps {
		if err := fn(stepArgs); err != nil {
			return nil, err
		}
	}

	return open(settings, args.GlobalSettings, waiting, opening)
}

// ---------------------------------------------------------
// BUILD STEP

// BuildStepFunc gets called during build time. It provides
// a client a direct hook to create or alter services.
type BuildStepFunc func(StepArgs) error

type StepArgs struct {
	GlobalSettings cfg.Settings
	Services       Services
	_services      *_services
}

func (a StepArgs) AddService(name string, value any) {
	a._services.all[name] = serviceEntry{service: value}
}
