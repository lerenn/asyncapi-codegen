package templates

import (
	"fmt"
	"html/template"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

// NamifyWithoutParams will convert a sentence to a golang conventional type name.
// and will remove all parameters that can appear between '{' and '}'.
func NamifyWithoutParams(sentence string) string {
	// Remove parameters
	re := regexp.MustCompile("{[^()]*}")
	sentence = string(re.ReplaceAll([]byte(sentence), []byte("_")))

	return Namify(sentence)
}

// Namify will convert a sentence to a golang conventional type name.
func Namify(sentence string) string {
	// Check if empty
	if len(sentence) == 0 {
		return sentence
	}

	// Upper letters that are preceded with an underscore
	previous := '_'
	for i, r := range sentence {
		if !unicode.IsLetter(previous) && !unicode.IsDigit(previous) {
			sentence = sentence[:i] + strings.ToUpper(string(r)) + sentence[i+1:]
		}
		previous = r
	}

	// Remove everything except alphanumerics
	re := regexp.MustCompile("[^a-zA-Z0-9]")
	sentence = string(re.ReplaceAll([]byte(sentence), []byte("")))

	// Remove leading numbers
	re = regexp.MustCompile("^[0-9]+")
	sentence = string(re.ReplaceAll([]byte(sentence), []byte("")))

	// Upper first letter
	sentence = strings.ToUpper(sentence[:1]) + sentence[1:]

	return sentence
}

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// SnakeCase will convert a sentence to snake case.
func SnakeCase(sentence string) string {
	snake := matchFirstCap.ReplaceAllString(sentence, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// ReferenceToTypeName will convert a reference to a type name in the form of
// golang conventional type names.
func ReferenceToTypeName(ref string) string {
	parts := strings.Split(ref, "/")
	return Namify(parts[3])
}

// HasField will check if a struct has a field with the given name.
func HasField(v any, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	return rv.FieldByName(name).IsValid()
}

// DescribeStruct will describe a struct in a human readable way using `%+v`
// format from the standard library.
func DescribeStruct(st any) string {
	return fmt.Sprintf("%+v", st)
}

// MultiLineComment will prefix each line of a comment with "// " in order to
// make it a valid multiline golang comment.
func MultiLineComment(comment string) string {
	comment = strings.TrimSuffix(comment, "\n")
	return strings.ReplaceAll(comment, "\n", "\n// ")
}

// Args is a function used to pass arguments to templates.
func Args(vs ...any) []any {
	return vs
}

// HelpersFunctions returns the functions that can be used as helpers
// in a golang template.
func HelpersFunctions() template.FuncMap {
	return template.FuncMap{
		"namifyWithoutParam":  NamifyWithoutParams,
		"namify":              Namify,
		"snakeCase":           SnakeCase,
		"referenceToTypeName": ReferenceToTypeName,
		"hasField":            HasField,
		"describeStruct":      DescribeStruct,
		"multiLineComment":    MultiLineComment,
		"args":                Args,
	}
}
