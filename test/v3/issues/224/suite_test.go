//go:generate go run ../../../../cmd/asyncapi-codegen -p issue224 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue224

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, NewSuite())
}

type Suite struct {
	suite.Suite
}

func NewSuite() *Suite {
	return &Suite{}
}

func (suite *Suite) TestMarshalAndUnmarshalRoundtrip() {
	f := float32(66.66)
	// Create an original schema with additional properties
	original := ColliderDictionarySchema{
		AdditionalProperties: map[string]ColliderSchema{
			"testCollider": {
				Margin: &f,
				Pose: &PoseSchema{
					Position:    &Vector3dSchema{1.1, 2.2, 3.3},
					Orientation: &Vector3dSchema{4.4, 5.5, 6.6, 7.7},
				},
				Shape: ShapePropertyFromColliderSchema{
					Radius:    10.0,
					ShapeType: "sphere",
				},
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	suite.Require().NoError(err)

	// Unmarshal back to a new instance
	var unmarshaled ColliderDictionarySchema
	err = json.Unmarshal(data, &unmarshaled)
	suite.Require().NoError(err)

	// Check that the additional properties are preserved
	suite.Equal(original, unmarshaled)
}
