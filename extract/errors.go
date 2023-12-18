package extract

// firstErr answers the first non-nil error in the list.
func firstErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
