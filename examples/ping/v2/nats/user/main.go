//go:generate go run ../../../../../cmd/asyncapi-codegen -g user,types -p main -i ../../asyncapi.yaml -o ./user.gen.go

package main

import (
	"context"

	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions/loggers"
	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions/middlewares"
	testutil "github.com/TheSadlig/asyncapi-codegen/pkg/utils/test"
)

func main() {
	// Get broker address based on the environment, it will returns an address like "nats://nats:4222"
	// Note: this is not needed in your application, you can directly use the address
	addr := testutil.BrokerAddress(testutil.BrokerAddressParams{
		Schema:         "nats",
		DockerizedAddr: "nats",
		Port:           "4222",
	})

	// Instantiate a NATS controller with a logger
	logger := loggers.NewText()
	broker, err := nats.NewController(
		addr,                              // Set URL to broker
		nats.WithLogger(logger),           // Attach an internal logger
		nats.WithQueueGroup("ping-users"), // Set a specific queue group to avoid collisions
	)
	if err != nil {
		panic(err)
	}
	defer broker.Close()

	// Create a new user controller
	ctrl, err := NewUserController(
		broker,             // Attach the NATS controller
		WithLogger(logger), // Attach an internal logger
		WithMiddlewares(middlewares.Logging(logger))) // Attach a middleware to log messages
	if err != nil {
		panic(err)
	}
	defer ctrl.Close(context.Background())

	// Make a new ping message
	req := NewPingMessage()
	req.Payload = "ping"

	// Create the publication function to send the message
	// Note: it will indefinitely wait to publish as context has no timeout
	publicationFunc := func(ctx context.Context) error {
		return ctrl.PublishPing(ctx, req)
	}

	// The following function will subscribe to the 'pong' channel, execute the publication
	// function and wait for a response. The response will be detected through its
	// correlation ID.
	//
	// This function is available only if the 'correlationId' field has been filled
	// for any channel in the AsyncAPI specification. You will then be able to use it
	// with the form WaitForXXX where XXX is the channel name.
	//
	// Note: it will indefinitely wait for messages as context has no timeout
	_, err = ctrl.WaitForPong(context.Background(), &req, publicationFunc)
	if err != nil {
		panic(err)
	}
}
