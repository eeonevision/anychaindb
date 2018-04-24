#!/bin/bash

type=""
export NODE_IP="$(ip route get 8.8.8.8 | awk '{print $NF; exit}')"
export CONFIG_PATH="~/leadschain/config"
export DB_PORT=27017
export P2P_PORT=46656
export GRPC_PORT=46657
export ABCI_PORT=46658
export REST_PORT=8889
export NODE_ARGS=""

for i in "$@"
do
case $i in
    --type=*)
        type="${i#*=}"
        shift
    ;;
    --node_ip=*)
        export NODE_IP="${i#*=}"
        shift
    ;;
    --config=*)
        export CONFIG_PATH="${i#*=}"
        shift
    ;;
    --db_port=*)
        export DB_PORT="${i#*=}"
        shift
    ;;
    --p2p_port=*)
        export P2P_PORT="${i#*=}"
        shift
    ;;
    --grpc_port=*)
        export GRPC_PORT="${i#*=}"
        shift
    ;;
    --abci_port=*)
        export ABCI_PORT="${i#*=}"
        shift
    ;;
    --api_port=*)
        export REST_PORT="${i#*=}"
        shift
    ;;
    --node_args=*)
        export NODE_ARGS="${i#*=}"
        shift
    ;;
    *)
        echo "Unknown argument: ${i}"
        exit 1
    ;;
esac
done

#create leadschain-compose folder
rm -r -f ~/leadschain/deploy
mkdir -p ~/leadschain/deploy

#move to leadschain-compose folder
cd ~/leadschain/deploy

if [ "$type" = "validator-dev" ]; then
    echo "Installing Leadschain Validator node [DEVELOP]"
    curl -L -O https://github.com/leadschain/leadschain/raw/develop/deploy/DOCKER/leadschain-validator-develop.yaml && \
    docker-compose -f leadschain-validator-develop.yaml up
elif [ "$type" = "node-dev" ]; then
    echo "Installing Leadschain Non-Validator node [DEVELOP]"
    curl -L -O https://github.com/leadschain/leadschain/raw/develop/deploy/DOCKER/leadschain-node-develop.yaml && \
    docker-compose -f leadschain-node-develop.yaml up
elif [ "$type" = "validator" ]; then
    echo "Installing Leadschain Validator node [RELEASE]"
    curl -L -O https://github.com/leadschain/leadschain/raw/master/deploy/DOCKER/leadschain-validator.yaml && \
    docker-compose -f leadschain-validator.yaml up
elif [ "$type" = "node" ]; then
    echo "Installing Leadschain Non-Validator node [RELEASE]"
    curl -L -O https://github.com/leadschain/leadschain/raw/master/deploy/DOCKER/leadschain-node.yaml && \
    docker-compose -f leadschain-node.yaml up
elif [ "$type" = "clean" ]; then
    echo "Prune docker images"
    curl -L -O https://github.com/leadschain/leadschain/raw/develop/deploy/DOCKER/leadschain-validator-develop.yaml && \
    curl -L -O https://github.com/leadschain/leadschain/raw/develop/deploy/DOCKER/leadschain-node-develop.yaml && \
    curl -L -O https://github.com/leadschain/leadschain/raw/master/deploy/DOCKER/leadschain-validator.yaml && \
    curl -L -O https://github.com/leadschain/leadschain/raw/master/deploy/DOCKER/leadschain-node.yaml && \
    docker-compose -f leadschain-validator-develop.yaml down -v --rmi all --remove-orphans && \
    docker-compose -f leadschain-node-develop.yaml down -v --rmi all --remove-orphans && \
    docker-compose -f leadschain-validator.yaml down -v --rmi all --remove-orphans && \
    docker-compose -f leadschain-node.yaml down -v --rmi all --remove-orphans && \
    docker system prune -f
fi
