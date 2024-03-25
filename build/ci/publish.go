package main

import (
	"context"
	"dagger/asyncapi-codegen-ci/internal/dagger"

	"github.com/lerenn/asyncapi-codegen/pkg/utils/git"
)

const (
	// dockerImageName is the name of the docker image.
	dockerImageName = "lerenn/asyncapi-codegen"
)

var (
	// platforms represents the different OS/Arch platform wanted for docker hub.
	platforms = []runnerType{
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

func tagAndPush(ctx context.Context, dir *Directory, tag string) error {
	if _, err := dir.Export(ctx, "/tmp/asyncapi-codegen"); err != nil {
		return err
	}

	// Stop here if this not main branch
	if name, err := git.ActualBranchName("/tmp/asyncapi-codegen"); err != nil {
		return err
	} else if name != "main" {
		return nil
	}

	// Tag commit
	if err := git.TagCommit("/tmp/asyncapi-codegen", tag); err != nil {
		return err
	}

	// Push the result
	return git.PushTags("/tmp/asyncapi-codegen", tag)
}

func publishDocker(ctx context.Context, dir *Directory, tag string) error {
	// Get images for each platform
	platformVariants := make([]*dagger.Container, len(platforms))
	for i, p := range platforms {
		platformVariants[i] = runnerFromDockerfile(dir, p)
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
	if _, err := dag.Container().Publish(ctx, dockerImageName+":"+hash, publishOpts); err != nil {
		return err
	}

	// Stop here if this not main branch
	if name, err := git.ActualBranchName("."); err != nil {
		return err
	} else if name != "main" {
		return nil
	}

	// Publish with tag passed in argument
	if _, err := dag.Container().Publish(ctx, dockerImageName+":"+tag, publishOpts); err != nil {
		return err
	}

	// Publish with "latest" as tag
	if _, err := dag.Container().Publish(ctx, dockerImageName+":latest", publishOpts); err != nil {
		return err
	}

	return nil
}
