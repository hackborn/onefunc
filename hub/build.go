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
	// Settings might be configured for services that are added in
	// build steps, so check for missing services, and only error
	// if they are still missing after the build steps.
	missing := []string{}
	for name := range settings.AllKeys() {
		if !hasFactory(name) {
			missing = append(missing, name)
			continue
		}
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

	// Verify any missing have been added
	for _, name := range missing {
		if _, ok := waiting.get(name); !ok {
			return nil, fmt.Errorf("No factory named %v", name)
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
