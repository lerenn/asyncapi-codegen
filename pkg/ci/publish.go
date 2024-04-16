package ci

import (
	"context"

	"dagger.io/dagger"
	"github.com/lerenn/asyncapi-codegen/pkg/utils/git"
)

const (
	// DockerImageName is the name of the docker image.
	DockerImageName = "lerenn/asyncapi-codegen"
)

var (
	// platforms represents the different OS/Arch platform wanted for docker hub.
	platforms = []RunnerType{
		{OS: "linux", Arch: "386", BuildBaseImage: "golang:alpine", TargetBaseImage: "alpine"},
		{OS: "linux", Arch: "amd64", BuildBaseImage: "golang:alpine", TargetBaseImage: "alpine"},
		{OS: "linux", Arch: "arm/v6", BuildBaseImage: "golang:alpine", TargetBaseImage: "alpine"},
		{OS: "linux", Arch: "arm/v7", BuildBaseImage: "golang:alpine", TargetBaseImage: "alpine"},
		{OS: "linux", Arch: "arm64/v8", BuildBaseImage: "golang:alpine", TargetBaseImage: "alpine"},
		{OS: "linux", Arch: "mips64le", BuildBaseImage: "golang", TargetBaseImage: "debian:12"},
		{OS: "linux", Arch: "ppc64le", BuildBaseImage: "golang:alpine", TargetBaseImage: "alpine"},
		{OS: "linux", Arch: "s390x", BuildBaseImage: "golang:alpine", TargetBaseImage: "alpine"},
	}
)

// Publish should publish tag on git repository and docker image(s) on Docker Hub
// Note: if this is not 'main' branch, then it will just push docker image with
// git tag.
func Publish(ctx context.Context, client *dagger.Client, tag string) error {
	if err := tagAndPush(tag); err != nil {
		return err
	}

	if err := publishDocker(ctx, client, tag); err != nil {
		return err
	}

	return nil
}

func tagAndPush(tag string) error {
	// Stop here if this not main branch
	if name, err := git.ActualBranchName("."); err != nil {
		return err
	} else if name != "main" {
		return nil
	}

	// Tag commit
	if err := git.TagCommit(".", tag); err != nil {
		return err
	}

	// Push the result
	return git.PushTags(".", tag)
}

func publishDocker(ctx context.Context, client *dagger.Client, tag string) error {
	// Get images for each platform
	platformVariants := make([]*dagger.Container, len(platforms))
	for i, p := range platforms {
		platformVariants[i] = RunnerFromDockerfile(client, p)
	}

	// Set publication options from images
	publishOpts := dagger.ContainerPublishOpts{
		PlatformVariants: platformVariants,
	}

	// Get last git commit hash
	hash, err := git.GetLastCommitHash(".")
	if err != nil {
		return err
	}

	// Publish with hash
	if _, err := client.Container().Publish(ctx, DockerImageName+":"+hash, publishOpts); err != nil {
		return err
	}

	// Stop here if this not main branch
	if name, err := git.ActualBranchName("."); err != nil {
		return err
	} else if name != "main" {
		return nil
	}

	// Publish with tag passed in argument
	if _, err := client.Container().Publish(ctx, DockerImageName+":"+tag, publishOpts); err != nil {
		return err
	}

	// Publish with "latest" as tag
	if _, err := client.Container().Publish(ctx, DockerImageName+":latest", publishOpts); err != nil {
		return err
	}

	return nil
}
