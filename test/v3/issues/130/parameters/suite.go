//go:generate go run ../../../../../cmd/asyncapi-codegen -p parameters -i ./asyncapi.yaml -o ./asyncapi.gen.go

//nolint:revive
package parameters

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

func (suite *Suite) TestParameter() {
	var wg sync.WaitGroup

	// Set parameters
	params := UserSignupChannelParameters{
		UserId: "1234",
	}

	// Listen to new messages
	err := suite.app.SubscribeToUsersFromUserSignupChannel(
		context.Background(),
		params,
		func(ctx context.Context, msg UserMessage) {
			suite.Require().NotNil(msg.Payload.Name)
			suite.Require().Equal("testing", *msg.Payload.Name)
			wg.Done()
		})
	suite.Require().NoError(err)
	defer suite.app.UnsubscribeFromUsersFromUserSignupChannel(context.Background(), params)
	wg.Add(1)

	// Set a new message
	var msg UserMessage
	msg.Payload.Name = utils.ToPointer("testing")

	// Send the new message
	err = suite.user.PublishUserOnUserSignupChannel(context.Background(), params, msg)
	suite.Require().NoError(err)

	wg.Wait()
}
