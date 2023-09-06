//go:generate go run ../../../cmd/asyncapi-codegen -g application,types -p main -i ../asyncapi.yaml -o ./app.gen.go

package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/loggers"
	"github.com/nats-io/nats.go"
)

type ServerSubscriber struct {
	Controller *AppController
}

func (s ServerSubscriber) Ping(ctx context.Context, req PingMessage, _ bool) {
	// Generate a pong message, set as a response of the request
	resp := NewPongMessage()
	resp.SetAsResponseFrom(req)
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
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}

	// Create a new server controller
	ctrl, err := NewAppController(brokers.NewNATS(nc))
	if err != nil {
		panic(err)
	}
	defer ctrl.Close(context.Background())

	// Attach a logger (optional)
	logger := loggers.NewECS()
	ctrl.SetLogger(logger)

	// Subscribe to all (we could also have just listened on the ping request channel)
	sub := ServerSubscriber{Controller: ctrl}
	if err := ctrl.SubscribeAll(context.Background(), sub); err != nil {
		panic(err)
	}

	// Process messages until interruption signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
