#!/bin/bash
export GO111MODULE=on
export COMMIT_SHA=$(git rev-parse HEAD)
export export BUILD_DATE=$(date +'%s')
export KO_DOCKER_REPO="docker.io/stevenacoffman/commentary"
export LATEST_TAG=$(git tag -l --sort=-version:refname v* | head -1)
ko publish --bare -t latest -t "$COMMIT_SHA" -t "$LATEST_TAG" .