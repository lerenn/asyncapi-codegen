package main

import (
	"asyncapi-codegen/ci/dagger/internal/dagger"
	"runtime"
)

// runnerType represents the type of runner based on various info.
type runnerType struct {
	OS              string
	Arch            string
	BuildBaseImage  string
	TargetBaseImage string
}

// runnerFromDockerfile returns a dagger container based on the repository Dockerfile.
func runnerFromDockerfile(dir *dagger.Directory, rt runnerType) *dagger.Container {
	// Get running OS, if that's an OS unsupported by Docker, replace by Linux
	os := runtime.GOOS
	if os == "darwin" {
		os = "linux"
	}

	return dir.DockerBuild(dagger.DirectoryDockerBuildOpts{
		BuildArgs: []dagger.BuildArg{
			{Name: "BUILDPLATFORM", Value: os + "/" + runtime.GOARCH},
			{Name: "TARGETOS", Value: rt.OS},
			{Name: "TARGETARCH", Value: rt.Arch},
			{Name: "BUILDBASEIMAGE", Value: rt.BuildBaseImage},
			{Name: "TARGETBASEIMAGE", Value: rt.TargetBaseImage},
		},
		Platform:   dagger.Platform(rt.OS + "/" + rt.Arch),
		Dockerfile: "/build/package/Dockerfile",
	})
}
