//go:generate go run ../../../../cmd/asyncapi-codegen -p issue99 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue99

import (
	"context"
	"sync"
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/middlewares"
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

	// Middleware that adds info on emi
	m := func(ctx context.Context, msg *extensions.BrokerMessage, _ extensions.NextMiddleware) error {
		extensions.IfContextValueEquals(ctx, extensions.ContextKeyIsDirection, "publication", func() {
			msg.Headers["additional"] = []byte("some-info")
		})
		return nil
	}

	// Create app
	app, err := NewAppController(suite.broker, WithMiddlewares(m))
	suite.Require().NoError(err)
	suite.app = app

	// Create user
	user, err := NewUserController(suite.broker, WithMiddlewares(middlewares.Intercepter(suite.interceptor), m))
	suite.Require().NoError(err)
	suite.user = user
}

func (suite *Suite) TearDownTest() {
	suite.app.Close(context.Background())
	suite.user.Close(context.Background())
	close(suite.interceptor)
}

func (suite *Suite) TestAddingHeader() {
	var wg sync.WaitGroup

	// Expected message
	sent := Issue99TestMessage{
		Payload: "hello!",
	}

	// Check what the app receive and translate
	var recvMsg Issue99TestMessage
	wg.Add(1)
	err := suite.app.SubscribeIssue99Test(context.Background(), func(_ context.Context, msg Issue99TestMessage) error {
		recvMsg = msg
		wg.Done()
		return nil
	})
	suite.Require().NoError(err)

	// Publish the message
	err = suite.user.PublishIssue99Test(context.Background(), sent)
	suite.Require().NoError(err)

	// Wait for the message to be received by the app
	wg.Wait()

	// Check received message
	suite.Require().Equal(sent, recvMsg)

	// Check sent message to broker
	bMsg := <-suite.interceptor

	// Check that additional field is in the header
	header, exists := bMsg.Headers["additional"]
	suite.Require().True(exists)
	suite.Require().Equal([]byte("some-info"), header)
}
