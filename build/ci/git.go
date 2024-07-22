package main

import (
	"context"
	"dagger/asyncapi-codegen-ci/internal/dagger"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Git struct {
	container *dagger.Container

	lastCommit struct {
		title    string
		shortSHA string
	}

	lastTag   string
	newSemVer string
}

func NewGit(srcDir, sshDir *dagger.Directory) Git {
	container := dag.Container().
		From("alpine/git").
		WithMountedDirectory("/git", srcDir).
		WithoutEntrypoint()

	if sshDir != nil {
		container = container.WithMountedDirectory("/root/.ssh", sshDir)
	}

	return Git{
		container: container,
	}
}

func (g *Git) GetLastCommitShortSHA(ctx context.Context) (string, error) {
	// Check if already doesn't exist
	if g.lastCommit.shortSHA != "" {
		return g.lastCommit.shortSHA, nil
	}

	res, err := g.container.
		WithExec([]string{"git", "rev-parse", "--short", "HEAD"}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}
	g.lastCommit.shortSHA = res

	return g.lastCommit.shortSHA, nil
}

func (g *Git) GetLastCommitTitle(ctx context.Context) (string, error) {
	// Check if already doesn't exist
	if g.lastCommit.title != "" {
		return g.lastCommit.title, nil
	}

	res, err := g.container.
		WithExec([]string{"git", "log", "-1", "--pretty=%B"}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}
	g.lastCommit.title = res

	return g.lastCommit.title, nil
}

func (g *Git) GetLastTag(ctx context.Context) (string, error) {
	// Check if already doesn't exist
	if g.lastTag != "" {
		return g.lastTag, nil
	}

	res, err := g.container.
		WithExec([]string{"git", "describe", "--tags", "--abbrev=0"}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}
	g.lastTag = res

	return g.lastTag, nil
}

func (g *Git) GetLastSemanticVersionedTag(ctx context.Context) (major, minor, patch int, err error) {
	// Get last tag
	tag, err := g.GetLastTag(ctx)
	if err != nil {
		return 0, 0, 0, err
	}

	// Remove wrong characters
	tag = strings.TrimPrefix(tag, "v")
	tag = strings.TrimSuffix(tag, "\n")

	// Split tag into parts
	parts := strings.Split(tag, ".")
	if len(parts) != 3 {
		return 0, 0, 0, errors.New("invalid tag format:" + tag)
	}

	// Convert parts to integers
	major, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, err
	}
	minor, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, err
	}
	patch, err = strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, 0, err
	}

	return major, minor, patch, nil
}

func (g *Git) GetNewSemVerIfNeeded(ctx context.Context) (string, error) {
	// Check if already doesn't exist
	if g.newSemVer != "" {
		return g.newSemVer, nil
	}

	// Get last commit title
	title, err := g.GetLastCommitTitle(ctx)
	if err != nil {
		return "", err
	}

	// Get last semantic versioned tag
	major, minor, patch, err := g.GetLastSemanticVersionedTag(ctx)
	if err != nil {
		return "", err
	}

	// Change the version
	switch {
	case strings.HasPrefix(title, "BREAKING CHANGE"):
		major += 1
	case strings.HasPrefix(title, "feat"):
		minor += 1
	case strings.HasPrefix(title, "fix"):
		patch += 1
	default:
		return "", nil
	}

	return fmt.Sprintf("v%d.%d.%d", major, minor, patch), nil
}

func (g *Git) PushNewSemVerIfNeeded(ctx context.Context) error {
	// Get new semver
	semver, err := g.GetNewSemVerIfNeeded(ctx)
	if err != nil {
		return err
	}

	// Check if semver is empty
	if semver == "" {
		return nil
	}

	// Push new tag
	_, err = g.container.
		WithExec([]string{"git", "tag", semver}).
		WithExec([]string{"git", "push", "--tags"}).
		Stdout(ctx)
	return err
}
