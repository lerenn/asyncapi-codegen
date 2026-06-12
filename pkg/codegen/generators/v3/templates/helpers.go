package templates

import (
	"fmt"
	"regexp"
	"strconv"
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

// HasScalarDefault reports whether the schema declares a default value that can
// be rendered as a Go scalar literal (string, boolean, integer or number).
// Date/time strings, objects, arrays and other complex defaults are not
// supported and return false.
func HasScalarDefault(s *asyncapi.Schema) bool {
	if s == nil || s.Default == nil {
		return false
	}

	switch s.Type {
	case "boolean", "integer", "number":
		return true
	case "string":
		// Generated date/time types have a non-trivial literal form, so they
		// are intentionally left out.
		return s.Format != "date" && s.Format != "date-time"
	default:
		return false
	}
}

// DefaultLiteral returns the Go literal for the schema's default value, typed to
// match the field type produced by the "schema-name" template. It assumes
// HasScalarDefault returned true for the schema.
func DefaultLiteral(s *asyncapi.Schema) string {
	switch s.Type {
	case "boolean":
		b, _ := s.Default.(bool)
		return strconv.FormatBool(b)
	case "string":
		str, _ := s.Default.(string)
		return strconv.Quote(str)
	case "integer":
		goType := "int64"
		if s.Format == "int32" {
			goType = "int32"
		}
		return fmt.Sprintf("%s(%s)", goType, formatDefaultNumber(s.Default))
	case "number":
		goType := "float64"
		if s.Format == "float" {
			goType = "float32"
		}
		return fmt.Sprintf("%s(%s)", goType, formatDefaultNumber(s.Default))
	default:
		return ""
	}
}

// formatDefaultNumber renders a numeric default (decoded from JSON as float64,
// but also tolerating int/int64) without a spurious trailing ".0" for integers.
func formatDefaultNumber(v any) string {
	switch n := v.(type) {
	case float64:
		return strconv.FormatFloat(n, 'f', -1, 64)
	case int:
		return strconv.Itoa(n)
	case int64:
		return strconv.FormatInt(n, 10)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// SubscribeFuncName returns the Go name of the SubscribeTo function generated
// for the given operation, honoring the x-go-subscribe-func override.
func SubscribeFuncName(op *asyncapi.Operation) string {
	if name := op.Follow().ExtGoSubscribeFunc; name != "" {
		return templateutil.Namify(name)
	}
	return "SubscribeTo" + templateutil.Namify(op.Follow().Name)
}

// ReceivedFuncName returns the Go name of the subscriber callback generated for
// the given operation, honoring the x-go-received-func override.
func ReceivedFuncName(op *asyncapi.Operation) string {
	if name := op.Follow().ExtGoReceivedFunc; name != "" {
		return templateutil.Namify(name)
	}
	return templateutil.Namify(op.Follow().Name) + "Received"
}

// ReplyFuncName returns the Go name of the ReplyTo function generated for the
// given operation, honoring the x-go-reply-func override.
func ReplyFuncName(op *asyncapi.Operation) string {
	if name := op.Follow().ExtGoReplyFunc; name != "" {
		return templateutil.Namify(name)
	}
	return "ReplyTo" + templateutil.Namify(op.Follow().Name)
}

// SendFuncName returns the Go name of the Send function generated for the given
// operation, honoring the x-go-send-func override. The prefix is the controller
// side ("App" or "User"), which selects the default "SendAs"/"SendTo" verb.
func SendFuncName(op *asyncapi.Operation, prefix string) string {
	if name := op.Follow().ExtGoSendFunc; name != "" {
		return templateutil.Namify(name)
	}
	verb := "SendAs"
	if prefix == "User" {
		verb = "SendTo"
	}
	return verb + templateutil.Namify(op.Follow().Name)
}

// RequestFuncName returns the Go name of the Request function generated for the
// given operation, honoring the x-go-request-func override. The prefix is the
// controller side ("App" or "User"), which selects the default
// "RequestAs"/"RequestTo" verb.
func RequestFuncName(op *asyncapi.Operation, prefix string) string {
	if name := op.Follow().ExtGoRequestFunc; name != "" {
		return templateutil.Namify(name)
	}
	verb := "RequestAs"
	if prefix == "User" {
		verb = "RequestTo"
	}
	return verb + templateutil.Namify(op.Follow().Name)
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
		"hasScalarDefault":               HasScalarDefault,
		"defaultLiteral":                 DefaultLiteral,
		"subscribeFuncName":              SubscribeFuncName,
		"receivedFuncName":               ReceivedFuncName,
		"replyFuncName":                  ReplyFuncName,
		"sendFuncName":                   SendFuncName,
		"requestFuncName":                RequestFuncName,
	}
}
