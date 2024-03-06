# Ping example (AsyncAPI v2)

This example will use a `ping` example that you can find
[here](./examples/ping/v2).

> The code for this example have already been generated and can be
read, in the subdirectories `app/` and `user/`. You can execute examples with
`make examples`, or just one with `EXAMPLE=<example> make examples` where the
example is `<example>/<broker>` (here `EXAMPLE=ping/nats`).

In order to recreate the code for user and application, you have to run this command:

```shell
# Install the tool
go install github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@latest

# Generate the code from the asyncapi file
asyncapi-codegen -i examples/ping/v2/asyncapi.yaml -p main -o ./ping.gen.go
```

We can then go through the possible application and user implementations that
use `ping.gen.go`.

#### Application

```golang
type Subscriber struct {
  Controller *AppController
}

func (s Subscriber) Ping(req PingMessage) {
  // Generate a pong message, set as a response of the request
  resp := NewPongMessage()
  resp.SetAsResponseFrom(&req)
  resp.Payload.Message = "pong"
  resp.Payload.Time = time.Now()

  // Publish the pong message
  s.Controller.PublishPong(cresp)
}

func main() {
  // ...

  // Create a new application controller
  ctrl, _ := NewAppController(/* Add corresponding broker controller */)
  defer ctrl.Close(context.Background())

  // Subscribe to all (we could also have just listened on the ping request channel)
  sub := AppSubscriber{Controller: ctrl}
  ctrl.SubscribeAll(context.Background(), sub)

  // Process messages until interruption signal
  // ...
}
```

#### User

```golang
// Create a new user controller
ctrl, _ := NewUserController(/* Add corresponding broker controller */)
defer ctrl.Close(context.Background())

// Make a new ping message
req := NewPingMessage()
req.Payload = "ping"

// Create the publication function to send the message
publicationFunc := func(ctx context.Context) error {
  return ctrl.PublishPing(ctx, req)
}

// The following function will subscribe to the 'pong' channel, execute the publication
// function and wait for a response. The response will be detected through its
// correlation ID.
//
// This function is available only if the 'correlationId' field has been filled
// for any channel in the AsyncAPI specification. You will then be able to use it
// with the form WaitForXXX where XXX is the channel name.
resp, _ := ctrl.WaitForPong(context.Background(), &req, publicationFunc)
```
