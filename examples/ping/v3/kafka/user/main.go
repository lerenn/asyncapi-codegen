//go:generate go run ../../../../../cmd/asyncapi-codegen -g user,types -p main -i ../../asyncapi.yaml -o ./user.gen.go

package main

import (
	"context"

	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions/brokers/kafka"
	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions/loggers"
	"github.com/TheSadlig/asyncapi-codegen/pkg/extensions/middlewares"
	testutil "github.com/TheSadlig/asyncapi-codegen/pkg/utils/test"
)

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
		[]string{addr},                  // List of hosts
		kafka.WithLogger(logger),        // Attach an internal logger
		kafka.WithGroupID("ping-users"), // Change group id
	)
	if err != nil {
		panic(err)
	}

	// Create a new user controller
	ctrl, err := NewUserController(
		broker,             // Attach the kafka controller
		WithLogger(logger), // Attach an internal logger
		WithMiddlewares(middlewares.Logging(logger))) // Attach a middleware to log messages
	if err != nil {
		panic(err)
	}
	defer ctrl.Close(context.Background())

	// Make a new ping message
	req := NewPingMessage()
	// -- you can modifiy the request here

	// The following function will subscribe to the 'pong' channel (reply channel
	// to PingRequest operation) and wait for a response. The response will be
	// detected through its correlation ID. However, if there is no correlation
	// ID, then it will return the first message on the reply channel.
	//
	// Note: it will indefinitely wait for messages as context has no timeout
	_, err = ctrl.RequestToPingRequestOperation(context.Background(), req)
	if err != nil {
		panic(err)
	}
}
