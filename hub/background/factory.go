package background

import (
	"github.com/hackborn/onefunc/cfg"

	"github.com/hackborn/onefunc/hub"
)

const (
	ServiceName = "background"
)

func init() {
	fac := hub.Factory{}
	fac.NewServiceFn = func(args hub.BuildArgs, settings cfg.Settings) (any, error) {
		return newService(settings), nil
	}
	hub.RegisterFactory(ServiceName, fac)
}
