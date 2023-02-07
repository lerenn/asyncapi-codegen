// Universal parts generation
//go:generate go run ../../../cmd/asyncapi-codegen -g client -p generated -i ../asyncapi.yaml -o ./generated/client.gen.go
//go:generate go run ../../../cmd/asyncapi-codegen -g broker -p generated -i ../asyncapi.yaml -o ./generated/broker.gen.go
//go:generate go run ../../../cmd/asyncapi-codegen -g types -p generated -i ../asyncapi.yaml -o ./generated/types.gen.go

// Specific brokers implementations generation
//go:generate go run ../../../cmd/asyncapi-codegen -g nats -p generated -i ../asyncapi.yaml -o ./generated/nats.gen.go

package main

import (
	"log"
	"time"

	"github.com/lerenn/asyncapi-codegen/examples/ping/client/generated"
	"github.com/nats-io/nats.go"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func main() {
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}

	// Create a new client controller
	ctrl := generated.NewClientController(generated.NewNATSController(nc))

	// Make a new ping message
	req := generated.NewPingMessage()
	req.Payload = "ping"

	// Create the publication function to send the message
	publicationFunc := func() error {
		log.Println("New ping request")
		return ctrl.PublishPing(req)
	}

	// The following function will subscribe to the 'pong' channel, execute the publication
	// function and wait for a response. The response will be detected through its
	// correlation ID.
	//
	// This function is available only if the 'correlationId' field has been filled
	// for any channel in the AsyncAPI specification. You will then be able to use it
	// with the form WaitForXXX where XXX is the channel name.
	resp, err := ctrl.WaitForPong(req.Headers.CorrelationID, publicationFunc, time.Second)
	if err != nil {
		panic(err)
	}

	log.Println("Got response:", resp.Payload)

	time.Sleep(time.Second)
}
