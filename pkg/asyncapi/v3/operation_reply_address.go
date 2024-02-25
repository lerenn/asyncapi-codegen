package asyncapiv3

// OperationReply is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v3.0.0#operationReplyAddressObject
type OperationReplyAddress struct {
	Description string `json:"description"`
	Location    string `json:"location"`
}
