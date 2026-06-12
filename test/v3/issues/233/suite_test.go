//go:generate go run ../../../../cmd/asyncapi-codegen -p issue233 -g types -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue233

import (
	"strings"
	"testing"
	"time"

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

// TestDateTimePreservesFractionalSeconds ensures a date-time payload keeps its
// sub-second precision when marshaled to the broker. Before the fix the
// generated code used time.RFC3339, which drops fractional seconds on output
// (issue #233); it now uses time.RFC3339Nano.
func (suite *Suite) TestDateTimePreservesFractionalSeconds() {
	original := time.Date(2024, 1, 2, 3, 4, 5, 123456789, time.UTC)

	// Marshal to the broker message.
	msg := TestMessageMessage{Payload: original}
	bMsg, err := msg.toBrokerMessage()
	suite.Require().NoError(err)

	// The wire representation must carry the fractional seconds.
	suite.Require().True(
		strings.Contains(string(bMsg.Payload), ".123456789"),
		"expected fractional seconds in payload, got %q", string(bMsg.Payload))

	// Round-trip back and ensure the nanoseconds survived.
	got, err := brokerMessageToTestMessageMessage(bMsg)
	suite.Require().NoError(err)
	suite.Require().True(original.Equal(got.Payload))
	suite.Require().Equal(123456789, got.Payload.Nanosecond())
}
