// Universal parts generation
//go:generate go run ../../../cmd/asyncapi-codegen -g client -p generated -i ../asyncapi.yaml -o ./generated/client.gen.go
//go:generate go run ../../../cmd/asyncapi-codegen -g broker -p generated -i ../asyncapi.yaml -o ./generated/broker.gen.go
//go:generate go run ../../../cmd/asyncapi-codegen -g types -p generated -i ../asyncapi.yaml -o ./generated/types.gen.go

// Specific brokers implementations generation
//go:generate go run ../../../cmd/asyncapi-codegen -g nats -p generated -i ../asyncapi.yaml -o ./generated/nats.gen.go

package main

import (
	"log"

	"github.com/lerenn/asyncapi-codegen/examples/helloworld/client/generated"
	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}

	// Create a new client controller
	ctrl, err := generated.NewClientController(generated.NewNATSController(nc))
	if err != nil {
		panic(err)
	}
	defer ctrl.Close()

	// Send HelloWorld
	log.Println("Publishing 'hello world' message")
	if err := ctrl.PublishHello(generated.HelloMessage{
		Payload: "HelloWorld!",
	}); err != nil {
		panic(err)
	}
}
