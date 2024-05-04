#!/bin/sh

# Check if git is available
if ! command -v git &> /dev/null; then
    # Lazily install it as Alpine distribution (feel free to raise an issue if
    # you need it for another distribution)
    apk add git
fi

# Execute golang code generation
go generate ./...

# Check that there is nothing to commit
git diff-index HEAD
git diff --minimal --color=always --compact-summary --exit-code HEAD || FAILED=true
if [[ $FAILED ]]; then
    echo "❗️ please run \"make generate\" locally and commit the changes"
    exit 1
fi