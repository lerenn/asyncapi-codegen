package asyncapi

import (
	"fmt"

	asyncapiv2 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v2"
	asyncapiv3 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
)

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
	Process()
}

// ToV2 returns an AsyncAPI specification V2 from interface, if compatible.
// Note: Before using this, you should make sure that parsed data is in version 2.
func ToV2(s Specification) (*asyncapiv2.Specification, error) {
	spec, ok := s.(*asyncapiv2.Specification)
	if !ok {
		return nil, fmt.Errorf("unknown spec format: should have been a v2 format")
	}

	return spec, nil
}

// ToV3 returns an AsyncAPI specification V3 from interface, if compatible.
// Note: Before using this, you should make sure that parsed data is in version 3.
func ToV3(s Specification) (*asyncapiv3.Specification, error) {
	spec, ok := s.(*asyncapiv3.Specification)
	if !ok {
		return nil, fmt.Errorf("unknown spec format: should have been a v3 format")
	}

	return spec, nil
}
