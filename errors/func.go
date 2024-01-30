package errors

import (
	"log"
)

// First answers the first non-nil error in the list.
func First(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// Panic panics if any error is non-nil.
func Panic(errs ...error) {
	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}

// LogFatal does a log.Fata if any error is non-nil.
func LogFatal(errs ...error) {
	for _, err := range errs {
		if err != nil {
			log.Fatal(err)
		}
	}
}
