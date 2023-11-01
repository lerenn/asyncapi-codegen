package templates

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
	"github.com/stoewer/go-strcase"
)

// Namify will convert a sentence to a golang conventional type name.
func Namify(sentence string) string {
	// Remove everything except alphanumerics and '_'
	re := regexp.MustCompile("[^a-zA-Z0-9_]")
	sentence = string(re.ReplaceAll([]byte(sentence), []byte("_")))

	// Remove leading numbers
	re = regexp.MustCompile("^[0-9]+")
	sentence = string(re.ReplaceAll([]byte(sentence), []byte("")))

	// Snake case to Upper Camel case
	sentence = strcase.UpperCamelCase(sentence)

	// Correct acronyms
	return correctAcronyms(sentence)
}

func correctAcronyms(sentence string) string {
	acronyms := []string{"ID"}
	for _, a := range acronyms {
		wronglyFormatedAcronym := strcase.UpperCamelCase(a)
		re := regexp.MustCompile(fmt.Sprintf("%s([A-Z]+|$)", wronglyFormatedAcronym))

		positions := re.FindAllStringIndex(sentence, -1)
		for _, p := range positions {
			sentence = sentence[:p[0]] + a + sentence[p[0]+len(a):]
		}
	}

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

	name := parts[3]
	if parts[2] == "messages" {
		name += "Message"
	} else if parts[2] == "schemas" {
		name += "Schema"
	}

	return Namify(name)
}

// referenceToSlicePath will convert a reference to a slice where each element is a
// step of the path.
func referenceToSlicePath(ref string) []string {
	ref = strings.ReplaceAll(ref, ".", "/")
	ref = strings.ReplaceAll(ref, "#", "")
	return strings.Split(ref, "/")[1:]
}

// ReferenceToStructAttributePath will convert a reference to a struct attribute
// path in the form of "a.b.c" where a, b and c are struct attributes in the
// form of golang conventional type names.
func ReferenceToStructAttributePath(ref string) string {
	path := referenceToSlicePath(ref)

	for k, v := range path {
		// If this is concerning the header, then it will be named "headers"
		if v == asyncapi.MessageTypeIsHeader.String() {
			v = "headers"
		}

		path[k] = Namify(v)
	}

	return strings.Join(path, ".")
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

// ChannelToMessageTypeName will convert a channel to a message type name in the
// form of golang conventional type names.
func ChannelToMessageTypeName(ch asyncapi.Channel) string {
	msg := ch.GetChannelMessage()

	if msg.Payload != nil || msg.OneOf != nil {
		return Namify(ch.Name) + "Message"
	}

	return ReferenceToTypeName(msg.Reference)
}

// IsRequired will check if a field is required in a asyncapi struct.
func IsRequired(schema asyncapi.Schema, field string) bool {
	return schema.IsFieldRequired(field)
}

// GenerateChannelPath will generate a channel path with the given channel.
func GenerateChannelPath(ch asyncapi.Channel) string {
	// If there is no parameter, then just return the path
	if ch.Parameters == nil {
		return fmt.Sprintf("%q", ch.Path)
	}

	parameterRegexp := regexp.MustCompile("{[^{}]*}")

	matches := parameterRegexp.FindAllString(ch.Path, -1)
	format := parameterRegexp.ReplaceAllString(ch.Path, "%v")

	sprint := fmt.Sprintf("fmt.Sprintf(%q, ", format)
	for _, m := range matches {
		sprint += fmt.Sprintf("params.%s,", Namify(m))
	}

	return sprint[:len(sprint)-1] + ")"
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

// OperationName returns `operationId` value from Publish or Subscribe operation if any.
// If no `operationID` exists — return provided default value (`name`).
func OperationName(channel *asyncapi.Channel) string {
	switch {
	case channel.Publish != nil && channel.Publish.OperationID != "":
		return channel.Publish.OperationID
	case channel.Subscribe != nil && channel.Subscribe.OperationID != "":
		return channel.Subscribe.OperationID
	default:
		return channel.Name
	}
}
