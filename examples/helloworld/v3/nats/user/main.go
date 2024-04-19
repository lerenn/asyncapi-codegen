//go:generate go run ../../../../../cmd/asyncapi-codegen -g user,types -p main -i ../../asyncapi.yaml -o ./user.gen.go

package main

import (
	"context"
	"log"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions/brokers/nats"
	testutil "github.com/lerenn/asyncapi-codegen/pkg/utils/test"
)

func main() {
	// Get broker address based on the environment, it will returns an address like "nats://nats:4222"
	// Note: this is not needed in your application, you can directly use the address
	addr := testutil.BrokerAddress(testutil.BrokerAddressParams{
		Schema:         "nats",
		DockerizedAddr: "nats",
		Port:           "4222",
	})

	// Create a new broker
	broker, err := nats.NewController(addr, nats.WithQueueGroup("helloworld-users"))
	if err != nil {
		panic(err)
	}
	defer broker.Close()

	// Create a new user controller
	ctrl, err := NewUserController(broker)
	if err != nil {
		panic(err)
	}
	defer ctrl.Close(context.Background())

	// Send HelloWorld
	// Note: it will indefinitely wait to publish as context has no timeout
	log.Println("Publishing 'hello world' message")
	if err := ctrl.SendToReceiveHelloOperation(context.Background(), SayHelloMessage{
		Payload: "HelloWorld!",
	}); err != nil {
		panic(err)
	}
}
