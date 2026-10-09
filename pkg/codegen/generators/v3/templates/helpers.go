package templates

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
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
		if IsGeneratedEnum(s) {
			return fmt.Sprintf("%s(%s)", templateutil.Namify(s.Name), strconv.Quote(str))
		}
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

// IsGeneratedEnum reports whether the schema is a string enumeration for which a
// dedicated Go type and constants are generated (issue #337). Enumerations with
// a custom Go type, a generated date/time format or non-string values are left
// as their plain underlying type.
func IsGeneratedEnum(s *asyncapi.Schema) bool {
	if s == nil || s.Type != "string" || len(s.Enum) == 0 || s.Name == "" || s.ExtGoType != "" {
		return false
	}
	if s.Format == "date" || s.Format == "date-time" {
		return false
	}
	for _, e := range s.Enum {
		if _, ok := e.(string); !ok {
			return false
		}
	}
	return true
}

// EnumChildren returns the generated enumeration schemas that are direct
// children (properties, items, additional properties) of the given schema,
// sorted by name so the generated code is stable.
func EnumChildren(s *asyncapi.Schema) []*asyncapi.Schema {
	// Each candidate is paired with the schema that names it. A schema can be
	// shared between several parents (e.g. through a merged reference), so an
	// enumeration is only generated by the parent it was named after, to avoid
	// declaring the same type twice.
	type candidate struct{ owner, schema *asyncapi.Schema }

	candidates := []candidate{{s, s.Items}, {s, s.AdditionalProperties}}
	for _, p := range s.Properties {
		candidates = append(candidates, candidate{s, p})
		if p.Type == "array" {
			// Arrays of enumerations declared inline in a property.
			candidates = append(candidates, candidate{p, p.Items})
		}
	}

	children := make([]*asyncapi.Schema, 0, len(candidates))
	seen := make(map[string]bool)
	for _, c := range candidates {
		if IsGeneratedEnum(c.schema) && strings.HasSuffix(c.schema.Name, "From"+c.owner.Name) && !seen[c.schema.Name] {
			seen[c.schema.Name] = true
			children = append(children, c.schema)
		}
	}
	sort.Slice(children, func(i, j int) bool { return children[i].Name < children[j].Name })

	return children
}

// EnumConstants returns the Go constant names of the enumeration values, in the
// order of the specification. Names are the type name followed by the namified
// value; empty or colliding names get the value index appended.
func EnumConstants(s *asyncapi.Schema) []string {
	typeName := templateutil.Namify(s.Name)
	names := make([]string, len(s.Enum))
	used := make(map[string]bool, len(s.Enum))
	for i, e := range s.Enum {
		v, _ := e.(string)
		name := typeName + templateutil.Namify(v)
		if templateutil.Namify(v) == "" || used[name] {
			name = fmt.Sprintf("%s%d", name, i)
		}
		used[name] = true
		names[i] = name
	}

	return names
}

// EnumValue returns the quoted Go string literal of the i-th enumeration value.
func EnumValue(s *asyncapi.Schema, i int) string {
	v, _ := s.Enum[i].(string)
	return strconv.Quote(v)
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

// OperationMessageCount returns the number of messages of an operation.
func OperationMessageCount(op *asyncapi.Operation) int {
	return len(op.Follow().GetMessages())
}

// MessageTypeName returns the generated Go type name of a message.
func MessageTypeName(msg *asyncapi.Message) string {
	return templateutil.Namify(msg.Follow().Name)
}

// SendFuncNameForMessage returns the name of the per-message send function
// generated when an operation carries more than one message (issue #140). It is
// the operation send function name suffixed with "For" and the message name.
func SendFuncNameForMessage(op *asyncapi.Operation, prefix string, msg *asyncapi.Message) string {
	base := strings.TrimSuffix(MessageTypeName(msg), "Message")
	return SendFuncName(op, prefix) + "For" + base
}

// messageNameBase returns the message type name without its "Message" suffix,
// used to build per-message function names.
func messageNameBase(msg *asyncapi.Message) string {
	return strings.TrimSuffix(MessageTypeName(msg), "Message")
}

// SubscribeFuncNameForMessage returns the name of the per-message subscribe
// function generated when a receive operation carries more than one message
// (issue #333). It is the operation subscribe function name suffixed with "For"
// and the message name.
func SubscribeFuncNameForMessage(op *asyncapi.Operation, msg *asyncapi.Message) string {
	return SubscribeFuncName(op) + "For" + messageNameBase(msg)
}

// UnsubscribeFuncNameForMessage returns the name of the per-message unsubscribe
// function generated when a receive operation carries more than one message
// (issue #333).
func UnsubscribeFuncNameForMessage(op *asyncapi.Operation, msg *asyncapi.Message) string {
	return "UnsubscribeFrom" + templateutil.Namify(op.Follow().Name) + "For" + messageNameBase(msg)
}

// ReceivedFuncNameForMessage returns the name of the per-message subscriber
// callback generated when a receive operation carries more than one message
// (issue #333).
func ReceivedFuncNameForMessage(op *asyncapi.Operation, msg *asyncapi.Message) string {
	return templateutil.Namify(op.Follow().Name) + "For" + messageNameBase(msg) + "Received"
}

// IsMultiMessageReceive reports whether the given receive operation carries more
// than one message and therefore needs per-message subscribe functions and
// try-each-until-valid dispatch (issue #333). Reply and reply-target operations
// are excluded as the request/reply path is handled separately.
func IsMultiMessageReceive(op *asyncapi.Operation) bool {
	o := op.Follow()
	return len(o.GetMessages()) > 1 && o.Reply == nil && o.ReplyOf == nil
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
		"isGeneratedEnum":                IsGeneratedEnum,
		"enumChildren":                   EnumChildren,
		"enumConstants":                  EnumConstants,
		"enumValue":                      EnumValue,
		"hasScalarDefault":               HasScalarDefault,
		"defaultLiteral":                 DefaultLiteral,
		"subscribeFuncName":              SubscribeFuncName,
		"receivedFuncName":               ReceivedFuncName,
		"replyFuncName":                  ReplyFuncName,
		"sendFuncName":                   SendFuncName,
		"requestFuncName":                RequestFuncName,
		"operationMessageCount":          OperationMessageCount,
		"messageTypeName":                MessageTypeName,
		"sendFuncNameForMessage":         SendFuncNameForMessage,
		"subscribeFuncNameForMessage":    SubscribeFuncNameForMessage,
		"unsubscribeFuncNameForMessage":  UnsubscribeFuncNameForMessage,
		"receivedFuncNameForMessage":     ReceivedFuncNameForMessage,
		"isMultiMessageReceive":          IsMultiMessageReceive,
	}
}
