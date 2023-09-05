// Universal parts generation
//go:generate go run ../../../cmd/asyncapi-codegen -g client -p generated -i ../asyncapi.yaml -o ./generated/client.gen.go
//go:generate go run ../../../cmd/asyncapi-codegen -g types -p generated -i ../asyncapi.yaml -o ./generated/types.gen.go

package main

import (
	"context"
	"log"

	"github.com/lerenn/asyncapi-codegen/examples/helloworld/client/generated"
	"github.com/lerenn/asyncapi-codegen/pkg/broker/controllers"
	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}

	// Create a new client controller
	ctrl, err := generated.NewClientController(controllers.NewNATS(nc))
	if err != nil {
		panic(err)
	}
	defer ctrl.Close(context.Background())

	// Send HelloWorld
	// Note: it will indefinitely wait to publish as context has no timeout
	log.Println("Publishing 'hello world' message")
	if err := ctrl.PublishHello(context.Background(), generated.HelloMessage{
		Payload: "HelloWorld!",
	}); err != nil {
		panic(err)
	}
}
