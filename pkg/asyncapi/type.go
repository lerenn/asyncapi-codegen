package asyncapi

// Type is a structure that represents the type of a field.
type Type string

// String returns the string representation of the type.
func (t Type) String() string {
	return string(t)
}

const (
	// TypeIsArray represents the type of an array.
	TypeIsArray Type = "array"
	// TypeIsHeader represents the type of a header.
	TypeIsHeader Type = "header"
	// TypeIsObject represents the type of an object.
	TypeIsObject Type = "object"
	// TypeIsString represents the type of a string.
	TypeIsString Type = "string"
	// TypeIsInteger represents the type of an integer.
	TypeIsInteger Type = "integer"
	// TypeIsPayload represents the type of a payload.
	TypeIsPayload Type = "payload"
)
