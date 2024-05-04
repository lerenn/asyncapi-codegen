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
	"dagger/asyncapi-codegen-ci/internal/dagger"
	"strings"

	"github.com/lerenn/asyncapi-codegen/pkg/utils"
)

const (
	// linterImage is the image used for linter.
	linterImage = "golangci/golangci-lint:v1.55"
	// golangImage is the image used as base for golang operations.
	golangImage = "golang:1.21.4-alpine"
)

// AsyncapiCodegenCi is the Dagger CI module for AsyncAPI Codegen.
type AsyncapiCodegenCi struct {
	brokers map[string]*dagger.Service
}

func (ci *AsyncapiCodegenCi) cachedBrokers() map[string]*dagger.Service {
	if ci.brokers == nil {
		ci.brokers = brokers()
	}
	return ci.brokers
}

// Execute all check operations (generate, lint, examples, and tests)
func (ci *AsyncapiCodegenCi) Check(
	ctx context.Context,
	dir *Directory,
) (string, error) {
	if _, err := ci.Generate(ctx, dir); err != nil {
		return "", err
	}

	if _, err := ci.Lint(ctx, dir); err != nil {
		return "", err
	}

	if _, err := ci.Examples(ctx, dir); err != nil {
		return "", err
	}

	if _, err := ci.Tests(ctx, dir); err != nil {
		return "", err
	}

	return "", nil
}

// Generate files from Golang generate command on AsyncAPI-Codegen source code.
func (ci *AsyncapiCodegenCi) Generate(
	ctx context.Context,
	dir *Directory,
) (string, error) {
	_, err := dag.Container().
		From(golangImage).
		With(sourceCodeAndGoCache(dir)).
		WithExec([]string{"go", "generate", "./..."}).
		Directory(".").
		Export(ctx, ".")

	return "", err
}

// Lint AsyncAPI-Codegen source code.
func (ci *AsyncapiCodegenCi) Lint(
	ctx context.Context,
	dir *Directory,
) (string, error) {
	return dag.Container().
		From(linterImage).
		With(sourceCodeAndGoCache(dir)).
		WithMountedCache("/root/.cache/golangci-lint", dag.CacheVolume("golangci-lint")).
		WithExec([]string{"golangci-lint", "run"}).
		Stdout(ctx)
}

// Run AsyncAPI-Codegen examples.
func (ci *AsyncapiCodegenCi) Examples(
	ctx context.Context,
	dir *Directory,
) (string, error) {
	// Get examples subdirs
	subdirs, err := directoriesAtSublevel(ctx, dir.Directory("examples"), 2, "./examples")
	if err != nil {
		return "", err
	}

	// Get examples containers
	containers := make(map[string]*dagger.Container, 0)
	for _, p := range subdirs {
		// Get corresponding broker (always at same level)
		brokerName := strings.Split(p, "/")[4]

		// Set app container
		app := dag.Container().
			From(golangImage).
			With(sourceCodeAndGoCache(dir)).
			WithServiceBinding(brokerName, ci.cachedBrokers()[brokerName]).
			WithExec([]string{"go", "run", p + "/app"}).
			WithExposedPort(1234).
			AsService()

		// Set user container
		user := dag.Container().
			From(golangImage).
			With(sourceCodeAndGoCache(dir)).
			WithServiceBinding(brokerName, ci.cachedBrokers()[brokerName]).
			WithServiceBinding("app", app).
			WithExec([]string{"go", "run", p + "/user"})

		// Add user containers to containers
		containers[p] = user
	}

	executeContainers(ctx, utils.MapToList(containers)...)
	return "", nil
}

// Run tests from AsyncAPICodegen
func (ci *AsyncapiCodegenCi) Tests(
	ctx context.Context,
	dir *Directory,
) (string, error) {
	// Get tests subdirs
	subdirs, err := directoriesAtSublevel(ctx, dir.Directory("test"), 2, "./test")
	if err != nil {
		return "", err
	}

	// Get tests containers
	containers := make(map[string]*dagger.Container, 0)
	for _, p := range subdirs {
		containers[p] = dag.Container().
			From(golangImage).
			With(sourceCodeAndGoCache(dir)).
			With(bindBrokers(ci.cachedBrokers())).
			WithExec([]string{"go", "test", p})
	}

	executeContainers(ctx, utils.MapToList(containers)...)
	return "", nil
}

// Publish tag on git repository and docker image(s) on Docker Hub
// Note: if this is not 'main' branch, then it will just push docker image with
// git tag.
func (ci *AsyncapiCodegenCi) Publish(
	ctx context.Context,
	dir *Directory,
	tag string,
) error {
	if err := publishDocker(ctx, dir, tag); err != nil {
		return err
	}

	return nil
}
