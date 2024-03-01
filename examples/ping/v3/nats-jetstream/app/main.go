//go:generate go run ../../../../../cmd/asyncapi-codegen -g application,types -p main -i ../../asyncapi.yaml -o ./app.gen.go

package main

import (
	"context"

	"github.com/lerenn/asyncapi-codegen/examples"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/natsjetstream"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/middlewares"
	"github.com/nats-io/nats.go/jetstream"
)

type ServerSubscriber struct {
	Controller *AppController
}

func (s ServerSubscriber) PingRequestOperationReceived(ctx context.Context, req Ping) {
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
	// Instantiate a NATS controller with a logger
	logger := loggers.NewText()
	broker, err := natsjetstream.NewController(
		"nats://nats-jetstream:4222",     // Set URL to broker
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
	sub := ServerSubscriber{Controller: ctrl}
	if err := ctrl.SubscribeToAllOperations(context.Background(), sub); err != nil {
		panic(err)
	}

	// Listen on port to let know that app is ready
	examples.ListenLocalPort(1234)
}