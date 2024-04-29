package asyncapi

// Specification only contains common functions between each version.
// This should be casted to get all other functions, base on the version.
type Specification interface {
	// MajorVersion returns the major version of the AsyncAPI specification.
	MajorVersion() int
	// Process processes all information in specification in order to link the
	// references, apply the traits and generating other information for code
	// generation.
	//
	// WARNING: this will alter the specification as you will find, by example,
	// traits applied in the specification.
	Process() error
	// AddDependency adds a dependency to the specification.
	AddDependency(path string, spec Specification) error
}
