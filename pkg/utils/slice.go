package utils

// RemoveDuplicateFromSlice removes duplicate values from a slice
func RemoveDuplicateFromSlice[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// IsInSlice checks if a string is in a slice
func IsInSlice(slice []string, match string) bool {
	for _, v := range slice {
		if v == match {
			return true
		}
	}
	return false
}
