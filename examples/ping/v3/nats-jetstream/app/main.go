//go:generate go run ../../../../../cmd/asyncapi-codegen -g application,types -p main -i ../../asyncapi.yaml -o ./app.gen.go

package main

import (
	"context"

	"github.com/TheSadlig/asyncapi-codegen/examples"
	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions/brokers/natsjetstream"
	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions/loggers"
	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions/middlewares"
	testutil "github.com/TheSadlig/asyncapi-codegen/pkg/utils/test"
	"github.com/nats-io/nats.go/jetstream"
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
	// Get broker address based on the environment, it will returns an address like "nats://nats-jetstream:4222"
	// Note: this is not needed in your application, you can directly use the address
	addr := testutil.BrokerAddress(testutil.BrokerAddressParams{
		Schema:         "nats",
		DockerizedAddr: "nats-jetstream",
		DockerizedPort: "4222",
		LocalPort:      "4225",
	})

	// Instantiate a NATS controller with a logger
	logger := loggers.NewText()
	broker, err := natsjetstream.NewController(
		addr,                             // Set URL to broker
		natsjetstream.WithLogger(logger), // Attach an internal logger
		natsjetstream.WithStreamConfig(jetstream.StreamConfig{
			Name: "pingv3",
			Subjects: []string{
				"ping.v3", "pong.v3",
			},
		}), // Create the stream "ping"
		natsjetstream.WithConsumerConfig(jetstream.ConsumerConfig{Name: "pingv3"}), // Create the corresponding consumer
	)
	if err != nil {
		panic(err)
	}
	defer broker.Close()

	// Create a new app controller
	ctrl, err := NewAppController(
		broker,             // Attach the NATS controller
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

	// Listen on port to let know that app is ready
	examples.ListenLocalPort(1234)
}
