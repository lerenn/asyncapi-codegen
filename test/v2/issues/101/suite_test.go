//go:generate go run ../../../../cmd/asyncapi-codegen -p issue101 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue101

import (
	"context"
	"sync"
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	testutil "github.com/lerenn/asyncapi-codegen/test"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	brokers, cleanup := testutil.BrokerControllers(t)
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
	// Middleware that adds info on context and check it
	m1 := func(ctx context.Context, _ *extensions.BrokerMessage, next extensions.NextMiddleware) error {
		ctx = context.WithValue(ctx, "test-ctx-passing-middlewares", "value passed") //nolint:staticcheck
		return next(ctx)
	}
	m2 := func(ctx context.Context, msg *extensions.BrokerMessage, _ extensions.NextMiddleware) error {
		suite.Require().NotNil(ctx.Value("test-ctx-passing-middlewares"))
		suite.Require().Equal("value passed", ctx.Value("test-ctx-passing-middlewares"))
		return nil
	}

	// Create app
	app, err := NewAppController(suite.broker, WithMiddlewares(m1, m2))
	suite.Require().NoError(err)
	suite.app = app

	// Create user
	user, err := NewUserController(suite.broker, WithMiddlewares(m1, m2))
	suite.Require().NoError(err)
	suite.user = user
}

func (suite *Suite) TearDownTest() {
	suite.app.Close(context.Background())
	suite.user.Close(context.Background())
}

func (suite *Suite) TestAddingHeader() {
	var wg sync.WaitGroup

	// Expected message
	sent := V2Issue101TestMessage{
		Payload: "hello!",
	}

	// Check what the app receive
	err := suite.app.SubscribeV2Issue101Test(
		context.Background(),
		func(_ context.Context, msg V2Issue101TestMessage) error {
			wg.Done()
			return nil
		})
	suite.Require().NoError(err)

	// Publish the message
	wg.Add(1)
	err = suite.user.PublishV2Issue101Test(context.Background(), sent)
	suite.Require().NoError(err)

	// Wait for the message to be received by the app
	wg.Wait()
}
