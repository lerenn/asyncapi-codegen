//go:generate go run ../../../../../cmd/asyncapi-codegen -g application,types -p main -i ../../asyncapi.yaml -o ./app.gen.go

package main

import (
	"context"
	"log"

	"github.com/lerenn/asyncapi-codegen/examples"
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

	// Create a new broker adapter
	broker, err := nats.NewController(addr, nats.WithQueueGroup("helloworld-apps"))
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

	// Subscribe to HelloWorld messages
	// Note: it will indefinitely wait for messages as context has no timeout
	log.Println("Subscribe to hello world...")
	ctrl.SubscribeHello(context.Background(), func(_ context.Context, msg HelloMessage) error {
		log.Println("Received message:", msg.Payload)
		return nil
	})

	// Listen on port to let know that app is ready
	examples.ListenLocalPort(1234)
}
