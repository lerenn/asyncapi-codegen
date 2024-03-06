//go:generate go run ../../../../../cmd/asyncapi-codegen -p decoupling -i ./asyncapi.yaml -o ./asyncapi.gen.go

//nolint:revive
package decoupling

import (
	"context"
	"sync"

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

func (suite *Suite) TestSendReceive() {
	var wg sync.WaitGroup

	// Listen to new messages
	err := suite.app.SubscribeToUserFromUserSignupChannel(
		context.Background(),
		func(ctx context.Context, msg UserMessage) {
			suite.Require().NotNil(msg.Payload.DisplayName)
			suite.Require().Equal("testing", *msg.Payload.DisplayName)
			wg.Done()
		})
	suite.Require().NoError(err)
	defer suite.app.UnsubscribeFromUserFromUserSignupChannel(context.Background())
	wg.Add(1)

	// Set a new message
	var msg UserMessage
	msg.Payload.DisplayName = utils.ToPointer("testing")

	// Send the new message
	err = suite.user.PublishUserOnUserSignupChannel(context.Background(), msg)
	suite.Require().NoError(err)

	wg.Wait()
}
