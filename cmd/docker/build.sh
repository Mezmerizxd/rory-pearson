#!/bin/bash

# Check if version argument is provided
if [ -z "$1" ]; then
  echo "Usage: ./build.sh <version>"
  exit 1
fi

# Assign version from the first argument
VERSION=$1

# Docker Hub username and repository name
DOCKER_REPO="mezmerizxd/rory-pearson-test"

# Build the image with version and latest tags
docker build -t ${DOCKER_REPO}:${VERSION} -t ${DOCKER_REPO}:latest .

# Push the version and latest tags to Docker Hub
docker push ${DOCKER_REPO}:${VERSION}
docker push ${DOCKER_REPO}:latest

echo "Build and push completed for version ${VERSION} and latest."
