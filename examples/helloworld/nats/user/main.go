//go:generate go run ../../../../cmd/asyncapi-codegen -g user,types -p main -i ../../asyncapi.yaml -o ./user.gen.go

package main

import (
	"context"
	"log"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers"
	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}

	// Create a new user controller
	ctrl, err := NewUserController(brokers.NewNATSController(nc))
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
