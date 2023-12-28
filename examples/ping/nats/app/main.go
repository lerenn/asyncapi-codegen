//go:generate go run ../../../../cmd/asyncapi-codegen -g application,types -p main -i ../../asyncapi.yaml -o ./app.gen.go

package main

import (
	"context"
	"time"

	"github.com/lerenn/asyncapi-codegen/examples"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/middlewares"
)

type ServerSubscriber struct {
	Controller *AppController
}

func (s ServerSubscriber) Ping(ctx context.Context, req PingMessage) {
	// Generate a pong message, set as a response of the request
	resp := NewPongMessage()
	resp.SetAsResponseFrom(&req)
	resp.Payload.Message = "pong"
	resp.Payload.Time = time.Now()

	// Publish the pong message
	// Note: it will indefinitely wait to publish as context has no timeout
	err := s.Controller.PublishPong(ctx, resp)
	if err != nil {
		panic(err)
	}
}

func main() {
	// Instantiate a NATS controller with a logger
	logger := loggers.NewText()
	broker, err := nats.NewController(
		"nats://nats:4222",               // Set URL to broker
		nats.WithLogger(logger),          // Attach an internal logger
		nats.WithQueueGroup("ping-apps"), // Set a specific queue group to avoid collisions
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

	// Subscribe to all (we could also have just listened on the ping request channel)
	sub := ServerSubscriber{Controller: ctrl}
	if err := ctrl.SubscribeAll(context.Background(), sub); err != nil {
		panic(err)
	}

	// Listen on port to let know that app is ready
	examples.ListenLocalPort(1234)
}
