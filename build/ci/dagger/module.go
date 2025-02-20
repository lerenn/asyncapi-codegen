// A generated module for AsyncapiCodegenCi functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"

	"asyncapi-codegen/ci/dagger/internal/dagger"
)

const (
	// linterImage is the image used for linter.
	linterImage = "golangci/golangci-lint:v1.62.0"
	// golangImage is the image used as base for golang operations.
	golangImage = "golang:1.21.4-alpine"
)

// AsyncapiCodegenCi is the Dagger CI module for AsyncAPI Codegen.
type AsyncapiCodegenCi struct {
	brokers map[string]*dagger.Service
}

func (ci *AsyncapiCodegenCi) cachedBrokers() map[string]*dagger.Service {
	if ci.brokers == nil {
		ci.brokers = brokerServices()
	}
	return ci.brokers
}

// Execute all check operations (generate, lint, examples, and tests)
func (ci *AsyncapiCodegenCi) Check(
	ctx context.Context,
	srcDir *dagger.Directory,
) ([]*dagger.Container, error) {
	containers := make([]*dagger.Container, 0)

	containers = append(containers, ci.CheckGeneration(ctx, srcDir))
	containers = append(containers, ci.Lint(ctx, srcDir))

	examples, err := ci.Examples(ctx, srcDir)
	if err != nil {
		return nil, err
	}
	containers = append(containers, examples...)

	containers = append(containers, ci.Test(ctx, srcDir))

	return containers, nil
}

// CheckGeneration generate files from Golang generate command on AsyncAPI-Codegen
// source code and check that there is no change.
func (ci *AsyncapiCodegenCi) CheckGeneration(
	ctx context.Context,
	srcDir *dagger.Directory,
) *dagger.Container {
	return dag.Container().
		From(golangImage).
		With(sourceCodeAndGoCache(srcDir)).
		WithExec([]string{"sh", "./scripts/check-generation.sh"})
}

// Lint AsyncAPI-Codegen source code.
func (ci *AsyncapiCodegenCi) Lint(
	ctx context.Context,
	srcDir *dagger.Directory,
) *dagger.Container {
	return dag.Container().
		From(linterImage).
		With(sourceCodeAndGoCache(srcDir)).
		WithMountedCache("/root/.cache/golangci-lint", dag.CacheVolume("golangci-lint")).
		WithExec([]string{"golangci-lint", "run"})
}

// Examples runs AsyncAPI-Codegen examples.
func (ci *AsyncapiCodegenCi) Examples(
	ctx context.Context,
	srcDir *dagger.Directory,
) ([]*dagger.Container, error) {
	// Get examples subdirs
	subdirs, err := directoriesAtSublevel(ctx, srcDir.Directory("examples"), 2, "./examples")
	if err != nil {
		return nil, err
	}

	// Get examples containers
	containers := make([]*dagger.Container, len(subdirs))
	for i, p := range subdirs {
		// Set app container
		app := dag.Container().
			From(golangImage).
			// Add source code as work directory
			With(sourceCodeAndGoCache(srcDir)).
			// Set broker as dependency
			With(bindBrokers(ci.cachedBrokers())).
			// Execute command
			WithExec([]string{"go", "run", p + "/app"}).
			// Add exposed port to let know when the service is ready
			WithExposedPort(1234).
			// Set as service
			AsService()

		// Set user container
		user := dag.Container().
			// Add base image
			From(golangImage).
			// Add source code as work directory
			With(sourceCodeAndGoCache(srcDir)).
			// Set broker as dependency
			With(bindBrokers(ci.cachedBrokers())).
			// Add app as dependency of user
			WithServiceBinding("app", app).
			// Execute command
			WithExec([]string{"go", "run", p + "/user"})

		// Add user container to list
		containers[i] = user
	}

	return containers, nil
}

// Test run tests from AsyncAPICodegen
func (ci *AsyncapiCodegenCi) Test(
	ctx context.Context,
	srcDir *dagger.Directory,
) *dagger.Container {
	return dag.Container().
		// Add base image
		From(golangImage).
		// Add source code as work directory
		With(sourceCodeAndGoCache(srcDir)).
		// Set brokers as dependencies of app and user
		With(bindBrokers(ci.cachedBrokers())).
		// Execute command
		WithExec([]string{"go", "test", "./..."})
}

// Publish tag on git repository and docker image(s) on Docker Hub
// Note: if this is not 'main' branch, then it will just push docker image with
// git tag.
func (ci *AsyncapiCodegenCi) Publish(
	ctx context.Context,
	srcDir *dagger.Directory,
	// +optional
	sshDir *dagger.Directory,
) error {
	gi := NewGit(srcDir, sshDir)

	// Push new commit tag if needed
	if err := gi.PushNewSemVerIfNeeded(ctx); err != nil {
		return err
	}

	// Publish docker image
	return publishDocker(ctx, srcDir, gi)
}
