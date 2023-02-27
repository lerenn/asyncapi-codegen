package templates

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
	"github.com/stoewer/go-strcase"
)

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
		re := regexp.MustCompile(fmt.Sprintf("%s[A-Z]*", wronglyFormatedAcronym))

		positions := re.FindAllIndex([]byte(sentence), -1)
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

func SnakeCase(sentence string) string {
	snake := matchFirstCap.ReplaceAllString(sentence, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ReferenceToTypeName(ref string) string {
	parts := strings.Split(ref, "/")

	name := parts[3]
	if parts[2] == "messages" {
		name += "Message"
	}

	return Namify(name)
}

func ReferenceToStructAttributePath(ref string) string {
	ref = strings.Replace(ref, ".", "/", -1)
	ref = strings.Replace(ref, "#", "", -1)

	elems := strings.Split(ref, "/")[1:]
	for k, v := range elems {
		// If this is concerning the header, then it will be named "headers"
		if v == "header" {
			v = "headers"
		}

		elems[k] = Namify(v)
	}

	return strings.Join(elems, ".")
}

func HasField(v interface{}, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	return rv.FieldByName(name).IsValid()
}

func ChannelToMessageTypeName(ch asyncapi.Channel) string {
	msg := ch.GetChannelMessage()

	if msg.Payload != nil || msg.OneOf != nil {
		return Namify(ch.Name) + "Message"
	}

	return ReferenceToTypeName(msg.Reference)
}

func IsRequired(any asyncapi.Any, field string) bool {
	return any.IsFieldRequired(field)
}
