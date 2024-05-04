#!/bin/sh

# Check if git is available
git --version 2>&1 >/dev/null # improvement by tripleee
GIT_IS_AVAILABLE=$?
if [ $GIT_IS_AVAILABLE -eq 0 ]; then
    # Lazily install it as Alpine distribution (feel free to raise an issue if you need it for another distribution)
    apk add git
fi

git diff-index HEAD
git diff --minimal --color=always --compact-summary --exit-code HEAD || FAILED=true ;
if [[ $FAILED ]];
    then echo "❗️ please run \"make generate\" locally and commit the changes"
    exit 1
fi