package pipeline

import (
	"os"

	"dagger.io/dagger"
)

// Tests returns containers for all tests.
func Tests(client *dagger.Client, brokers map[string]*dagger.Service) map[string]*dagger.Container {
	containers := make(map[string]*dagger.Container, 0)

	// Set examples
	for _, p := range testsPaths(3, "./test") {
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
		containers[p] = t
	}

	return containers
}

func testsPaths(level int, path string) []string {
	paths := make([]string, 0)

	tests, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	if level == 0 {
		for _, t := range tests {
			if !t.Type().IsDir() {
				continue
			}

			paths = append(paths, path+"/"+t.Name())
		}

		return paths
	}

	for _, t := range tests {
		if !t.Type().IsDir() {
			continue
		}

		paths = append(paths, testsPaths(level-1, path+"/"+t.Name())...)
	}

	return paths
}
