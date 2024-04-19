package ci

import (
	"dagger.io/dagger"
)

// Tests returns containers for all tests.
func Tests(client *dagger.Client, brokers map[string]*dagger.Service, path string) *dagger.Container {
	return client.Container().
		// Add base image
		From(GolangImage).
		// Add source code as work directory
		With(sourceAsWorkdir(client)).
		// Set brokers as dependencies of app and user
		With(BindBrokers(brokers)).
		// Execute command
		WithExec([]string{"go", "test", path})
}
