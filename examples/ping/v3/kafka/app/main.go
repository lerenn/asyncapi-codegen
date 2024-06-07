//go:generate go run ../../../../../cmd/asyncapi-codegen -g application,types -p main -i ../../asyncapi.yaml -o ./app.gen.go

package main

import (
	"context"

	"github.com/TheSadlig/asyncapi-codegen/examples"
	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions/brokers/kafka"
	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions/loggers"
	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions/middlewares"
	testutil "github.com/TheSadlig/asyncapi-codegen/pkg/utils/test"
)

var _ AppSubscriber = (*Subscriber)(nil)

type Subscriber struct {
	Controller *AppController
}

func (s Subscriber) PingRequestOperationReceived(ctx context.Context, ping PingMessage) error {
	// Publish the pong message, with the callback function to modify it
	// Note: it will indefinitely wait to publish as context has no timeout
	err := s.Controller.ReplyToPingRequestOperation(ctx, ping, func(pong *PongMessage) {
		// Reply with the same event than the ping
		pong.Payload.Event = ping.Payload.Event
	})

	// Error management
	if err != nil {
		panic(err)
	}

	return nil
}

func main() {
	// Get broker address based on the environment, it will returns an address like "kafka:9092"
	// Note: this is not needed in your application, you can directly use the address
	addr := testutil.BrokerAddress(testutil.BrokerAddressParams{
		DockerizedAddr: "kafka",
		Port:           "9092",
	})

	// Instantiate a Kafka controller with a logger
	logger := loggers.NewText()
	broker, err := kafka.NewController(
		[]string{addr},                 // List of hosts
		kafka.WithLogger(logger),       // Attach an internal logger
		kafka.WithGroupID("ping-apps"), // Change group id
	)
	if err != nil {
		panic(err)
	}

	// Create a new app controller
	ctrl, err := NewAppController(
		broker,             // Attach the kafka controller
		WithLogger(logger), // Attach an internal logger
		WithMiddlewares(middlewares.Logging(logger))) // Attach a middleware to log messages
	if err != nil {
		panic(err)
	}
	defer ctrl.Close(context.Background())

	// Subscribe to all (we could also have just subscribed to the ping request operation)
	sub := Subscriber{Controller: ctrl}
	if err := ctrl.SubscribeToAllChannels(context.Background(), sub); err != nil {
		panic(err)
	}
	defer ctrl.UnsubscribeFromAllChannels(context.Background())

	// Listen on port to let know that app is ready
	examples.ListenLocalPort(1234)
}
