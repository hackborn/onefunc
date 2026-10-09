package hub

import (
	"cmp"
	"fmt"
	"iter"

	"github.com/hackborn/onefunc/cfg"
)

// Opener is called during the construction of
// the services list.
// Services can also implement io.Closer to close.
type Opener interface {
	// Open the service, returning an error. Clients
	// can return WaitingErr to indicate they are waiting
	// for a prerequisite service, but this error does not
	// impact functionality. Services are opened repeaetedly
	// until they have no error or no more services can open.
	Open(OpenArgs) error
}

type OpenArgs struct {
	Settings       cfg.Settings
	GlobalSettings cfg.Settings
	Services       Services
}

// Closing is called prior to Close(), so any services
// dependent on other services can cleanup. You can control
// the order of close operations in the Factory.
type Closing interface {
	Closing(Services)
}

// Updater is a contract that is not implemented
// in the services. Clients are required to extract
// and handle any services that Update.
type Updater interface {
	// Update is called every frame.
	Update()
}

var WaitingErr = fmt.Errorf("Waiting")

// open services. Repeatedly walk the list opening
// whomever doesn't return an error. Stop when all
// services are opened or we only get errors.
func open(settings, globalSettings cfg.Settings, a, opening *_services) (*_services, error) {
	b := newServices()
	composite := &compositeService{a: b, b: opening}
	var lastErr error
	lastSize := len(a.all) + 1
	for lastSize != len(a.all) {
		lastSize = len(a.all)

		for k, v := range a.all {
			if opener, ok := v.service.(Opener); ok {
				args := OpenArgs{Settings: settings.Subset(k),
					GlobalSettings: globalSettings,
					//					Services:       b,
					Services: composite,
				}
				if err := opener.Open(args); err == nil {
					a.moveTo(k, b)
				} else {
					lastErr = fmt.Errorf("%v: %w", k, err)
				}
			} else {
				a.moveTo(k, b)
			}
		}
	}
	opening.Close()
	if len(a.all) > 0 {
		// We failed to open everyone, which means the client
		// will receive an empty services, so at least clean up.
		b.Close()
		a.Close()
		return nil, cmp.Or(lastErr, fmt.Errorf("Can't open services"))
	}
	return b, nil
}

type compositeService struct {
	a, b *_services
}

func (s *compositeService) All() iter.Seq2[string, any] {
	fmt.Println("Unfinished: compositeService.All() only iterates a")
	return s.a.All()
}

func (s *compositeService) Close() error {
	err := s.a.Close()
	return cmp.Or(err, s.b.Close())
}

func (s *compositeService) get(name string) (any, bool) {
	if res, ok := s.a.get(name); ok {
		return res, ok
	}
	return s.b.get(name)
}
