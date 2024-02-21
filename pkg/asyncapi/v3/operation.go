package asyncapiv3

// Operation is a representation of the corresponding asyncapi object filled
// from an asyncapi specification that will be used to generate code.
// Source: https://www.asyncapi.com/docs/reference/specification/v2.6.0#operationObject
type Operation struct {
	OperationID string  `json:"operationId"`
	Message     Message `json:"message"`
}
