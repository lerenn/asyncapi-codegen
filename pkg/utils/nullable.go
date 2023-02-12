package utils

func ToNullable[T any](t T) *T {
	return &t
}
