//go:generate go run ../../../../cmd/asyncapi-codegen -g application,types -p main -i ../../asyncapi.yaml -o ./app.gen.go

package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
)

func main() {
	// Create a new application controller
	ctrl, err := NewAppController(nats.NewController("nats://nats:4222"))
	if err != nil {
		panic(err)
	}
	defer ctrl.Close(context.Background())

	// Subscribe to HelloWorld messages
	// Note: it will indefinitely wait for messages as context has no timeout
	log.Println("Subscribe to hello world...")
	ctrl.SubscribeHello(context.Background(), func(_ context.Context, msg HelloMessage, _ bool) {
		log.Println("Received message:", msg.Payload)
	})

	// Process messages until interruption signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
