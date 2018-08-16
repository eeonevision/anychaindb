#!/usr/bin/env bash
set -e

# Get the branch from env.
if [ -z "$BRANCH" ]; then
		echo "Please specify a branch for build."
		exit 1
fi

echo "Build docker images for branch: ${BRANCH} ..."
docker build -t "anychaindb/node:${BRANCH}" -f dockerfiles/anychaindb-node.Dockerfile .
docker build --build-arg branch=${BRANCH} -t "anychaindb/abci:${BRANCH}" -f dockerfiles/anychaindb-abci.Dockerfile .
docker build --build-arg branch=${BRANCH} -t "anychaindb/api:${BRANCH}" -f dockerfiles/anychaindb-rest-api.Dockerfile .