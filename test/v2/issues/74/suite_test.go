//go:generate go run ../../../../cmd/asyncapi-codegen -p issue74 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue74

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/middlewares"
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
	broker      extensions.BrokerController
	app         *AppController
	user        *UserController
	interceptor chan extensions.BrokerMessage
	suite.Suite
}

func NewSuite(broker extensions.BrokerController) *Suite {
	return &Suite{
		broker: broker,
	}
}

func (suite *Suite) SetupTest() {
	// Create a channel to intercept message before sending to broker and after
	// reception from broker
	suite.interceptor = make(chan extensions.BrokerMessage, 8)

	// Create app
	app, err := NewAppController(suite.broker, WithMiddlewares(middlewares.Intercepter(suite.interceptor)))
	suite.Require().NoError(err)
	suite.app = app

	// Create user
	user, err := NewUserController(suite.broker, WithMiddlewares(middlewares.Intercepter(suite.interceptor)))
	suite.Require().NoError(err)
	suite.user = user
}

func (suite *Suite) TearDownTest() {
	suite.app.Close(context.Background())
	suite.user.Close(context.Background())
	close(suite.interceptor)
}

func (suite *Suite) TestHeaders() {
	var wg sync.WaitGroup

	// Expected message
	sent := TestMessage{
		Headers: HeaderSchema{
			DateTime: utils.Must(time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")).UTC(),
			Version:  "1.0.0",
		},
	}

	// Check what the app receive and translate
	var recvMsg TestMessage
	wg.Add(1)
	err := suite.app.SubscribeIssue74TestChannel(context.Background(), func(_ context.Context, msg TestMessage) {
		recvMsg = msg
		wg.Done()
	})
	suite.Require().NoError(err)

	// Publish the message
	err = suite.user.PublishIssue74TestChannel(context.Background(), sent)
	suite.Require().NoError(err)

	// Wait for the message to be received by the app
	wg.Wait()

	// Check received message
	suite.Require().Equal(sent, recvMsg)

	// Check sent message to broker
	bMsg := <-suite.interceptor

	// Check that version is in the header
	version, exists := bMsg.Headers["version"]
	suite.Require().True(exists)
	suite.Require().Equal([]byte("1.0.0"), version)

	// Check that datetime is in the header
	datetime, exists := bMsg.Headers["dateTime"]
	suite.Require().True(exists)
	suite.Require().Equal([]byte("2020-01-01T00:00:00Z"), datetime)
}
