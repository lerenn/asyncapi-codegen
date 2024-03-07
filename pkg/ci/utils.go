package ci

import (
	"os"

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

func directoriesAtSublevel(sublevel int, path string) []string {
	paths := make([]string, 0)

	tests, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	if sublevel == 0 {
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

		paths = append(paths, directoriesAtSublevel(sublevel-1, path+"/"+t.Name())...)
	}

	return paths
}
