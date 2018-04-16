#!/bin/bash

echo "Export variables"
export NODE_IP="$(ip route get 8.8.8.8 | awk '{print $NF; exit}')"
#export EXTERNAL_NODE_IP="$(curl -s checkip.dyndns.org | sed -e 's/.*Current IP Address: //' -e 's/<.*$//')"
export CONFIG_PATH="$2"
export P2P_PORT="$3"
export GRPC_PORT="$4"
export ABCI_PORT="$5"
export REST_PORT="$6"
export NODE_ARGS="$7"

# move up in path
cd ../

if [ "$1" = "validator-dev" ]; then
    docker-compose -f leadschain-validator-develop.yaml up
elif [ "$1" = "node-dev" ]; then
    docker-compose -f leadschain-node-develop.yaml up
elif [ "$1" = "validator" ]; then
    docker-compose -f leadschain-validator.yaml up
elif [ "$1" = "node" ]; then
    docker-compose -f leadschain-node.yaml up
elif [ "$1" = "clean" ]; then
    echo "Prune docker images"
    docker-compose -f leadschain-validator-develop.yaml down -v --rmi all --remove-orphans && \
    docker-compose -f leadschain-node-develop.yaml down -v --rmi all --remove-orphans && \
    docker-compose -f leadschain-validator.yaml down -v --rmi all --remove-orphans && \
    docker-compose -f leadschain-node.yaml down -v --rmi all --remove-orphans && \
    docker system prune -f
fi

