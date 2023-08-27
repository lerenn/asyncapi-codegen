package utils

// ToReference returns a pointer to the given value
func ToReference[T any](t T) *T {
	return &t
}
