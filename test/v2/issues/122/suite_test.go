//go:generate go run ../../../../cmd/asyncapi-codegen -p issue122 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue122

import (
	"context"
	"errors"
	"fmt"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	asyncapi_test "github.com/lerenn/asyncapi-codegen/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
)

var errTest = errors.New("some test error")

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
	sent := Issue122MsgMessage{
		Payload: "test some errors",
	}

	// return some error on message
	err := suite.app.SubscribeIssue122Msg(context.Background(), func(_ context.Context, msg Issue122MsgMessage) error {
		return errTest
	})
	suite.Require().NoError(err)

	suite.wg.Add(1)

	// Publish the message
	err = suite.user.PublishIssue122Msg(context.Background(), sent)
	suite.Require().NoError(err)

	// Wait for errorhandler is called
	suite.wg.Wait()
}

func (suite *Suite) TestErrorHandlerForUser() {
	// Test message
	sent := Issue122MsgMessage{
		Payload: "test some errors",
	}

	// return some error on message
	err := suite.user.SubscribeIssue122Msg(context.Background(), func(_ context.Context, msg Issue122MsgMessage) error {
		return errTest
	})
	suite.Require().NoError(err)

	suite.wg.Add(1)

	// Publish the message
	err = suite.user.PublishIssue122Msg(context.Background(), sent)
	suite.Require().NoError(err)

	// Wait for errorhandler is called
	suite.wg.Wait()
}
