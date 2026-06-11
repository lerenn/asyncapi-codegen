//go:generate go run ../../../../cmd/asyncapi-codegen -p issue290 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue290

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func stringPtr(s string) *string {
	return &s
}

func (suite *Suite) TestMessagePayloadTypeGeneration() {
	// Test that the message payload type is properly generated from allOf schema reference
	// Before the fix, this would fail to compile because Payload had no type
	msg := NewTestCreatedMessage()

	// Verify the message has a properly typed payload field from allOf composition
	msg.Payload.Id = stringPtr("test-id")
	msg.Payload.Timestamp = stringPtr("2023-01-01T00:00:00Z")
	msg.Payload.Data = stringPtr("test-data")

	// Verify fields are accessible and properly typed
	assert.Equal(suite.T(), "test-id", *msg.Payload.Id)
	assert.Equal(suite.T(), "2023-01-01T00:00:00Z", *msg.Payload.Timestamp)
	assert.Equal(suite.T(), "test-data", *msg.Payload.Data)
}
