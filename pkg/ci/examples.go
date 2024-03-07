package ci

import (
	"strings"

	"dagger.io/dagger"
)

// Examples returns a container that runs all examples.
func Examples(client *dagger.Client, brokers map[string]*dagger.Service) map[string]*dagger.Container {
	containers := make(map[string]*dagger.Container, 0)

	// Set examples
	for _, p := range directoriesAtSublevel(2, "./examples") {
		// Get corresponding broker
		brokerName := strings.Split(p, "/")[4]

		// Set app container
		app := client.Container().
			From(GolangImage).
			// Add source code as work directory
			With(sourceAsWorkdir(client)).
			// Set broker as dependency
			WithServiceBinding(brokerName, brokers[brokerName]).
			// Execute command
			WithExec([]string{"go", "run", p + "/app"}).
			// Add exposed port to let know when the service is ready
			WithExposedPort(1234).
			// Set as service
			AsService()

		// Set user container
		user := client.Container().
			// Add base image
			From(GolangImage).
			// Add source code as work directory
			With(sourceAsWorkdir(client)).
			// Set broker as dependency
			WithServiceBinding(brokerName, brokers[brokerName]).
			// Add app as dependency of user
			WithServiceBinding("app", app).
			// Execute command
			WithExec([]string{"go", "run", p + "/user"})

		// Add user containers to containers
		containers[p] = user
	}

	return containers
}
