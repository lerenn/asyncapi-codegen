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

	"github.com/google/uuid"
	"github.com/lerenn/asyncapi-codegen/examples/library/client/generated"
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

	// Create a new application controller
	clientController := generated.NewClientController(generated.NewNATSController(nc))

	// Make a new book list request
	var req generated.BooksListRequestMessage
	req.Payload.Genre = "famous"
	req.Headers.CorrelationID = uuid.New().String()

	// Create the publication function
	publicationFunc := func() error {
		log.Println("New book list request for:", req.Payload.Genre)
		return clientController.PublishBooksListRequest(req)
	}

	// Send request and wait for response
	resp, err := clientController.WaitForBooksListResponse(req.Headers.CorrelationID, publicationFunc, time.Second)
	if err != nil {
		panic(err)
	}

	log.Println("Got books:", resp.Payload.Books)

	time.Sleep(time.Second)
}
