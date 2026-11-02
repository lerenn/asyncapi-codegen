//go:generate go run ../../../../cmd/asyncapi-codegen -p issue341 -i asyncapi.yaml -o ./asyncapi.gen.go
package issue341

import (
	"context"
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

// An operation message that references a channel message, which in turn references a component
// message, used to resolve to the channel message only. The generated controller code then used
// a type name that is never defined (in this example, EventMessageFromEventsChannel), so the
// package did not compile.
// The ability for this package to compile asserts that the issue is no longer present.
func (suite *Suite) TestOperationMessageResolvesToComponentMessage() {
	var _ func(*AppController, context.Context, EventMessage) error = (*AppController).SendAsSendEventOperation
	var _ func(AppSubscriber, context.Context, EventMessage) error = AppSubscriber.ReceiveEventOperationReceived

	msg := NewEventMessage()
	id := "1"
	msg.Payload.Id = &id
	suite.Require().Equal("1", *msg.Payload.Id)
}
