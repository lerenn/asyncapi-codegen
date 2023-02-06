// Universal parts generation
//go:generate go run ../../../cmd/asyncapi-codegen -g application -p generated -i ../asyncapi.yaml -o ./generated/app.gen.go
//go:generate go run ../../../cmd/asyncapi-codegen -g broker -p generated -i ../asyncapi.yaml -o ./generated/broker.gen.go
//go:generate go run ../../../cmd/asyncapi-codegen -g types -p generated -i ../asyncapi.yaml -o ./generated/types.gen.go

// Specific brokers implementations generation
//go:generate go run ../../../cmd/asyncapi-codegen -g nats -p generated -i ../asyncapi.yaml -o ./generated/nats.gen.go

package main

import (
	"log"

	"github.com/lerenn/asyncapi-codegen/examples/ping/server/generated"
	"github.com/nats-io/nats.go"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

type AppSubscriber struct {
	Controller *generated.AppController
}

func (as AppSubscriber) Ping(req generated.PingMessage) {
	var resp generated.PongMessage

	log.Println("Received a ping request")

	// Respond with pong
	resp.Payload = "pong"

	// And with same correlation Id
	resp.Headers.CorrelationID = req.Headers.CorrelationID

	err := as.Controller.PublishPong(resp)
	if err != nil {
		panic(err)
	}
}

func main() {
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}

	// Create a new application controller
	appController := generated.NewAppController(generated.NewNATSController(nc))

	// Subscribe to all
	log.Println("Subscribe to all...")
	sub := AppSubscriber{Controller: appController}
	if err := appController.SubscribeAll(sub); err != nil {
		panic(err)
	}

	// Listen
	log.Println("Listening to subscriptions...")
	irq := make(chan interface{})
	appController.Listen(irq)
}
