//go:generate go run ../../../../cmd/asyncapi-codegen -p issue130 -i ./asyncapi.yaml -o ./asyncapi.gen.go

// More info: https://github.com/lerenn/asyncapi-codegen/issues/130

package issue130

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

	// Only do it with one broker as this is not testing the broker
	suite.Run(t, NewSuite(brokers[0]))
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

func (suite *Suite) TestSendReceive() {
	var wg sync.WaitGroup

	// Listen to new messages
	suite.app.ListenConsumeUserSignup(context.Background(), func(ctx context.Context, msg UserMessage) {
		suite.Require().NotNil(msg.Payload.DisplayName)
		suite.Require().Equal("testing", *msg.Payload.DisplayName)
		wg.Done()
	})
	wg.Add(1)

	// Send a new message
	var msg UserMessage
	msg.Payload.DisplayName = utils.ToPointer("testing")
	suite.user.SendConsumeUserSignup(context.Background(), msg)

	wg.Wait()
}
