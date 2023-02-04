package generators

import (
	"embed"
	"fmt"
	"path"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
	"github.com/stoewer/go-strcase"
)

const (
	templatesDir = "templates"

	importsTemplatePath = templatesDir + "/imports.tmpl"
	typesTemplatePath   = templatesDir + "/types.tmpl"

	appDir                    = templatesDir + "/app"
	appControllerTemplatePath = appDir + "/controller.tmpl"
	appSubscriberTemplatePath = appDir + "/subscriber.tmpl"

	clientDir                    = templatesDir + "/client"
	clientControllerTemplatePath = clientDir + "/controller.tmpl"
	clientSubscriberTemplatePath = clientDir + "/subscriber.tmpl"

	brokerDir                    = templatesDir + "/broker"
	brokerControllerTemplatePath = brokerDir + "/controller.tmpl"
	brokerNATSTemplatePath       = brokerDir + "/nats.tmpl"

	elementsDir         = templatesDir + "/elements"
	anyTemplatePath     = elementsDir + "/any.tmpl"
	messageTemplatePath = elementsDir + "/message.tmpl"
)

var (
	//go:embed templates/*
	files embed.FS
)

func loadTemplate(paths ...string) (*template.Template, error) {
	return template.
		New(path.Base(paths[0])).
		Funcs(templateFunctions()).
		ParseFS(files, paths...)
}

func templateFunctions() template.FuncMap {
	return template.FuncMap{
		"namify":                         namify,
		"snakeCase":                      snakeCase,
		"referenceToStructAttributePath": referenceToStructAttributePath,
		"referenceToTypeName":            referenceToTypeName,
		"channelToMessageTypeName":       channelToMessageTypeName,
		"hasField":                       hasField,
	}
}

func namify(sentence string) string {
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
		sentence = string(re.ReplaceAll([]byte(sentence), []byte(a)))
	}

	return sentence
}

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func snakeCase(sentence string) string {
	snake := matchFirstCap.ReplaceAllString(sentence, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func referenceToTypeName(ref string) string {
	parts := strings.Split(ref, "/")

	name := parts[3]
	if parts[2] == "messages" {
		name += "Message"
	}

	return namify(name)
}

func referenceToStructAttributePath(ref string) string {
	ref = strings.Replace(ref, ".", "/", -1)
	ref = strings.Replace(ref, "#", "", -1)

	elems := strings.Split(ref, "/")[1:]
	for k, v := range elems {
		// If this is  concerning the header, then it will be named "headers"
		if v == "header" {
			v = "headers"
		}

		elems[k] = namify(v)
	}

	return strings.Join(elems, ".")
}

func hasField(v interface{}, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	return rv.FieldByName(name).IsValid()
}

func channelToMessageTypeName(ch asyncapi.Channel) string {
	msg := ch.GetMessageWithoutReferenceRedirect()

	if msg.Payload != nil {
		return namify(ch.Name) + "Message"
	}

	return referenceToTypeName(msg.Reference)
}
