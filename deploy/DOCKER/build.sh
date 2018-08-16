#!/usr/bin/env bash
set -e

# Get the tag from the version, or try to figure it out.
if [ -z "$TAG" ]; then
	TAG=$(awk -F\" '/Version =/ { print $2; exit }' < ../../version/version.go)
fi
if [ -z "$TAG" ]; then
		echo "Please specify a tag."
		exit 1
fi

echo "Build two docker images with latest and ${TAG} tags ..."
docker build -t "anychaindb/node:latest" -t "anychaindb/node:$TAG" -f dockerfiles/anychaindb-node.Dockerfile .
docker build --build-arg branch=master -t "anychaindb/abci:latest" -t "anychaindb/abci:$TAG" -f dockerfiles/anychaindb-abci.Dockerfile .
docker build --build-arg branch=master -t "anychaindb/api:latest" -t "anychaindb/api:$TAG" -f dockerfiles/anychaindb-rest-api.Dockerfile .
