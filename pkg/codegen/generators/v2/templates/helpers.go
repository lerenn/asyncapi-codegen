package templates

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"

	asyncapi "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v2"
	templateutil "github.com/lerenn/asyncapi-codegen/pkg/utils/template"
)

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

		path[k] = templateutil.Namify(v)
	}

	return strings.Join(path, ".")
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

	return templateutil.Namify(name)
}

// ChannelToMessageTypeName will convert a channel to a message type name in the
// form of golang conventional type names.
func ChannelToMessageTypeName(ch asyncapi.Channel) string {
	msg := ch.GetChannelMessage()

	if msg.Payload != nil || msg.OneOf != nil {
		return templateutil.Namify(ch.Name) + "Message"
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
		sprint += fmt.Sprintf("params.%s,", templateutil.Namify(m))
	}

	return sprint[:len(sprint)-1] + ")"
}

// OperationName returns `operationId` value from Publish or Subscribe operation if any.
// If no `operationID` exists — return provided default value (`name`).
func OperationName(channel asyncapi.Channel) string {
	var name string

	switch {
	case channel.Publish != nil && channel.Publish.OperationID != "":
		name = channel.Publish.OperationID
	case channel.Subscribe != nil && channel.Subscribe.OperationID != "":
		name = channel.Subscribe.OperationID
	default:
		name = channel.Name
	}

	return templateutil.Namify(name)
}

// HelpersFunctions returns the functions that can be used as helpers
// in a golang template.
func HelpersFunctions() template.FuncMap {
	return template.FuncMap{
		"channelToMessageTypeName":       ChannelToMessageTypeName,
		"isRequired":                     IsRequired,
		"generateChannelPath":            GenerateChannelPath,
		"referenceToStructAttributePath": ReferenceToStructAttributePath,
		"operationName":                  OperationName,
		"referenceToTypeName":            ReferenceToTypeName,
	}
}
