package ci

import "dagger.io/dagger"

func RunnerFromDockerfile(client *dagger.Client) *dagger.Container {
	return client.Host().Directory(".").DockerBuild(dagger.DirectoryDockerBuildOpts{
		Dockerfile: "/build/package/Dockerfile",
	})
}
