#!/usr/bin/env bash
set -e

# Get the tag from env.
if [ -z "$TAG" ]; then
		echo "Please specify a tag for push."
		exit 1
fi

echo "Push docker images for tag: ${TAG} ..."
docker push "anychaindb/node:${TAG}"
docker push "anychaindb/abci:${TAG}"
docker push "anychaindb/api:${TAG}"
