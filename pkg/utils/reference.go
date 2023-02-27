package utils

func ToReference[T any](t T) *T {
	return &t
}
