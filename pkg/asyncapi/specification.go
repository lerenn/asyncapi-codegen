package asyncapi

// Specification only contains common functions between each version.
// This should be casted to get all other functions, base on the version.
type Specification interface {
	AsyncAPIVersion() string
}
