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
	"time"

	"github.com/lerenn/asyncapi-codegen/examples/ping/server/generated"
	"github.com/nats-io/nats.go"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

type ServerSubscriber struct {
	Controller *generated.AppController
}

func (s ServerSubscriber) Ping(req generated.PingMessage, _ bool) {
	log.Println("Received a ping request")

	// Generate a pong message, set as a response of the request
	resp := generated.NewPongMessage()
	resp.SetAsResponseFrom(req)
	resp.Payload.Message = "pong"
	resp.Payload.Time = time.Now()

	// Publish the pong message
	err := s.Controller.PublishPong(resp)
	if err != nil {
		panic(err)
	}
}

func main() {
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}

	// Create a new server controller
	ctrl, err := generated.NewAppController(generated.NewNATSController(nc))
	if err != nil {
		panic(err)
	}
	defer ctrl.Close()

	// Attach a logger (optional)
	// ctrl.SetLogger(SimpleLogger{})

	// Subscribe to all (we could also have just listened on the ping request channel)
	log.Println("Subscribe to all...")
	sub := ServerSubscriber{Controller: ctrl}
	if err := ctrl.SubscribeAll(sub); err != nil {
		panic(err)
	}

	// Process messages until interruption signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
