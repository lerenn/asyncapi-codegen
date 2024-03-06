package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TagCommit(path, tag string) error {
	// Open git repository
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	// Get last commit
	commit, err := lastCommit(path)
	if err != nil {
		return err
	}

	// Create tag
	_, err = repo.CreateTag(tag, commit.Hash, nil)
	return err
}

func GetLastCommitHash(path string) (string, error) {
	commit, err := lastCommit(path)
	if err != nil {
		return "", err
	}

	return commit.Hash.String(), nil
}

func lastCommit(path string) (*object.Commit, error) {
	// Open git repository
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	// Get commit iterator
	commitIter, err := repo.CommitObjects()
	if err != nil {
		return nil, err
	}

	var commit *object.Commit = &object.Commit{}
	if err := commitIter.ForEach(func(c *object.Commit) error {
		if c.Author.When.After(commit.Author.When) {
			commit = c
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return commit, nil
}

func ActualBranchName(path string) (string, error) {
	// Open git repository
	repo, err := git.PlainOpen(path)
	if err != nil {
		return "", err
	}

	ref, err := repo.Head()
	if err != nil {
		return "", err
	}

	return ref.Name().Short(), nil
}

func Push(path string) error {
	// Open git repository
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	return repo.Push(&git.PushOptions{})
}
