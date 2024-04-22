package ci

import (
	"dagger.io/dagger"
)

// Linter returns a container that runs the linter.
func Linter(client *dagger.Client) *dagger.Container {
	return client.Container().
		// Add base image
		From(LinterImage).
		// Add source code as work directory
		With(sourceAsWorkdir(client)).
		// Add golangci-lint cache
		WithMountedCache("/root/.cache/golangci-lint", client.CacheVolume("golangci-lint")).
		// Add command
		WithExec([]string{"golangci-lint", "run", "--timeout=10m"})
}
