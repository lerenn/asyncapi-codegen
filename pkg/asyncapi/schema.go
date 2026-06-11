package asyncapi

// Validations is a representation of the JSON-Object validation fields supported by asyncapi
// These fields are in common for v2 and v3.
type Validations[T any] struct {
	Required         []string `json:"required,omitempty"`
	MultipleOf       []string `json:"multipleOf,omitempty"`
	Maximum          float64  `json:"maximum,omitempty"`
	ExclusiveMaximum float64  `json:"exclusiveMaximum,omitempty"`
	Minimum          float64  `json:"minimum,omitempty"`
	ExclusiveMinimum float64  `json:"exclusiveMinimum,omitempty"`
	MaxLength        uint     `json:"maxLength,omitempty"`
	MinLength        uint     `json:"minLength,omitempty"`
	Pattern          string   `json:"pattern,omitempty"`
	MaxItems         uint     `json:"maxItems,omitempty"`
	MinItems         uint     `json:"minItems,omitempty"`
	UniqueItems      bool     `json:"uniqueItems,omitempty"`
	MaxProperties    uint     `json:"maxProperties,omitempty"`
	MinProperties    uint     `json:"minProperties,omitempty"`
	Enum             []any    `json:"enum,omitempty"`
	Const            any      `json:"const,omitempty"`

	AllOf []*T `json:"allOf,omitempty"`
	AnyOf []*T `json:"anyOf,omitempty"`
	OneOf []*T `json:"oneOf,omitempty"`

	// --- Non JSON Schema/AsyncAPI fields -------------------------------------
	IsRequired      bool  `json:"-"`
	ShouldOmitEmpty *bool `json:"-"`
}

// Merge merges the newV into the current Validations.
//
//nolint:cyclop,funlen // This function is a merge function and it is expected to have a high cyclomatic complexity.
func (v *Validations[T]) Merge(newV Validations[T]) {
	if len(newV.Required) > 0 {
		v.Required = newV.Required
	}
	if len(newV.MultipleOf) > 0 {
		v.MultipleOf = newV.MultipleOf
	}
	if newV.Maximum != 0 {
		v.Maximum = newV.Maximum
	}
	if newV.ExclusiveMaximum != 0 {
		v.ExclusiveMaximum = newV.ExclusiveMaximum
	}
	if newV.Minimum != 0 {
		v.Minimum = newV.Minimum
	}
	if newV.ExclusiveMinimum != 0 {
		v.ExclusiveMinimum = newV.ExclusiveMinimum
	}
	if newV.MaxLength != 0 {
		v.MaxLength = newV.MaxLength
	}
	if newV.MinLength != 0 {
		v.MinLength = newV.MinLength
	}
	if newV.Pattern != "" {
		v.Pattern = newV.Pattern
	}
	if newV.MaxItems != 0 {
		v.MaxItems = newV.MaxItems
	}
	if newV.MinItems != 0 {
		v.MinItems = newV.MinItems
	}
	if newV.UniqueItems {
		v.UniqueItems = newV.UniqueItems
	}
	if newV.MaxProperties != 0 {
		v.MaxProperties = newV.MaxProperties
	}
	if newV.MinProperties != 0 {
		v.MinProperties = newV.MinProperties
	}
	if len(newV.Enum) > 0 {
		v.Enum = newV.Enum
	}
	if newV.Const != nil {
		v.Const = newV.Const
	}
	if len(newV.AllOf) > 0 {
		v.AllOf = newV.AllOf
	}
	if len(newV.AnyOf) > 0 {
		v.AnyOf = newV.AnyOf
	}
	if len(newV.OneOf) > 0 {
		v.OneOf = newV.OneOf
	}
	if newV.IsRequired {
		v.IsRequired = newV.IsRequired
	}
	if newV.ShouldOmitEmpty != nil {
		v.ShouldOmitEmpty = newV.ShouldOmitEmpty
	}
}
