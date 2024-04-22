# How to contribute

Welcome on the CONTRIBUTING page of asyncapi-codegen ! Thanks for your interest
in contributing to this project, any contribution/suggestion is warmly welcomed.

If you need further guidance, you can find our team on the following:
* [Discussions on Github](https://github.com/lerenn/asyncapi-codegen/discussions)

## Contributing

### 1. Create an issue

If you have a suggestion, a bug, or any kind of thing that should impact the
source code, please open an issue on the
[Github repository](https://github.com/lerenn/asyncapi-codegen/issues).

Expose your problem and/or the suggestion and feel free to mention if you are
willing to do the change by yourself or prefer to let someone else do it.

It is important to open an issue in order to discuss the matter and avoid any
unecessary back-and-forth exchanges that would result in uneeded change in
the code you would write. 

### 2. Write your contribution

#### Start the development environment

If and once you're ready to write the contribution, you can start the development
environment by typing this command:

```bash 
make dev/up
```

Once you're done, you can terminate your environment by typing this command:

```bash
make dev/down
```

#### Explore the source code

Code is separated between the following directories:
* `build/` contains the files related to the CI and the packaging (Docker, ...)
  of asyncapi-codegen project
* `cmd/` contains all executable produces by this project (no the project
  internal tools)
* `examples/` contains workable examples of generated code along to their asyncapi
  document against the brokers supported in this project. You can test them if
  you start the `app` in a terminal then the `user` in another by using
  `go run ./examples/<example>/v<version>/<broker>/<app|user>`
* `pkg/`contains all the project source code that can be reused in other projects:
    * `pkg/asyncapi` contains the Go code for asyncapi specification
    * `pkg/ci` contains the Go code for the CI
    * `pkg/codegen` contains the Go code for the code generation
    * `pkg/extensions` contains the Go code for the extensions used by asyncapi-codegen users
    * `pkg/utils` contains the Go code for the utilities used by asyncapi-codegen
* `tests/` contains the tests of the project by version and by type (issue, etc).
  This is where you should implement your tests linked to your issues if this implies
  code generation.
* `tools/` contains the tools used by the project (like the certs generation tool
  for testing)

#### Write your code

Once you've explored the source code, you can start writing your code.

Please follow the following rules:
* Write tests for your code if possible
* Respect the code style (linter)

#### Test your code

##### Without code generation

If you don't need code generation (testing broker implementation, asyncapi
parsing, etc), feel free to add tests close to your changes in a `*_test.go` file.

##### With code generation

If you're code implies some code generation, you can write test in the corresponding
directory `./test/<version>/issue<#>/` where you can put the following files:
* `asyncapi.yml` the asyncapi document that will be used to generate the code
* `asyncapi.gen.go` the generated code
* `suite_test.go` the test file that will be used to test the generated code, it
  should have a `//go:generate` command to generate the `asyncapi.gen.go`.

Please respect the following rules:
* Channels should be prefixed with the version and the issue number (example:
  `v2.issue1`) to avoid collision in case of parallel tests
* The test package should be named `issue<#>` where `#` is the issue number
* Use the testify framework to write your tests (see the existing tests for
  examples)
* Brokers from `test/brokers.go` should be used to ensure that tests works with
  all brokers.

Of course, do not hesitate to ask for help if you need it.

### 3. Open a pull request

Once you're done with your code, you can open a pull request on the
[Github repository](https://github.com/lerenn/asyncapi-codegen/pulls).

## Missing information here?

If you think that some information is missing here, feel free to open an issue
on the [Github repository](https://github.com/lerenn/asyncapi-codegen/issues).