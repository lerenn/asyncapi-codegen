//go:generate go run ../../../../cmd/asyncapi-codegen -p issue262 --force-pointers -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue262

import (
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
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

func (suite *Suite) TestToBrokerMessage() {
	t := TestMessage{
		Headers: HeadersFromTestMessage{
			// All fields should be pointers (nillable)
			FieldNonReq:  nil,
			FieldReq:     nil,
			SomeDateTime: nil,
		},
		Payload: "",
	}

	// toBrokerMessage should return an error due to FieldReq being required
	_, err := t.toBrokerMessage()
	assert.ErrorContains(suite.T(), err, "field FieldReq should not be nil")
	str := "Something"
	t.Headers.FieldReq = &str
	_, err = t.toBrokerMessage()
	assert.NoError(suite.T(), err)
}

func (suite *Suite) TestBrokerMessageToTestMessage() {
	// BrokerMessageToTestMessage should be valid
	_, err := brokerMessageToTestMessage(
		extensions.BrokerMessage{
			Headers: map[string][]byte{},
		},
	)
	// There is currently no check on fields returned by brokerMessageToTestMessage
	assert.NoError(suite.T(), err)
}
