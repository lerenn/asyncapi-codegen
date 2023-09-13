package utils

// Must panics if the passed error is not nil
//
// The purpose of this function is ignoring error checking when writing tests
// and only tests. This should not be used on production code.
func Must[T any](a T, err error) T { //nolint:ireturn,nolintlint
	if err != nil {
		panic(err)
	}

	return a
}
