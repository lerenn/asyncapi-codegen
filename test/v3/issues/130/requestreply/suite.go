//go:generate go run ../../../../../cmd/asyncapi-codegen -p requestreply -i ./asyncapi.yaml -o ./asyncapi.gen.go

//nolint:revive
package requestreply

import (
	"context"

	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions"
	"github.com/TheSadlig/asyncapi-codegen/pkg/utils"
	"github.com/stretchr/testify/suite"
)

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

func (suite *Suite) TestRequestReply() {
	// Listen for pings on the application
	err := suite.app.SubscribeToPingOperation(
		context.Background(),
		func(ctx context.Context, msg PingMessage) error {
			var respMsg PongMessage
			respMsg.Payload.Event = msg.Payload.Event
			err := suite.app.SendAsReplyToPingOperation(ctx, respMsg)
			suite.Require().NoError(err)
			return nil
		})
	suite.Require().NoError(err)
	defer suite.app.UnsubscribeFromPingOperation(context.Background())

	// Set a new ping
	var msg PingMessage
	msg.Payload.Event = utils.ToPointer("testing")

	// Send a request
	resp, err := suite.user.RequestToPingOperation(context.Background(), msg)
	suite.Require().NoError(err)

	// Check response
	suite.Require().Equal(*msg.Payload.Event, *resp.Payload.Event)
}

func (suite *Suite) TestRequestReplyWithID() {
	// Listen to new pings
	err := suite.app.SubscribeToPingWithIDOperation(context.Background(),
		func(ctx context.Context, msg PingWithIDMessage) error {
			// Set response
			var respMsg PongWithIDMessage
			respMsg.SetAsResponseFrom(&msg)
			respMsg.Payload.Event = msg.Payload.Event

			// Send response
			callbackErr := suite.app.SendAsReplyToPingWithIDOperation(ctx, respMsg)
			suite.Require().NoError(callbackErr)

			return nil
		})
	suite.Require().NoError(err)
	defer suite.app.UnsubscribeFromPingWithIDOperation(context.Background())

	// Set a new ping
	var msg PingWithIDMessage
	msg.Payload.Event = utils.ToPointer("testing")

	// Send a request
	resp, err := suite.user.RequestToPingWithIDOperation(
		context.Background(),
		msg,
	)
	suite.Require().NoError(err)

	// Check response
	suite.Require().Equal(*msg.Payload.Event, *resp.Payload.Event)
}
