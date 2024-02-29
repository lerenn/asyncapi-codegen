//go:generate go run ../../../../../cmd/asyncapi-codegen -p requestreply -i ./asyncapi.yaml -o ./asyncapi.gen.go

//nolint:revive
package requestreply

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/utils"
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
	err := suite.app.SubscribeToPingRequestOperation(context.Background(), func(ctx context.Context, msg Ping) {
		var respMsg Pong
		respMsg.Payload.Event = msg.Payload.Event
		err := suite.app.ReplyToPingRequestOperation(ctx, respMsg)
		suite.Require().NoError(err)
	})
	suite.Require().NoError(err)
	defer suite.app.UnsubscribeFromPingRequestOperation(context.Background())

	// Set a new ping
	var msg Ping
	msg.Payload.Event = utils.ToPointer("testing")

	// Send a request
	resp, err := suite.user.RequestToPingRequestOperation(context.Background(), msg)
	suite.Require().NoError(err)

	// Check response
	suite.Require().Equal(*msg.Payload.Event, *resp.Payload.Event)
}

func (suite *Suite) TestRequestReplyWithCorrelationID() {
	// Listen to new pings
	err := suite.app.SubscribeToPingRequestWithCorrelationIDOperation(context.Background(),
		func(ctx context.Context, msg PingWithCorrelationID) {
			// Set response
			var respMsg PongWithCorrelationID
			respMsg.SetAsResponseFrom(&msg)
			respMsg.Payload.Event = msg.Payload.Event

			// Send response
			err := suite.app.ReplyToPingRequestWithCorrelationIDOperation(ctx, respMsg)
			suite.Require().NoError(err)
		})
	suite.Require().NoError(err)
	defer suite.app.UnsubscribeFromPingRequestWithCorrelationIDOperation(context.Background())

	// Set a new ping
	var msg PingWithCorrelationID
	msg.Payload.Event = utils.ToPointer("testing")

	// Send a request
	resp, err := suite.user.RequestToPingRequestWithCorrelationIDOperation(context.Background(), msg)
	suite.Require().NoError(err)

	// Check response
	suite.Require().Equal(*msg.Payload.Event, *resp.Payload.Event)
}
