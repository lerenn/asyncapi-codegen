# Ping example (AsyncAPI v3)

This example will use a `ping` example that you can find
[here](./examples/ping/v3).

> The code for this example have already been generated and can be
read, in the subdirectories `app/` and `user/`. You can execute examples with
`make examples`, or just one with `EXAMPLE=<example> make examples` where the
example is `<example>/<broker>` (here `EXAMPLE=ping/nats`).

In order to recreate the code for user and application, you have to run this command:

```shell
# Install the tool
go install github.com/TheSadlig/asyncapi-codegen/cmd/asyncapi-codegen@latest

# Generate the code from the asyncapi file
asyncapi-codegen -i examples/ping/v3/asyncapi.yaml -p main -o ./ping.gen.go
```

We can then go through the possible application and user implementations that
use `ping.gen.go`.

#### Application

```golang
type Subscriber struct {
  Controller *AppController
}

func (s Subscriber) PingRequestOperationReceived(ctx context.Context, ping PingMessage) {
  // Publish the pong message, with the callback function to modify it
  // Note: it will indefinitely wait to publish as context has no timeout
  err := s.Controller.ReplyToPingRequestOperation(ctx, ping, func(pong *PongMessage) {
  	// Reply with the same event than the ping
  	pong.Payload.Event = ping.Payload.Event
  })

  // ...
}

func main() {
  // ...

  // Create a new application controller
  ctrl, _ := NewAppController(/* Add corresponding broker controller */)
  defer ctrl.Close(context.Background())

  // Subscribe to all (we could also have just listened on the ping request channel)
  sub := ServerSubscriber{Controller: ctrl}
	if err := ctrl.SubscribeToAllChannels(context.Background(), sub); err != nil {
	  // -- Error management
	}
	defer ctrl.UnsubscribeFromAllChannels(context.Background())

  // Process messages until interruption signal
  // ...
}
```

#### User

```golang
// Create a new user controller
ctrl, err := NewUserController(/* Add corresponding broker controller */)
// -- Error management
defer ctrl.Close(context.Background())

// Make a new ping message
req := NewPingMessage()
// -- you can modifiy the request here

// The following function will subscribe to the 'pong' channel (reply channel
// to PingRequest operation) and wait for a response. The response will be
// detected through its correlation ID. However, if there is no correlation
// ID, then it will return the first message on the reply channel.
//
// Note: it will indefinitely wait for messages as context has no timeout
resp, err := ctrl.RequestToPingRequestOperation(context.Background(), req)
// -- Error management

// Use the response
```
