package pipeline

import (
	"dagger.io/dagger"
)

// Generator returns a container that generates code.
func Generator(client *dagger.Client) *dagger.Container {
	return client.Container().
		// Add base image
		From(GolangImage).
		// Add source code as work directory
		With(sourceAsWorkdir(client)).
		// Add command to generate code
		WithExec([]string{"go", "generate", "./..."})
}
