#!/bin/bash

echo "$(date --utc +%FT%TZ): Fetching remote repository..."
git fetch

UPSTREAM=${1:-'@{u}'}
LOCAL=$(git rev-parse @)
REMOTE=$(git rev-parse "$UPSTREAM")
BASE=$(git merge-base @ "$UPSTREAM")

if [ $LOCAL = $REMOTE ]; then
    echo "$(date --utc +%FT%TZ): No changes detected in git"
elif [ $LOCAL = $BASE ]; then
    BUILD_VERSION=$(git rev-parse HEAD)
    echo "$(date --utc +%FT%TZ): Deploying new version $BUILD_VERSION"
    ./scripts/rebuild.sh
elif [ $REMOTE = $BASE ]; then
    echo "$(date --utc +%FT%TZ): Local changes detected, stashing"
    git stash
    ./scripts/deploy.sh
else
    echo "$(date --utc +%FT%TZ): git is diverged. What did you do???"
fi