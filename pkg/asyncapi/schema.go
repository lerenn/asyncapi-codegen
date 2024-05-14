package asyncapi

// Validations is a representation of the JSON-Object validation fields supported by asyncapi
// These fields are in common for v2 and v3.
type Validations[T any] struct {
	Required         []string `json:"required"`
	MultipleOf       []string `json:"multipleOf"`
	Maximum          float64  `json:"maximum"`
	ExclusiveMaximum float64  `json:"exclusiveMaximum"`
	Minimum          float64  `json:"minimum"`
	ExclusiveMinimum float64  `json:"exclusiveMinimum"`
	MaxLength        uint     `json:"maxLength"`
	MinLength        uint     `json:"minLength"`
	Pattern          string   `json:"pattern"`
	MaxItems         uint     `json:"maxItems"`
	MinItems         uint     `json:"minItems"`
	UniqueItems      bool     `json:"uniqueItems"`
	MaxProperties    uint     `json:"maxProperties"`
	MinProperties    uint     `json:"minProperties"`
	Enum             []any    `json:"enum"`
	Const            any      `json:"const"`

	AllOf []*T `json:"allOf"`
	AnyOf []*T `json:"anyOf"`
	OneOf []*T `json:"oneOf"`

	// --- Non JSON Schema/AsyncAPI fields -------------------------------------
	IsRequired bool `json:"-"`
}
