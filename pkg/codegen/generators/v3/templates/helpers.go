package templates

import (
	"fmt"
	"regexp"
	"text/template"

	asyncapi "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
	"github.com/lerenn/asyncapi-codegen/pkg/codegen/generators"
	templateutil "github.com/lerenn/asyncapi-codegen/pkg/utils/template"
)

// ChannelToMessageTypeName will convert a channel to a message type name in the
// form of golang conventional type names.
func ChannelToMessageTypeName(ch asyncapi.Channel) (string, error) {
	msg, err := ch.Follow().GetMessage()
	if err != nil {
		return "", err
	}
	return templateutil.Namify(msg.Follow().Name), nil
}

// OpToMsgTypeName will convert an operation to a message type name in the
// form of golang conventional type names.
func OpToMsgTypeName(op asyncapi.Operation) (string, error) {
	msg, err := op.Follow().GetMessage()
	if err != nil {
		return "", err
	}
	return templateutil.Namify(msg.Follow().Name), nil
}

// OpToChannelTypeName will convert an operation to a channel type name in the
// form of golang conventional type names.
func OpToChannelTypeName(op asyncapi.Operation) string {
	ch := op.Channel.Follow()
	return templateutil.Namify(ch.Name)
}

// GenerateChannelAddrFromOp will generate a channel path with the given operation.
func GenerateChannelAddrFromOp(op asyncapi.Operation) string {
	ch := op.Channel.Follow()
	return GenerateChannelAddr(ch)
}

// GenerateChannelAddr will generate a channel path with the given channel.
func GenerateChannelAddr(ch *asyncapi.Channel) string {
	// Be sure this is the final channel, not a proxy
	ch = ch.Follow()

	// If there is no parameter, then just return the path
	if ch.Parameters == nil {
		return fmt.Sprintf("%q", ch.Address)
	}

	parameterRegexp := regexp.MustCompile("{[^{}]*}")

	matches := parameterRegexp.FindAllString(ch.Address, -1)
	if len(matches) == 0 {
		return fmt.Sprintf("%q", ch.Address)
	}
	format := parameterRegexp.ReplaceAllString(ch.Address, "%s")

	sprint := fmt.Sprintf("fmt.Sprintf(%q, ", format)
	for _, m := range matches {
		sprint += fmt.Sprintf("params.%s,", templateutil.Namify(m))
	}

	return sprint[:len(sprint)-1] + ")"
}

// HelpersFunctions returns the functions that can be used as helpers
// in a golang template.
func HelpersFunctions() template.FuncMap {
	return template.FuncMap{
		"getChildrenObjectSchemas":       generators.GetChildrenObjectSchemas[asyncapi.Schema],
		"channelToMessageTypeName":       ChannelToMessageTypeName,
		"opToMsgTypeName":                OpToMsgTypeName,
		"opToChannelTypeName":            OpToChannelTypeName,
		"isRequired":                     generators.IsRequired[asyncapi.Schema],
		"isFieldPointer":                 generators.IsFieldPointer[asyncapi.Schema],
		"generateChannelAddr":            GenerateChannelAddr,
		"generateChannelAddrFromOp":      GenerateChannelAddrFromOp,
		"referenceToStructAttributePath": generators.ReferenceToStructAttributePath,
		"generateValidateTags":           generators.GenerateValidateTags[asyncapi.Schema],
		"generateJSONTags":               generators.GenerateJSONTags[asyncapi.Schema],
	}
}
