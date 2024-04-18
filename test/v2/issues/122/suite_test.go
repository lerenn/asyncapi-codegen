//go:generate go run ../../../../cmd/asyncapi-codegen -p issue122 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue122

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	testutil "github.com/lerenn/asyncapi-codegen/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var errTest = errors.New("some test error")

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

	wg sync.WaitGroup
}

func NewSuite(broker extensions.BrokerController) *Suite {
	return &Suite{
		broker: broker,
	}
}

func (suite *Suite) SetupTest() {
	// Errorhandler that should do some error handling for app and user controller
	testErrorHandler := func(ctx context.Context, topic string, msg *extensions.AcknowledgeableBrokerMessage, err error) {
		assert.ErrorIs(suite.T(), err, errTest, fmt.Sprintf("want %v, have %v", errTest, err))
		suite.wg.Done()
	}

	// Create app
	app, err := NewAppController(suite.broker, WithErrorHandler(testErrorHandler))
	suite.Require().NoError(err)
	suite.app = app

	// Create user
	user, err := NewUserController(suite.broker, WithErrorHandler(testErrorHandler))
	suite.Require().NoError(err)
	suite.user = user
}

func (suite *Suite) TearDownTest() {
	suite.app.Close(context.Background())
	suite.user.Close(context.Background())
}

func (suite *Suite) TestErrorHandlerForApp() {
	// Test message
	sent := V2Issue122MsgPublishMessage{
		Payload: "test some errors",
	}

	// return some error on message
	err := suite.app.SubscribeV2Issue122Msg(
		context.Background(),
		func(_ context.Context, msg V2Issue122MsgSubscribeMessage) error {
			return errTest
		})
	suite.Require().NoError(err)
	defer suite.app.UnsubscribeV2Issue122Msg(context.Background())

	// Publish the message
	suite.wg.Add(1)
	err = suite.user.PublishV2Issue122Msg(context.Background(), sent)
	suite.Require().NoError(err)

	// Wait for errorhandler is called
	suite.wg.Wait()
}

func (suite *Suite) TestErrorHandlerForUser() {
	// Test message
	sent := V2Issue122MsgPublishMessage{
		Payload: "test some errors",
	}

	// return some error on message
	err := suite.user.SubscribeV2Issue122Msg(
		context.Background(),
		func(_ context.Context, msg V2Issue122MsgSubscribeMessage) error {
			return errTest
		})
	suite.Require().NoError(err)
	defer suite.user.UnsubscribeV2Issue122Msg(context.Background())

	// Publish the message
	suite.wg.Add(1)
	err = suite.user.PublishV2Issue122Msg(context.Background(), sent)
	suite.Require().NoError(err)

	// Wait for errorhandler is called
	suite.wg.Wait()
}
