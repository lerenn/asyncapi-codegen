package ci

import (
	"runtime"

	"dagger.io/dagger"
)

// RunnerType represents the type of runner based on various info.
type RunnerType struct {
	OS              string
	Arch            string
	BuildBaseImage  string
	TargetBaseImage string
}

var (
	// DefaultRunnerType is the default runner type to use in case of doubt.
	DefaultRunnerType = RunnerType{
		OS:              "linux",
		Arch:            "amd64",
		BuildBaseImage:  "golang:alpine",
		TargetBaseImage: "alpine",
	}
)

// RunnerFromDockerfile returns a dagger container based on the repository Dockerfile.
func RunnerFromDockerfile(client *dagger.Client, rt RunnerType) *dagger.Container {
	// Get running OS, if that's an OS unsupported by Docker, replace by Linux
	os := runtime.GOOS
	if os == "darwin" {
		os = "linux"
	}

	return client.Host().Directory(".").DockerBuild(dagger.DirectoryDockerBuildOpts{
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
