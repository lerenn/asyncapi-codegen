package ci

import (
	"os"
	"strings"

	"dagger.io/dagger"
)

const (
	examplesPath = "./examples/"
)

// Examples returns a container that runs all examples.
func Examples(client *dagger.Client, brokers map[string]*dagger.Service) map[string]*dagger.Container {
	containers := make(map[string]*dagger.Container, 0)

	// Set examples
	for _, p := range examplesPaths() {
		// Get corresponding broker
		brokerName := strings.Split(p, "/")[1]

		// Set app container
		app := client.Container().
			From(GolangImage).
			// Add source code as work directory
			With(sourceAsWorkdir(client)).
			// Set broker as dependency
			WithServiceBinding(brokerName, brokers[brokerName]).
			// Execute command
			WithExec([]string{"go", "run", examplesPath + p + "/app"}).
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
			WithExec([]string{"go", "run", examplesPath + p + "/user"})

		// Add user containers to containers
		containers[p] = user
	}

	return containers
}

func examplesPaths() []string {
	paths := make([]string, 0)

	examples, err := os.ReadDir(examplesPath)
	if err != nil {
		panic(err)
	}

	for _, e := range examples {
		if !e.Type().IsDir() {
			continue
		}

		brokers, err := os.ReadDir("./examples/" + e.Name())
		if err != nil {
			panic(err)
		}

		for _, b := range brokers {
			if !b.Type().IsDir() {
				continue
			}

			paths = append(paths, e.Name()+"/"+b.Name())
		}
	}

	return paths
}
