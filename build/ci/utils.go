package main

import (
	"context"
	"dagger/asyncapi-codegen-ci/internal/dagger"
	"sync"
)

func sourceCodeAndGoCache(dir *Directory) func(r *dagger.Container) *dagger.Container {
	sourceCodeMountPath := "/go/src/github.com/TheSadlig/asyncapi-codegen"
	return func(r *dagger.Container) *dagger.Container {
		return r.
			WithMountedCache("/root/.cache/go-build", dag.CacheVolume("gobuild")).
			WithMountedCache("/go/pkg/mod", dag.CacheVolume("gocache")).
			WithMountedDirectory(sourceCodeMountPath, dir).
			WithWorkdir(sourceCodeMountPath)
	}
}

func directoriesAtSublevel(ctx context.Context, dir *Directory, sublevel int, basePath string) ([]string, error) {
	paths := make([]string, 0)

	entries, err := dir.Entries(ctx)
	if err != nil {
		return nil, err
	}

	if sublevel == 0 {
		for _, e := range entries {
			if check, err := isDir(ctx, dir, e); err != nil {
				return nil, err
			} else if !check {
				continue
			}

			paths = append(paths, basePath+"/"+e)
		}

		return paths, nil
	}

	for _, e := range entries {
		if check, err := isDir(ctx, dir, e); err != nil {
			return nil, err
		} else if !check {
			continue
		}

		d := dir.Directory(e)
		subentries, err := directoriesAtSublevel(ctx, d, sublevel-1, basePath+"/"+e)
		if err != nil {
			return nil, err
		}
		paths = append(paths, subentries...)
	}

	return paths, nil
}

func isDir(ctx context.Context, parentDir *Directory, path string) (bool, error) {
	_, isNotDirErr := parentDir.Directory(path).Sync(ctx)
	if isNotDirErr == nil {
		// If it is a directory do not keep further checking
		return true, nil
	}

	_, isNotFileErr := parentDir.File(path).Sync(ctx)
	if isNotFileErr == nil {
		return false, nil
	}

	// At this point we know that the path does not exist or a graphql error occurred
	// We also assume that isNotDirErr and isNotFileErr are the same error
	return false, isNotFileErr
}

func executeContainers(ctx context.Context, containers ...*dagger.Container) {
	funcs := make([]func(context.Context) error, 0)
	for _, c := range containers {
		local := c
		fn := func(ctx context.Context) error {
			_, err := local.Stderr(ctx)
			return err
		}
		funcs = append(funcs, fn)
	}

	execute(ctx, funcs...)
}

func execute(ctx context.Context, funcs ...func(context.Context) error) {
	var wg sync.WaitGroup
	for _, fn := range funcs {
		go func(callback func(context.Context) error) {
			if err := callback(ctx); err != nil {
				panic(err)
			}
			wg.Done()
		}(fn)

		wg.Add(1)
	}

	wg.Wait()
}
