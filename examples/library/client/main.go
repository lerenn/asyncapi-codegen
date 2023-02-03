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

	"github.com/lerenn/asyncapi-codegen/examples/library/client/generated"
	"github.com/nats-io/nats.go"
)

func BooksListResponseHandler(msg generated.BooksListResponseMessage) {
	log.Println("received books list response!")
}

func main() {
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}

	// Create a new application controller
	clientController := generated.NewClientController(generated.NewNATSController(nc))

	// Subscribe specifically to the response handler
	if err := clientController.SubscribeBooksListResponse(BooksListResponseHandler); err != nil {
		panic(err)
	}

	// Make a new book list request
	log.Println("New book list request!")
	var req generated.BooksListRequestMessage
	req.Payload.Genre = "famous"
	if err := clientController.PublishBooksListRequest(req); err != nil {
		panic(err)
	}

	time.Sleep(time.Second)
}
