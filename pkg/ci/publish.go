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
	runner := RunnerFromDockerfile(client)

	// Publish with git commit hash as tag
	hash, err := git.GetLastCommitHash(".")
	if err != nil {
		return err
	}
	if _, err := runner.Publish(ctx, DockerImageName+":"+hash); err != nil {
		return err
	}

	// Stop here if this not main branch
	if name, err := git.ActualBranchName("."); err != nil {
		return err
	} else if name != "main" {
		return nil
	}

	// Publish with tag
	if _, err = runner.Publish(ctx, DockerImageName+":"+tag); err != nil {
		return err
	}

	// Publish as latest
	_, err = runner.Publish(ctx, DockerImageName+":latest")
	return err
}
