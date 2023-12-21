package pipeline

import (
	"os"

	"dagger.io/dagger"
)

// Examples returns a container that runs all examples.
func Examples(client *dagger.Client, brokers map[string]*dagger.Service) []*dagger.Container {
	containers := make([]*dagger.Container, 0)

	// Set examples
	for _, p := range examplesPaths() {
		// Set app container
		app := client.Container().
			From(GolangImage).
			// Add source code as work directory
			With(sourceAsWorkdir(client)).
			// Add brokers as dependencies
			With(BindBrokers(brokers)).
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
			// Set brokers as dependencies of app and user
			With(BindBrokers(brokers)).
			// Add app as dependency of user
			WithServiceBinding("app", app).
			// Execute command
			WithExec([]string{"go", "run", p + "/user"})

		// Add user containers to containers
		containers = append(containers, user)
	}

	return containers
}

func examplesPaths() []string {
	paths := make([]string, 0)

	examples, err := os.ReadDir("./examples")
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

			paths = append(paths, "./examples/"+e.Name()+"/"+b.Name())
		}
	}

	return paths
}
