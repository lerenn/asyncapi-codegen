//go:generate go run ../../../../../cmd/asyncapi-codegen -g user,types -p main -i ../../asyncapi.yaml -o ./user.gen.go

package main

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/middlewares"
)

func main() {
	// Instantiate a NATS controller with a logger
	logger := loggers.NewText()
	broker, err := nats.NewController(
		"nats://nats:4222",                // Set URL to broker
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
	// -- you can modifiy the request here

	// The following function will subscribe to the 'pong' channel (reply channel
	// to PingRequest operation) and wait for a response. The response will be
	// detected through its correlation ID. However, if there is no correlation
	// ID, then it will return the first message on the reply channel.
	//
	// Note: it will indefinitely wait for messages as context has no timeout
	_, err = ctrl.RequestPongMessageOnPongChannelWithPingMessageOnPingChannel(context.Background(), req)
	if err != nil {
		panic(err)
	}
}
