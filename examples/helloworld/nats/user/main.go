//go:generate go run ../../../../cmd/asyncapi-codegen -g user,types -p main -i ../../asyncapi.yaml -o ./user.gen.go

package main

import (
	"context"
	"log"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
)

func main() {
	// Create a new broker
	broker := nats.NewController("nats://nats:4222", nats.WithQueueGroup("helloworld-users"))

	// Create a new user controller
	ctrl, err := NewUserController(broker)
	if err != nil {
		panic(err)
	}
	defer ctrl.Close(context.Background())

	// Send HelloWorld
	// Note: it will indefinitely wait to publish as context has no timeout
	log.Println("Publishing 'hello world' message")
	if err := ctrl.PublishHello(context.Background(), HelloMessage{
		Payload: "HelloWorld!",
	}); err != nil {
		panic(err)
	}
}
