//go:generate go run ../../../../cmd/asyncapi-codegen -p issue258 -g types -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue258

import (
	"encoding/json"
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

// TestOneOfHasOneFieldPerMember ensures a multi-member oneOf payload generates
// a struct with one optional field per member instead of flattening every
// member's properties into a single struct (issue #258).
func (suite *Suite) TestOneOfHasOneFieldPerMember() {
	// Compile-time proof both members exist as distinct fields.
	payload := TestMessageMessagePayload{
		IntValueSchema: &IntValueSchema{Value: utils.ToPointer(int64(42))},
	}
	suite.Require().NotNil(payload.IntValueSchema)
	suite.Require().Nil(payload.StrValueSchema)
}

// TestOneOfMarshalsTheSetMember ensures only the member that is set is emitted
// on the wire.
func (suite *Suite) TestOneOfMarshalsTheSetMember() {
	payload := TestMessageMessagePayload{
		StrValueSchema: &StrValueSchema{Value: utils.ToPointer("hello")},
	}

	data, err := json.Marshal(payload)
	suite.Require().NoError(err)
	suite.Require().JSONEq(`{"value":"hello"}`, string(data))
}

// TestOneOfUnmarshalsIntoMatchingMember ensures JSON is decoded into the member
// whose shape it matches. The two members use the same key with disjoint types
// (integer vs string), so decoding is unambiguous.
func (suite *Suite) TestOneOfUnmarshalsIntoMatchingMember() {
	var asInt TestMessageMessagePayload
	suite.Require().NoError(json.Unmarshal([]byte(`{"value":7}`), &asInt))
	suite.Require().NotNil(asInt.IntValueSchema)
	suite.Require().Equal(int64(7), *asInt.IntValueSchema.Value)
	suite.Require().Nil(asInt.StrValueSchema)

	var asStr TestMessageMessagePayload
	suite.Require().NoError(json.Unmarshal([]byte(`{"value":"seven"}`), &asStr))
	suite.Require().NotNil(asStr.StrValueSchema)
	suite.Require().Equal("seven", *asStr.StrValueSchema.Value)
	suite.Require().Nil(asStr.IntValueSchema)
}
