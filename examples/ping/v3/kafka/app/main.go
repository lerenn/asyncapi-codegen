//go:generate go run ../../../../../cmd/asyncapi-codegen -g application,types -p main -i ../../asyncapi.yaml -o ./app.gen.go

package main

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/examples"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/kafka"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/middlewares"
)

type Subscriber struct {
	Controller *AppController
}

func (s Subscriber) PingRequestOperationReceived(ctx context.Context, req Ping) {
	// Generate a pong message, set as a response of the request
	resp := NewPong()
	resp.SetAsResponseFrom(&req)
	// -- You can modifiy the response here

	// Publish the pong message
	// Note: it will indefinitely wait to publish as context has no timeout
	err := s.Controller.ReplyToPingRequestOperation(ctx, resp)
	if err != nil {
		panic(err)
	}
}

func main() {
	// Instantiate a Kafka controller with a logger
	logger := loggers.NewText()
	broker, err := kafka.NewController(
		[]string{"kafka:9092"},         // List of hosts
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
	if err := ctrl.SubscribeToAllOperations(context.Background(), sub); err != nil {
		panic(err)
	}

	// Listen on port to let know that app is ready
	examples.ListenLocalPort(1234)
}
