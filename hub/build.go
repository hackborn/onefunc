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

func (a BuildArgs) AppendStep(stepFn BuildStepFunc) BuildArgs {
	a.Steps = append(a.Steps, stepFn)
	return a
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
	onargs := OnNewServiceArgs{done: &atomic.Bool{},
		services: waiting,
	}

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

	// Run the steps
	if err := runSteps(args, waiting); err != nil {
		return nil, err
	}

	// Check for any missing
	if err := verifyMissing(waiting, missing); err != nil {
		return nil, err
	}

	// Open the services
	ready, err := open(settings, args.GlobalSettings, waiting)
	if err != nil {
		return nil, err
	}

	// Remove the finishers
	removeBuildFinishers(ready)
	return ready, nil
}

func runSteps(args BuildArgs, s *_services) error {
	stepArgs := StepArgs{GlobalSettings: args.GlobalSettings,
		Services:  s,
		_services: s,
	}
	for _, fn := range args.Steps {
		if err := fn(stepArgs); err != nil {
			return err
		}
	}
	return nil
}

// verifyMissing returns an error if I don't have any items
// in the missing.
func verifyMissing(s *_services, missing []string) error {
	for _, name := range missing {
		if _, ok := s.get(name); !ok {
			return fmt.Errorf("No service named %v", name)
		}
	}
	return nil
}

// removeBuildFinishers removes all BuildFinisher services.
func removeBuildFinishers(s *_services) {
	fn := func(name string, service any) bool {
		if f, ok := service.(BuildFinisher); ok {
			f.BuildFinished()
			return true
		}
		return false
	}
	s.deleteFunc(fn)
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
