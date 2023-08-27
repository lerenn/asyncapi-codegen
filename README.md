# AsyncAPI Codegen

Generate Go client and server boilerplate from AsyncAPI specifications.

**⚠️ Do not hesitate raise an issue on any bug or missing feature.**
**Contributions are welcomed!**

*Inspired from popular [deepmap/oapi-codegen](https://github.com/deepmap/oapi-codegen)*

## Contents

* [Supported functionalities](#supported-functionalities)
* [Concepts](#concepts)
* [Examples](#examples):
  * [Basic example](#basic-example)
  * [Request/Response example](#request-response-example)
* [CLI options](#cli-options)
* [Advanced topics](#advanced-topics)
* [Contributing and support](#contributing-and-support)

## Supported functionalities

* AsyncAPI versions:
  * 2.6.0
* Brokers:
  * NATS
  * Custom (implementation specified by the developer)
  * *Open a ticket for any missing one that you would want to have here!*
* Formats:
  * JSON
* Logging:
  * JSON (ECS compatible)

## Concepts

![basic schema](assets/basic-schema.svg)

Let's imagine a message broker centric architecture: you have the application
that you are developing on the right and the potential client(s) on the left.

Being a two directional communication, both of them can communicate to each
other through the broker. They can even communicate with themselves, in case
of multiple clients or application replication.

For more information about this, please refere to the [official AsyncAPI
concepts](https://www.asyncapi.com/docs/concepts).

### With Async API generated code

![with codegen schema](assets/with-codegen-schema.svg)

* <span style="color:yellow">Yellow parts</span>: when using the codegen tool,
you will generate the code that will act as an adapter (or controller) between
the client, the broker, and the application.
* <span style="color:red">Red parts</span>: you will need to fill these parts
between client, broker and application. These will allow message production and
reception with the generated code.
* <span style="color:orange">Orange parts</span>: these will be also generated
automatically if you use an implemented message broker. You can also use the
`none` type in order to implement it yourself.

## Examples

Here is a list of example, from basic to advanced ones.

### Basic example

This example will use the AsyncAPI official example of the
[HelloWorld](https://www.asyncapi.com/docs/tutorials/getting-started/hello-world).

> The code for this example have already been generated and can be
[read here](./examples/helloworld/), in the subdirectories `app/generated/`
and `client/generated/`. You can execute the example with `docker-compose up`.

In order to recreate the code for client and application, you have to run this command:

```shell
# Install the tool
go install github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@latest

# Generate the code from the asyncapi file
asyncapi-codegen -i examples/helloworld/asyncapi.yaml -o ./helloworld.gen.go
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
func NewAppController(bs BrokerController) *AppController

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
func (ac *AppController) SubscribeHello(ctx context.Context, fn func(msg HelloMessage, done bool)) error

// UnsubscribeHello will unsubscribe only the subscription on the "hello" channel.
// It should be only used when wanting specifically that, otherwise the clean up
// will be handled by the Close function.
func (ac *AppController) UnsubscribeHello(ctx context.Context)
```

And here is an example of the application that could be written to use this generated
code with NATS (you can also find it [here](./examples/helloworld/app/main.go)):

```go
// Connect to NATS
nc, _ := nats.Connect("nats://nats:4222")

// Create a new application controller
ctrl, _ := generated.NewAppController(generated.NewNATSController(nc))
defer ctrl.Close(context.Background())

// Subscribe to HelloWorld messages
// Note: it will indefinitely wait for messages as context has no timeout
log.Println("Subscribe to hello world...")
ctrl.SubscribeHello(context.Background(), func(_ context.Context, msg generated.HelloMessage, _ bool) {
  log.Println("Received message:", msg.Payload)
})

// Process messages until interruption signal
/* ... */
```

#### Client

Here is the code that is generated for the client side, with corresponding
comments:

```go
// ClientController is the struct that you will need in order to interact with the
// event broker from the client side. You will generate this with the 
// NewClientController function below.
type ClientController struct

// NewClientController will create a new Client Controller and will connect the
// BrokerController that you pass in argument to subscription and publication method.
func NewClientController(bs BrokerController) *ClientController

// Close function will clean up all resources and subscriptions left in the
// application controller. This should be call right after NewAppController
// with a `defer`
func (cc *ClientController) Close(ctx context.Context)

// PublishHello will publish a hello world message on the "hello" channel as
// specified in the AsyncAPI specification.
func (cc *ClientController) PublishHello(ctx context.Context, msg HelloMessage) error
```

And here is an example of the client that could be written to use this generated
code with NATS (you can also find it [here](./examples/helloworld/app/main.go)):

```go
// Connect to NATS
nc, _ := nats.Connect("nats://nats:4222")

// Create a new application controller
ctrl, _ := generated.NewClientController(generated.NewNATSController(nc))
defer ctrl.Close(context.Background())

// Send HelloWorld
log.Println("Publishing 'hello world' message")
ctrl.PublishHello(context.Background(), generated.HelloMessage{Payload: "HelloWorld!"})
```

#### Types

According to the specification that you pass in parameter, some types will also
be generated. Here is the ones generated for the HelloWorld example:

```go
// HelloMessage will contain all the information that will be sent on the 'hello'
// channel. There is only a payload here, but you could find also headers,
// correlation id, and more.
type HelloMessage struct {
	Payload string
}
```

#### Broker

In order to connect your application and your client to your broker, we need to
provide an adapter to it. Here is the interface that you need to satisfy:

```go
type BrokerController interface {
  	// Publish will be called under the hood by any PublishXXX function
	Publish(ctx context.Context, channel string, mw UniversalMessage) error
  	// Subscribe will be called under the hood by any SubscribeXXX function
	Subscribe(ctx context.Context, channel string) (msgs chan UniversalMessage, stop chan interface{}, err error)
}
```

You can find that there is an `UniversalMessage` structure that is provided and
that aims to abstract the event broker technology.

You can either generate an already existing adapter, or write your own if it doesn't
exists or if it doesn't suit your needs.

### Request/Response example

This example will use a `ping` example that you can find 
[here](./examples/ping/asyncapi.yaml).

> The code for this example have already been generated and can be
[read here](./examples/ping/), in the subdirectories `server/generated/`
and `client/generated/`. You can execute the example with `docker-compose up`.

In order to recreate the code for client and application, you have to run this command:

```shell
# Install the tool
go install github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@latest

# Generate the code from the asyncapi file
asyncapi-codegen -i examples/ping/asyncapi.yaml -o ./ping.gen.go
```

We can then go through the possible application and client implementations that
use `ping.gen.go`. 

#### Application (or server in this case)

```golang
type ServerSubscriber struct {
	Controller *generated.AppController
}

func (s ServerSubscriber) Ping(req generated.PingMessage, _ bool) {
	// Generate a pong message, set as a response of the request
	resp := generated.NewPongMessage()
	resp.SetAsResponseFrom(req)
	resp.Payload.Message = "pong"
	resp.Payload.Time = time.Now()

	// Publish the pong message
	s.Controller.PublishPong(cresp)
}

func main() {
	/* ... */

	// Create a new server controller
	ctrl, _ := generated.NewAppController(generated.NewNATSController(nc))
	defer ctrl.Close(context.Background())

	// Subscribe to all (we could also have just listened on the ping request channel)
	sub := ServerSubscriber{Controller: ctrl}
	ctrl.SubscribeAll(context.Background(), sub)

	// Process messages until interruption signal
	/* ... */
}
```

#### Client

```golang
// Create a new client controller
ctrl, _ := generated.NewClientController(/* Add corresponding broker controller */)
defer ctrl.Close(context.Background())

// Make a new ping message
req := generated.NewPingMessage()
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
resp, _ := ctrl.WaitForPong(context.Background(), req, publicationFunc)
```

## CLI options

The default options for oapi-codegen will generate everything; client, application,
broker, type definitions, and broker implementations but you can generate subsets
of those via the -generate flag. It defaults to client,application,broker,types
but you can specify any combination of those.

Here are the universal parts that you can generate:

* `application`: generate the application boilerplate. `application` requires
  the types in the same package to compile.
* `client`: generate the client boilerplate. It, too, requires the types to be
  present in its package.
* `broker`: generate the broker controller that you have to fill either with an
  existing implementation (more below), or by implementing your own.
* `types`: all type definitions for all types in the AsyncAPI spec.
  This will be everything under `#components`, as well as request parameter,
  request body, and response type objects.

You can also specify some specific implementation for the broker of your choice:

* `nats`: generate the NATS message broker boilerplate.

## Advanced topics

### Middlewares

You can use middlewares that will be executing when receiving and publishing
messages. You can add one or multiple middlewares using the following function
on a controller:

```golang
// Create a new app controller with a NATS controller for example
ctrl, _ := generated.NewAppController(generated.NewNATSController(nc))

// Add middleware
ctrl.AddMiddlewares(myMiddleware1, myMiddleware2 /*, ... */)
```

Here the function signature that should be satisfied:

```golang
func(ctx context.Context, next middleware.Next) context.Context
```

**Note:** the returned context will be the one that will be passed to following
middlewares, and finally to the generated code (and subscription callback).

#### Filtering messages

If you want to target specific messages, you can use the context passed in argument:

```golang
import(
	apiContext "github.com/lerenn/asyncapi-codegen/pkg/context"
)

func myMiddleware(ctx context.Context, _ middleware.Next) context.Context {
	// Execute this middleware only if this is a received message
	apiContext.IfEquals(ctx, apiContext.KeyIsDirection, "reception", func() {
		// Do specific stuff if message is received
	})

	return ctx
}
```

You can even discriminate on more specification. Please see the [Context section](#context).

#### Executing code after receiving/publishing the message

By default, middlewares will be executed right before the operation. If there is
a need to execute code before and/or after the operation, you can call the `next`
argument that represents the next middleware that should be executed or the
operation corresponding code if this was the last middleware.

Here is an example:

```golang
func surroundingMiddleware(ctx context.Context, next middleware.Next) context.Context {
	// Pre-operation
	fmt.Println("This will be displayed BEFORE the reception/publication")

	// Calling next middleware or reception/publication code
	// The given context will be the one propagated to other middlewares and operation source code
	next(ctx)

	// Post-operation
	fmt.Println("This will be displayed AFTER the reception/publication")

	return ctx
}
```

### Context

When receiving the context from generated code (either in subscription,
middleware, logging, etc), you can get some information embedded in context.

To get these information, please use the functions from
`github.com/lerenn/asyncapi-codegen/pkg/context`:

```golang
// Execute this middleware only if this is from "ping" channel
apiContext.IfEquals(ctx, apiContext.KeyIsChannel, "ping", func() {
	// Do specific stuff if the channel is ping
})
```

You can find other keys in the package `pkg/context`.

### Logging

#### Enable logging

To enable logging, the only thing you have to do is to set a logger to your
controller with the function `SetLogger()`:

```golang
// Create a new app controller with a NATS controller for example
ctrl, _ := generated.NewAppController(generated.NewNATSController(nc))
	
// Attach a logger (optional)
// You can find loggers in `github.com/lerenn/asyncapi-codegen/pkg/log` or create your own
logger := log.NewECS()
ctrl.SetLogger(logger)
```

You can find all loggers in the directory `pkg/log`.

#### Custom logging

It is possible to set your own logger to the generated code, all you have to do
is to fill the following interface:

```golang
type Logger interface {
    // Info logs information based on a message and key-value elements
    Info(ctx log.Context, msg string, info ...log.AdditionalInfo)

    // Error logs error based on a message and key-value elements
    Error(ctx log.Context, msg string, info ...log.AdditionalInfo)
}
```

Here is a basic implementation example:

```golang
type SimpleLogger struct{}

func (logger SimpleLogger) formatLog(ctx log.Context, info ...log.AdditionalInfo) string {
	var formattedLogInfo string
	for i := 0; i < len(keyvals)-1; i += 2 {
		formattedLogInfo = fmt.Sprintf("%s, %s: %+v", formattedLogInfo, info.Key, info.Value)
	}
	return fmt.Sprintf("%s, context: %+v", formattedLogInfo, ctx)
}

func (logger SimpleLogger) Info(ctx log.Context, msg string, info ...log.AdditionalInfo) {
	log.Printf("INFO: %s%s", msg, logger.formatLog(ctx, info...))
}

func (logger SimpleLogger) Error(ctx log.Context, msg string, info ...log.AdditionalInfo) {
	log.Printf("ERROR: %s%s", msg, logger.formatLog(ctx, info...))
}
```

You can then create a controller with a logger using similar lines:

```golang
// Create a new app controller with a NATS controller for example
ctrl, _ := generated.NewAppController(generated.NewNATSController(nc))

// Set a logger
ctrl.SetLogger(SimpleLogger{})
```

### Use of queue groups and queue name customization

Queues are used under the hood by implemented brokers to have replicates of a
same service that automatically load-balance request between instances. By
default, the queue name used by implemented brokers is `asyncapi`.

However, it is possible to set a custom queue name if you want to have multiple
applications which uses code generated by `asyncapi-codegen` but on different
queues:

```golang
broker.SetQueueName("my-custom-queue-name")
```

## Contributing and support

If you find any bug or lacking a feature, please raise an issue on the Github repository!

Also please do not hesitate to propose any improvment or bug fix on PR.
Any contribution is warmly welcomed!

And if you find this project useful, please support it through the Support feature
on Github.
