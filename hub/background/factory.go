package background

import (
	"github.com/hackborn/onefunc/cfg"

	"github.com/hackborn/onefunc/hub"
)

const (
	ServiceName = "background"
	OpenerName  = "background#opener"
)

func init() {
	fac := hub.Factory{}
	fac.NewServiceFn = func(args hub.BuildArgs, settings cfg.Settings) (any, error) {
		return newService(settings), nil
	}
	fac.OnNewServiceFn = func(args hub.OnNewServiceArgs) {
		// Add an opener service, so clients have access to a mechanism
		// during the Open stage to add their own background threads.
		if s, ok := args.Service.(*service); ok {
			openerS := &openerService{s: s}
			args.SetOpeningService(OpenerName, openerS)
		}
	}
	hub.RegisterFactory(ServiceName, fac)
}
