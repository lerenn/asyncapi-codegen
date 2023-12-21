package pipeline

import (
	"dagger.io/dagger"
)

func sourceAsWorkdir(client *dagger.Client) func(r *dagger.Container) *dagger.Container {
	// Set path where the source code is mounted.
	containerDir := "/go/src/github.com/lerenn/asyncapi-codegen"

	return func(r *dagger.Container) *dagger.Container {
		return r.
			// Add Go caches
			WithMountedCache("/root/.cache/go-build", client.CacheVolume("gobuild")).
			WithMountedCache("/go/pkg/mod", client.CacheVolume("gocache")).

			// Add source code
			WithMountedDirectory(containerDir, client.Host().Directory(".")).

			// Add workdir
			WithWorkdir(containerDir)
	}
}
