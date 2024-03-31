//go:generate go run ../../../../../cmd/asyncapi-codegen -g application,types -p main -i ../../asyncapi.yaml -o ./app.gen.go

package main

import (
	"context"
	"log"

	"github.com/lerenn/asyncapi-codegen/examples"
	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
)

func main() {
	// Create a new broker adapter
	broker, err := nats.NewController("nats://nats:4222", nats.WithQueueGroup("helloworld-apps"))
	if err != nil {
		panic(err)
	}
	defer broker.Close()

	// Create a new application controller
	ctrl, err := NewAppController(broker)
	if err != nil {
		panic(err)
	}
	defer ctrl.Close(context.Background())

	// Subscribe to SayHelloMessage messages
	// Note: it will indefinitely wait for messages as context has no timeout
	log.Println("Subscribe to hello world...")
	err = ctrl.SubscribeToReceiveHelloOperation(context.Background(), func(_ context.Context, msg SayHelloMessage) error {
		log.Println("Received message:", msg.Payload)
		return nil
	})
	defer ctrl.UnsubscribeFromReceiveHelloOperation(context.Background())

	if err != nil {
		panic(err)
	}

	// Listen on port to let know that app is ready
	examples.ListenLocalPort(1234)
}
