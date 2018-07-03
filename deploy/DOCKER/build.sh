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
docker build -t "leadschain/node:latest" -t "leadschain/node:$TAG" -f dockerfiles/leadschain-node.Dockerfile .
docker build -t "leadschain/abci:latest" -t "leadschain/abci:$TAG" -f dockerfiles/leadschain-abci.Dockerfile .
docker build -t "leadschain/api:latest" -t "leadschain/api:$TAG" -f dockerfiles/leadschain-rest-api.Dockerfile .
