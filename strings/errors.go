package strings

// firstErr answers the first non-nil error in the list.
func firstErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
