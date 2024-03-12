//go:generate go run ../../../../cmd/asyncapi-codegen -p issue145 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue145

import (
	"context"
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
	err := suite.app.SubscribeToPingFromPingChannel(
		context.Background(),
		func(ctx context.Context, ping PingMessage) {
			err := suite.app.ReplyToPingWithPongOnPongChannel(ctx, ping, func(pong *PongMessage) {
				pong.Payload.Event = ping.Payload.Event
			})
			suite.Require().NoError(err)
		})
	suite.Require().NoError(err)
	defer suite.app.UnsubscribeFromPingFromPingChannel(context.Background())

	// Set a new ping
	var msg PingMessage
	msg.Payload.Event = utils.ToPointer("testing")
	msg.Headers.ReplyTo = utils.ToPointer("pong.1234")

	// Send a request
	resp, err := suite.user.RequestPongOnPongChannelWithPingOnPingChannel(context.Background(), msg)
	suite.Require().NoError(err)

	// Check response
	suite.Require().Equal(*msg.Payload.Event, *resp.Payload.Event)
}

func (suite *Suite) TestRequestReplyWithManualReply() {
	// TODO
}

func (suite *Suite) TestRequestReplyOnRawChannel() {
	// TODO
}
