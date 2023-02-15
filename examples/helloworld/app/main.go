// Universal parts generation
//go:generate go run ../../../cmd/asyncapi-codegen -g application -p generated -i ../asyncapi.yaml -o ./generated/app.gen.go
//go:generate go run ../../../cmd/asyncapi-codegen -g broker -p generated -i ../asyncapi.yaml -o ./generated/broker.gen.go
//go:generate go run ../../../cmd/asyncapi-codegen -g types -p generated -i ../asyncapi.yaml -o ./generated/types.gen.go

// Specific brokers implementations generation
//go:generate go run ../../../cmd/asyncapi-codegen -g nats -p generated -i ../asyncapi.yaml -o ./generated/nats.gen.go

package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/lerenn/asyncapi-codegen/examples/helloworld/app/generated"
	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}

	// Create a new application controller
	ctrl, err := generated.NewAppController(generated.NewNATSController(nc))
	if err != nil {
		panic(err)
	}
	defer ctrl.Close()

	// Subscribe to HelloWorld messages
	log.Println("Subscribe to hello world...")
	ctrl.SubscribeHello(func(msg generated.HelloMessage) {
		log.Println("Received message:", msg.Payload)
	})

	// Process messages until interruption signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
