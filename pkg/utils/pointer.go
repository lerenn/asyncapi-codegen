package utils

// ToPointer returns a pointer to the given value.
func ToPointer[T any](t T) *T {
	return &t
}

// ToValue returns the value pointed by the given pointer.
func ToValue[T any](t *T) T { //nolint:ireturn
	if t == nil {
		t = new(T)
	}
	return *t
}
