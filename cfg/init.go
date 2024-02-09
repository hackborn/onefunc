package cfg

import (
	"errors"
)

// func InitErr answers a composition of errors generated
// during init() funcs, if any.
func InitErr() error {
	if appDataPath == "" {
		return errors.Join(initErr, errors.New("platform did not set appDataPath"))
	}
	return initErr
}

// func AddInitErr adds an error to the init error list.
// This is checked only during startup, to deal with any
// errors generated by init() funcs.
func AddInitErr(err error) {
	if err == nil {
		return
	}
	if initErr == nil {
		initErr = err
	} else {
		initErr = errors.Join(initErr, err)
	}
}

var (
	initErr error = nil
)
