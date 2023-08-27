package utils

import "unicode"

// UpperFirstLetter returns the given string with the first letter in uppercase
func UpperFirstLetter(str string) string {
	r := []rune(str)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
