#!/usr/bin/env bash
set -e

# Get the tag from env.
if [ -z "$TAG" ]; then
		echo "Please specify a tag for build."
		exit 1
fi

echo "Build docker images for tag: ${TAG} ..."
docker build -t "anychaindb/node:${TAG}" -f dockerfiles/anychaindb-node.Dockerfile .
docker build --build-arg branch=${TAG} -t "anychaindb/abci:${TAG}" -f dockerfiles/anychaindb-abci.Dockerfile .
docker build --build-arg branch=${TAG} -t "anychaindb/api:${TAG}" -f dockerfiles/anychaindb-rest-api.Dockerfile .