//go:generate go run ../../../../cmd/asyncapi-codegen -g types,user -p none -i ./asyncapi.yaml -o ./none/asyncapi.gen.go
//go:generate go run ../../../../cmd/asyncapi-codegen -g types,user -p snake --convert-keys snake -i ./asyncapi.yaml -o ./snake/asyncapi.gen.go
//go:generate go run ../../../../cmd/asyncapi-codegen -g types,user -p camel --convert-keys camel -i ./asyncapi.yaml -o ./camel/asyncapi.gen.go
//go:generate go run ../../../../cmd/asyncapi-codegen -g types,user -p kebab --convert-keys kebab -i ./asyncapi.yaml -o ./kebab/asyncapi.gen.go

package issue129

import (
	"context"
	"testing"

	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions"
	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions/middlewares"
	"github.com/TheSadlig/asyncapi-codegen/pkg/utils"
	testutil "github.com/TheSadlig/asyncapi-codegen/test"
	"github.com/TheSadlig/asyncapi-codegen/test/v2/issues/129/camel"
	"github.com/TheSadlig/asyncapi-codegen/test/v2/issues/129/kebab"
	"github.com/TheSadlig/asyncapi-codegen/test/v2/issues/129/none"
	"github.com/TheSadlig/asyncapi-codegen/test/v2/issues/129/snake"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	brokers, cleanup := testutil.BrokerControllers(t)
	defer cleanup()

	for _, b := range brokers {
		suite.Run(t, NewSuite(b))
	}
}

type Suite struct {
	broker extensions.BrokerController

	suite.Suite
}

func NewSuite(broker extensions.BrokerController) *Suite {
	return &Suite{
		broker: broker,
	}
}

func (suite *Suite) TestWithNoneKeyConversion() {
	// Create a channel to intercept message before sending to broker and after
	// reception from broker
	interceptor := make(chan extensions.BrokerMessage, 8)
	defer close(interceptor)

	// Create none user
	user, err := none.NewUserController(suite.broker, none.WithMiddlewares(middlewares.Intercepter(interceptor)))
	suite.Require().NoError(err)
	defer user.Close(context.Background())

	// Send the message
	err = user.PublishV2Issue129Test(context.Background(), none.V2Issue129TestMessage{
		Payload: none.TestSchema{
			ThisIsAProperty: utils.ToPointer("value"),
		},
	})
	suite.Require().NoError(err)

	// Check sent message to broker
	bMsg := <-interceptor

	// Check that the additional properties are at the level 0 of payload
	suite.Require().Equal("{\"This_is a-Property\":\"value\"}", string(bMsg.Payload))
}

func (suite *Suite) TestWithSnakeKeyConversion() {
	// Create a channel to intercept message before sending to broker and after
	// reception from broker
	interceptor := make(chan extensions.BrokerMessage, 8)
	defer close(interceptor)

	// Create snake user
	user, err := snake.NewUserController(suite.broker, snake.WithMiddlewares(middlewares.Intercepter(interceptor)))
	suite.Require().NoError(err)
	defer user.Close(context.Background())

	// Send the message
	err = user.PublishV2Issue129Test(context.Background(), snake.V2Issue129TestMessage{
		Payload: snake.TestSchema{
			ThisIsAProperty: utils.ToPointer("value"),
		},
	})
	suite.Require().NoError(err)

	// Check sent message to broker
	bMsg := <-interceptor

	// Check that the additional properties are at the level 0 of payload
	suite.Require().Equal("{\"this_is_a_property\":\"value\"}", string(bMsg.Payload))
}

func (suite *Suite) TestWithKebabKeyConversion() {
	// Create a channel to intercept message before sending to broker and after
	// reception from broker
	interceptor := make(chan extensions.BrokerMessage, 8)
	defer close(interceptor)

	// Create kebab user
	user, err := kebab.NewUserController(suite.broker, kebab.WithMiddlewares(middlewares.Intercepter(interceptor)))
	suite.Require().NoError(err)
	defer user.Close(context.Background())

	// Send the message
	err = user.PublishV2Issue129Test(context.Background(), kebab.V2Issue129TestMessage{
		Payload: kebab.TestSchema{
			ThisIsAProperty: utils.ToPointer("value"),
		},
	})
	suite.Require().NoError(err)

	// Check sent message to broker
	bMsg := <-interceptor

	// Check that the additional properties are at the level 0 of payload
	suite.Require().Equal("{\"this-is-a-property\":\"value\"}", string(bMsg.Payload))
}

func (suite *Suite) TestWithCamelKeyConversion() {
	// Create a channel to intercept message before sending to broker and after
	// reception from broker
	interceptor := make(chan extensions.BrokerMessage, 8)
	defer close(interceptor)

	// Create camel user
	user, err := camel.NewUserController(suite.broker, camel.WithMiddlewares(middlewares.Intercepter(interceptor)))
	suite.Require().NoError(err)
	defer user.Close(context.Background())

	// Send the message
	err = user.PublishV2Issue129Test(context.Background(), camel.V2Issue129TestMessage{
		Payload: camel.TestSchema{
			ThisIsAProperty: utils.ToPointer("value"),
		},
	})
	suite.Require().NoError(err)

	// Check sent message to broker
	bMsg := <-interceptor

	// Check that the additional properties are at the level 0 of payload
	suite.Require().Equal("{\"ThisIsAProperty\":\"value\"}", string(bMsg.Payload))
}
