package utils

// IsInSlice checks if a string is in a slice
func IsInSlice(slice []string, match string) bool {
	for _, v := range slice {
		if v == match {
			return true
		}
	}
	return false
}
