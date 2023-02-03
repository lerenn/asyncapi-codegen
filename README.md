# AsyncAPI Client and Server Code Generator

Generate Go client and server boilerplate from AsyncAPI specifications.

*Inspired from popular [deepmap/oapi-codegen](https://github.com/deepmap/oapi-codegen)*

## Supported functionalities

* AsyncAPI versions:
  * 2.5.0
* Brokers:
  * Custom (implementation specified by the developer)
* Formats:
  * JSON

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

When using the codegen tool, you will generate the code, represented by the <span style="color:yellow">yellow parts</span>, that will act as an
adapter (or controller) between the client, the broker, and the application.

You will need to fill the <span style="color:red">red parts</span> between client, broker and application. These will allow message production and reception with the generated code.

The <span style="color:orange">orange parts</span> will be also generated automatically if you use an implemented
message broker. You can also use the `none` type in order to implement it yourself.

<!-- TODO: ## Example -->

## Using `asyncapi-codegen`

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
