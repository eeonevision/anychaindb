#!/usr/bin/env bash
set -e

# Get the branch from env.
if [ -z "$BRANCH" ]; then
		echo "Please specify a branch for push."
		exit 1
fi

echo "Push docker images for branch: ${BRANCH} ..."
docker push "anychaindb/node:${BRANCH}"
docker push "anychaindb/abci:${BRANCH}"
docker push "anychaindb/api:${BRANCH}"
