package utils

import "unicode"

func UpperFirstLetter(str string) string {
	r := []rune(str)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
