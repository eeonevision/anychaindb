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

echo "Push two docker images with latest and ${TAG} tags ..."
docker push "leadschain/node"
docker push "leadschain/node:$TAG"
docker push "leadschain/abci"
docker push "leadschain/abci:$TAG"
docker push "leadschain/api"
docker push "leadschain/api:$TAG"
