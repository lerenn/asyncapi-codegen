# HelloWorld example (AsyncAPI v2)

This example will use the AsyncAPI official example of the
[HelloWorld](https://www.asyncapi.com/docs/tutorials/getting-started/hello-world).

The code for this example have already been generated and can be
[read here](./examples/helloworld/v2), in the subdirectories `app/`
and `user/`. You can execute examples with `make examples`, or just one with
`EXAMPLE=<example> make examples` where the example is `<example>/<broker>`
(here `EXAMPLE=helloworld/nats`).

In order to recreate the code for user and application, you have to run this command:

```shell
# Install the tool
go install github.com/TheSadlig/asyncapi-codegen/cmd/asyncapi-codegen@latest

# Generate the code from the asyncapi file
asyncapi-codegen -i examples/helloworld/v2/asyncapi.yaml -p main -o ./helloworld.gen.go
```

We can then go through the `helloworld.gen.go` file to understand what will be used.

#### Application

Here is the code that is generated for the application side, with corresponding
comments:

```go
// AppController is the struct that you will need in order to interact with the
// event broker from the application side. You will generate this with the
// NewAppController function below.
type AppController struct

// NewAppController will create a new App Controller and will connect the
// BrokerController that you pass in argument to subscription and publication method.
func NewAppController(bs BrokerController, options ...ControllerOption) *AppController

// Close function will clean up all resources and subscriptions left in the
// application controller. This should be call right after NewAppController
// with a `defer`
func (ac *AppController) Close(ctx context.Context)

// SubscribeAll will subscribe to all channel that the application should listen to.
//
// In order to use it, you'll have to implement the AppSubscriber interface and
// pass it as an argument to this function. Thus, the subscription will automatically
// call the corresponding function when it will receive a message.
//
// In the HelloWorld example, only one function will listen on application side,
// making it a bit overkill. You can directly use the SubscribeHello method.
func (ac *AppController) SubscribeAll(ctx context.Context, as AppSubscriber) error

// UnsubscribeAll will unsubscribe all channel that have subscribed to through
// SubscribeAll or SubscribeXXX where XXX correspond to the channel name.
func (ac *AppController) UnsubscribeAll(ctx context.Context)

// SubscribeHello will subscribe to new messages on the "hello" channel.
// It will expect messages as specified in the AsyncAPI specification.
//
// You just have to give a function that match the signature of the callback and
// then process the received message.
//
// The `done` argument is true when the subscription is closed. It can be used to
// cleanup resources, such as channels.
//
// The subscription will be canceled if the context is canceled, if the subscription
// is explicitely unsubscribed or if the controller is closed
func (ac *AppController) SubscribeHello(ctx context.Context, fn func(msg HelloMessage)) error

// UnsubscribeHello will unsubscribe only the subscription on the "hello" channel.
// It should be only used when wanting specifically that, otherwise the clean up
// will be handled by the Close function.
func (ac *AppController) UnsubscribeHello(ctx context.Context)
```

And here is an example of the application that could be written to use this generated
code (you can also find it [here](./examples/helloworld/v2)):

```go
import(
  "github.com/TheSadlig/asyncapi-codegen/pkg/extensions/brokers/nats"
  // ...
)

func main() {
  // Create a NATS controller
  broker, _ := nats.NewController("nats://nats:4222")
  defer broker.Close()

  // Create a new application controller
  ctrl, _ := NewAppController(broker)
  defer ctrl.Close(context.Background())

  // Subscribe to HelloWorld messages
  // Note: it will indefinitely wait for messages as context has no timeout
  log.Println("Subscribe to hello world...")
  ctrl.SubscribeHello(context.Background(), func(_ context.Context, msg HelloMessage) {
    log.Println("Received message:", msg.Payload)
  })

  // Process messages until interruption signal
  // ...
}
```

#### User

Here is the code that is generated for the user side, with corresponding
comments:

```go
// UserController is the struct that you will need in order to interact with the
// event broker from the user side. You will generate this with the
// NewUserController function below.
type UserController struct

// NewUserController will create a new User Controller and will connect the
// BrokerController that you pass in argument to subscription and publication method.
func NewUserController(bs BrokerController, options ...ControllerOption) *UserController

// Close function will clean up all resources and subscriptions left in the
// application controller. This should be call right after NewAppController
// with a `defer`
func (cc *UserController) Close(ctx context.Context)

// PublishHello will publish a hello world message on the "hello" channel as
// specified in the AsyncAPI specification.
func (cc *UserController) PublishHello(ctx context.Context, msg HelloMessage) error
```

And here is an example of the user that could be written to use this generated
code (you can also find it [here](./examples/helloworld/v2)):

```go
import(
  "github.com/TheSadlig/asyncapi-codegen/pkg/extensions/brokers/nats"
  // ...
)

func main() {
  // Create a NATS controller
  broker, _ := nats.NewController("nats://nats:4222")
  defer broker.Close()

  // Create a new user controller
  ctrl, _ := NewUserController(broker)
  defer ctrl.Close(context.Background())

  // Send HelloWorld
  log.Println("Publishing 'hello world' message")
  ctrl.PublishHello(context.Background(), HelloMessage{Payload: "HelloWorld!"})

  // ...
}
```

#### Types

According to the specification that you pass in parameter, some types will also
be  Here is the ones generated for the HelloWorld example:

```go
// HelloMessage will contain all the information that will be sent on the 'hello'
// channel. There is only a payload here, but you could find also headers,
// correlation id, and more.
type HelloMessage struct {
  Payload string
}
```