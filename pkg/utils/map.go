package utils

// MapToList will change a map to a list.
func MapToList[T1 comparable, T2 any](m map[T1]T2) []T2 {
	l := make([]T2, 0, len(m))
	for _, v := range m {
		l = append(l, v)
	}
	return l
}
