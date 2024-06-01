package platform

// Set platform properties. This should be done once, at
// at app startup or in an init.
func Set(keys ...any) {
	for _, key := range keys {
		switch k := key.(type) {
		case _setAppName:
			appName = k.name
		}
	}
}

// SetAppName sets the global application name. They return
// value is a parameter for Set().
func SetAppName(name string) any {
	return _setAppName{name: name}
}

var appName string

type _setAppName struct {
	name string
}
