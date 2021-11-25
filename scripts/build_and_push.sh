#!/usr/bin/env bash

# Exit on fail or ctrl-c
set -e
trap "exit" INT

# Variable definitions
ROOT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
DOCKER_BASE_IMAGE_NAME="ghcr.io/papetier/scraper"

# Git commands
GIT_SHORTHASH="$(git rev-parse --short HEAD)"
GIT_VERSION="$(git describe --tags --always --abbrev=10)"

printf "> Building the Docker image \`%s\` with following tags:\n- %s\n- %s\n\n" "$DOCKER_BASE_IMAGE_NAME" "$GIT_SHORTHASH" "$GIT_VERSION"

# Docker tags
DOCKER_COMMIT_TAG="$DOCKER_BASE_IMAGE_NAME:$GIT_SHORTHASH"
DOCKER_VERSION_TAG="$DOCKER_BASE_IMAGE_NAME:$GIT_VERSION"

# Build Docker image
DOCKER_BUILDKIT=1 docker build --rm=true --ssh=default -t "$DOCKER_COMMIT_TAG" -f  "$ROOT_DIR"/../Dockerfile "$ROOT_DIR"/..

# Add version tag
docker tag "$DOCKER_COMMIT_TAG" "$DOCKER_VERSION_TAG"

# Push Docker image
printf "\n\n> Sending built image to the registry...\n"
docker push --all-tags "$DOCKER_BASE_IMAGE_NAME"
