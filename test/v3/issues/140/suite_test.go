package issue140

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

// capturingBroker records the channel address and payload of the last published
// message.
type capturingBroker struct {
	channel string
	payload []byte
}

func (b *capturingBroker) Publish(_ context.Context, channel string, msg extensions.BrokerMessage) error {
	b.channel = channel
	b.payload = msg.Payload
	return nil
}

func (b *capturingBroker) Subscribe(
	_ context.Context,
	_ string,
) (extensions.BrokerChannelSubscription, error) {
	return extensions.NewBrokerChannelSubscription(
		make(chan extensions.AcknowledgeableBrokerMessage),
		make(chan any),
	), nil
}

type Suite struct {
	suite.Suite
}

func NewSuite() *Suite {
	return &Suite{}
}

// TestEachMessageHasItsOwnSendFunction ensures an operation that carries more
// than one message generates one send function per message, each sending its
// own payload on the channel (issue #140).
func (suite *Suite) TestEachMessageHasItsOwnSendFunction() {
	broker := &capturingBroker{}
	ctrl, err := NewAppController(broker)
	suite.Require().NoError(err)
	defer ctrl.Close(context.Background())

	created := NewUserCreatedMessage()
	created.Payload.Id = utils.ToPointer("created-1")
	suite.Require().NoError(ctrl.SendAsSendEventsOperationForUserCreated(context.Background(), created))
	suite.Require().Equal("v3.issue140.events", broker.channel)
	suite.Require().Contains(string(broker.payload), "created-1")

	deleted := NewUserDeletedMessage()
	deleted.Payload.Id = utils.ToPointer("deleted-1")
	suite.Require().NoError(ctrl.SendAsSendEventsOperationForUserDeleted(context.Background(), deleted))
	suite.Require().Equal("v3.issue140.events", broker.channel)
	suite.Require().Contains(string(broker.payload), "deleted-1")
}
