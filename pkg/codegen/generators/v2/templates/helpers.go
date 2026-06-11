package templates

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"

	asyncapi "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v2"
	"github.com/lerenn/asyncapi-codegen/pkg/codegen/generators"
	templateutil "github.com/lerenn/asyncapi-codegen/pkg/utils/template"
)

// ReferenceToTypeName will convert a reference to a type name in the form of
// golang conventional type names.
func ReferenceToTypeName(ref string) (string, error) {
	parts := strings.Split(ref, "/")
	if len(parts) < 4 {
		return "", fmt.Errorf("invalid reference %q: expected at least 4 '/'-separated segments", ref)
	}
	return templateutil.Namify(parts[3]), nil
}

// ChannelToMessage will convert a channel to its message, based on publish/subscribe.
//
//nolint:cyclop
func ChannelToMessage(ch asyncapi.Channel, direction string) (*asyncapi.Message, error) {
	switch {
	case ch.Publish != nil && ch.Subscribe == nil:
		return ch.Publish.Message.Follow(), nil
	case ch.Subscribe != nil && ch.Publish == nil:
		return ch.Subscribe.Message.Follow(), nil
	case direction == "publish":
		if ch.Publish == nil {
			return nil, fmt.Errorf("channel %q has no publish operation", ch.Name)
		}
		return ch.Publish.Message.Follow(), nil
	case direction == "subscribe":
		if ch.Subscribe == nil {
			return nil, fmt.Errorf("channel %q has no subscribe operation", ch.Name)
		}
		return ch.Subscribe.Message.Follow(), nil
	case ch.Subscribe == nil && ch.Publish == nil:
		return nil, fmt.Errorf("channel %q has no publish or subscribe operation", ch.Name)
	default:
		return nil, fmt.Errorf("direction must be either 'publish' or 'subscribe', got %q", direction)
	}
}

// GenerateChannelPath will generate a channel path with the given channel.
func GenerateChannelPath(ch asyncapi.Channel) string {
	// If there is no parameter, then just return the path
	if ch.Parameters == nil {
		return fmt.Sprintf("%q", ch.Path)
	}

	parameterRegexp := regexp.MustCompile("{[^{}]*}")

	matches := parameterRegexp.FindAllString(ch.Path, -1)
	if len(matches) == 0 {
		return fmt.Sprintf("%q", ch.Path)
	}
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
		"getChildrenObjectSchemas":       generators.GetChildrenObjectSchemas[asyncapi.Schema],
		"channelToMessage":               ChannelToMessage,
		"isRequired":                     generators.IsRequired[asyncapi.Schema],
		"isFieldPointer":                 generators.IsFieldPointer[asyncapi.Schema],
		"generateChannelPath":            GenerateChannelPath,
		"referenceToStructAttributePath": generators.ReferenceToStructAttributePath,
		"operationName":                  OperationName,
		"referenceToTypeName":            ReferenceToTypeName,
		"generateValidateTags":           generators.GenerateValidateTags[asyncapi.Schema],
		"generateJSONTags":               generators.GenerateJSONTags[asyncapi.Schema],
	}
}
