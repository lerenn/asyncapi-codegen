package ci

import "dagger.io/dagger"

// RunnerFromDockerfile returns a dagger container based on the repository Dockerfile.
func RunnerFromDockerfile(client *dagger.Client) *dagger.Container {
	return client.Host().Directory(".").DockerBuild(dagger.DirectoryDockerBuildOpts{
		Dockerfile: "/build/package/Dockerfile",
	})
}
