//go:generate go run ../../../../cmd/asyncapi-codegen -p issue139 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue139

import (
	"context"
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/utils"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, NewSuite())
}

// capturingBroker records the channel address of the last published message.
type capturingBroker struct {
	publishedChannel string
}

func (b *capturingBroker) Publish(_ context.Context, channel string, _ extensions.BrokerMessage) error {
	b.publishedChannel = channel
	return nil
}

func (b *capturingBroker) Subscribe(
	_ context.Context,
	_ string,
) (extensions.BrokerChannelSubscription, error) {
	return extensions.NewBrokerChannelSubscription(make(chan extensions.AcknowledgeableBrokerMessage), make(chan any)), nil
}

type Suite struct {
	suite.Suite
}

func NewSuite() *Suite {
	return &Suite{}
}

// TestParameterFilledFromMessageLocation ensures a channel parameter that
// declares a `location` is auto-filled from the outgoing message when the
// caller leaves it empty (issue #139).
func (suite *Suite) TestParameterFilledFromMessageLocation() {
	broker := &capturingBroker{}
	ctrl, err := NewAppController(broker)
	suite.Require().NoError(err)
	defer ctrl.Close(context.Background())

	msg := NewUserSignedUpMessage()
	msg.Payload.UserId = utils.ToPointer("abc")

	// userId left empty in params: it must be filled from the message payload.
	err = ctrl.SendAsSendUserOperation(context.Background(), UserChannelParameters{}, msg)
	suite.Require().NoError(err)
	suite.Require().Equal("v3.issue139.user.abc", broker.publishedChannel)
}

// TestParameterOverridesMessageLocation ensures an explicitly provided
// parameter is not overwritten by the message location value.
func (suite *Suite) TestParameterOverridesMessageLocation() {
	broker := &capturingBroker{}
	ctrl, err := NewAppController(broker)
	suite.Require().NoError(err)
	defer ctrl.Close(context.Background())

	msg := NewUserSignedUpMessage()
	msg.Payload.UserId = utils.ToPointer("abc")

	err = ctrl.SendAsSendUserOperation(context.Background(), UserChannelParameters{UserId: "override"}, msg)
	suite.Require().NoError(err)
	suite.Require().Equal("v3.issue139.user.override", broker.publishedChannel)
}
