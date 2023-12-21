package pipeline

import (
	"os"

	"dagger.io/dagger"
)

// Tests returns containers for all tests.
func Tests(client *dagger.Client, brokers map[string]*dagger.Service) []*dagger.Container {
	containers := make([]*dagger.Container, 0)

	// Set examples
	for _, p := range testsPaths() {
		t := client.Container().
			// Add base image
			From(GolangImage).
			// Add source code as work directory
			With(sourceAsWorkdir(client)).
			// Set brokers as dependencies of app and user
			With(BindBrokers(brokers)).
			// Execute command
			WithExec([]string{"go", "test", p})

		// Add user containers to containers
		containers = append(containers, t)
	}

	return containers
}

func testsPaths() []string {
	paths := make([]string, 0)

	test, err := os.ReadDir("./test/issues")
	if err != nil {
		panic(err)
	}

	for _, t := range test {
		if !t.Type().IsDir() {
			continue
		}

		paths = append(paths, "./test/issues/"+t.Name())
	}

	return paths
}
