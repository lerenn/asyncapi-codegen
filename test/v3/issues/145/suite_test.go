//go:generate go run ../../../../cmd/asyncapi-codegen -p issue145 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue145

import (
	"context"
	"sync"
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/utils"
	asyncapi_test "github.com/lerenn/asyncapi-codegen/test"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	brokers, cleanup := asyncapi_test.BrokerControllers(t)
	defer cleanup()

	for _, b := range brokers {
		suite.Run(t, NewSuite(b))
	}
}

type Suite struct {
	broker extensions.BrokerController
	app    *AppController
	user   *UserController
	suite.Suite
}

func NewSuite(broker extensions.BrokerController) *Suite {
	return &Suite{
		broker: broker,
	}
}

func (suite *Suite) SetupTest() {
	// Create app
	app, err := NewAppController(suite.broker)
	suite.Require().NoError(err)
	suite.app = app

	// Create user
	user, err := NewUserController(suite.broker)
	suite.Require().NoError(err)
	suite.user = user
}

func (suite *Suite) TearDownTest() {
	suite.app.Close(context.Background())
	suite.user.Close(context.Background())
}

func (suite *Suite) TestRequestReplyWithReplyHelper() {
	// Listen for pings on the application
	err := suite.app.SubscribeToPingRequestOperation(
		context.Background(),
		func(ctx context.Context, ping PingMessage) error {
			callbackErr := suite.app.ReplyToPingRequestOperation(ctx, ping, func(pong *PongMessage) {
				pong.Payload.Event = ping.Payload.Event
			})
			suite.Require().NoError(callbackErr)
			return nil
		})
	suite.Require().NoError(err)
	defer suite.app.UnsubscribeFromPingRequestOperation(context.Background())

	// Set a new ping
	var msg PingMessage
	msg.Payload.Event = utils.ToPointer("testing")
	msg.Headers.ReplyTo = utils.ToPointer("issue145.pong.1234")

	// Send a request
	resp, err := suite.user.RequestToPingRequestOperation(context.Background(), msg)
	suite.Require().NoError(err)

	// Check response
	suite.Require().Equal(*msg.Payload.Event, *resp.Payload.Event)
}

func (suite *Suite) TestRequestReplyOnRawChannel() {
	// Listen for pings on the application
	err := suite.app.SubscribeToPingRequestOperation(
		context.Background(),
		func(ctx context.Context, ping PingMessage) error {
			callbackErr := suite.app.ReplyToPingRequestOperation(ctx, ping, func(pong *PongMessage) {
				pong.Payload.Event = ping.Payload.Event
			})
			suite.Require().NoError(callbackErr)
			return nil
		})
	suite.Require().NoError(err)
	defer suite.app.UnsubscribeFromPingRequestOperation(context.Background())

	// Listen directly for reply from the broker
	sub, err := suite.broker.Subscribe(context.Background(), "issue145.pong.2345")
	suite.Require().NoError(err)
	defer sub.Cancel(context.Background())

	var wg sync.WaitGroup
	go func() {
		rawReply := <-sub.MessagesChannel()
		reply, err := newPongMessageFromBrokerMessage(rawReply.BrokerMessage)
		suite.Require().NoError(err)
		suite.Require().NotNil(reply.Payload.Event)
		suite.Require().Equal("testing.2345", *reply.Payload.Event)

		wg.Done()
	}()
	wg.Add(1)

	// Set a new ping
	var msg PingMessage
	msg.Payload.Event = utils.ToPointer("testing.2345")
	msg.Headers.ReplyTo = utils.ToPointer("issue145.pong.2345")

	// Send a request
	err = suite.user.SendToPingRequestOperation(context.Background(), msg)
	suite.Require().NoError(err)

	// Wait for the end
	wg.Wait()
}
