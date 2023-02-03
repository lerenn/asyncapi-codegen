package asyncapi

type Message struct {
	Description   string         `json:"description"`
	Headers       *Any           `json:"headers"`
	Payload       *Any           `json:"payload"`
	CorrelationId *CorrelationId `json:"correlationId"`
	Reference     string         `json:"$ref"`

	// Non AsyncAPI fields
	Name string `json:"-"`
}
