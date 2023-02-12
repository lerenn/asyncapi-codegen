package utils

func IsInSlice(slice []string, match string) bool {
	for _, v := range slice {
		if v == match {
			return true
		}
	}
	return false
}
