//go:generate go run ../../../../cmd/asyncapi-codegen -p issue164 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue164

import (
	"context"
	"sync"
	"testing"

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

func (suite *Suite) TestAdditionalProperties() {
	var wg sync.WaitGroup

	// Expected message
	sent := TestMapMessage{
		Payload: TestMapSchema{
			Property: utils.ToPointer("value"),
			AdditionalProperties: map[string]string{
				"hello": "there",
			},
		},
	}

	// Check what the app receive and translate
	var recvMsg TestMapMessage
	wg.Add(1)
	err := suite.app.SubscribeIssue164TestMap(
		context.Background(),
		func(_ context.Context, msg TestMapMessage) {
			recvMsg = msg
			wg.Done()
		})
	suite.Require().NoError(err)

	// Send the message
	err = suite.user.PublishIssue164TestMap(context.Background(), sent)
	suite.Require().NoError(err)

	// Wait for the message to be received by the app
	wg.Wait()

	// Check received message
	suite.Require().Equal(sent, recvMsg)

	// Check sent message to broker
	bMsg := <-suite.interceptor

	// Check that the additional properties are at the level 0 of payload
	suite.Require().Equal(
		"{\"property\":\"value\",\"hello\":\"there\"}",
		string(bMsg.Payload),
	)
}
