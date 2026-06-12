//go:generate go run ../../../../cmd/asyncapi-codegen -p issue261 -g types -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue261

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/utils"
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

// TestAdditionalPropertiesAnyValue ensures that a schema using
// `additionalProperties: true` parses (it used to crash) and generates a
// map[string]any whose values keep their JSON type when marshaled
// (issue #261). Before the fix, non-string additional properties were
// serialized as quoted strings.
func (suite *Suite) TestAdditionalPropertiesAnyValue() {
	original := FreeFormSchema{
		Name: utils.ToPointer("example"),
		AdditionalProperties: map[string]any{
			"count": float64(42),
			"flag":  true,
		},
	}

	data, err := json.Marshal(original)
	suite.Require().NoError(err)

	// Non-string values must not be quoted.
	suite.Require().True(strings.Contains(string(data), `"count":42`), string(data))
	suite.Require().True(strings.Contains(string(data), `"flag":true`), string(data))

	// Round-trip back.
	var got FreeFormSchema
	suite.Require().NoError(json.Unmarshal(data, &got))
	suite.Require().Equal("example", *got.Name)
	suite.Require().Equal(float64(42), got.AdditionalProperties["count"])
	suite.Require().Equal(true, got.AdditionalProperties["flag"])
}
