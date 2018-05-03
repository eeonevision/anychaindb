#!/bin/bash

# default values for environment variables
type=""
export NODE_IP="$(dig +short myip.opendns.com @resolver1.opendns.com)"
export DATA_ROOT="$HOME/leadschain"
export CONFIG_PATH="$DATA_ROOT/config"
export DB_PORT=27017
export P2P_PORT=46656
export GRPC_PORT=46657
export ABCI_PORT=46658
export REST_PORT=8889
export NODE_ARGS=""

# set environment variables
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
    --data=*)
        export DATA_ROOT="${i#*=}"
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

# prepare method erases data root directory
prepare () {
    rm -r -f $DATA_ROOT && \
    mkdir -p $DATA_ROOT/config && \
    mkdir -p $DATA_ROOT/deploy && \
    mkdir -p $DATA_ROOT/mongo && \
    cp -a $CONFIG_PATH/. $DATA_ROOT/config/ && \
    cd $DATA_ROOT/deploy
}

# clean method removes leadschain docker containers and prunes volumes/networks, 
# that was used by leadschain
clean () {
    echo "Removing Leadschain images"
    docker stop $(docker ps | grep "leadschain-" | awk '/ / { print $1 }')
    docker rm $(docker ps -a | grep "leadschain-" | awk '/ / { print $1 }')
    docker volume rm $(docker volume ls -qf dangling=true)
    docker image -a -f
    docker system prune -f
}

# check conditions
if [ "$type" = "node" ]; then
    prepare
    clean
    echo "[RELEASE] Deploying Leadschain node..."
    curl -L -O https://github.com/leadschain/leadschain/raw/master/deploy/DOCKER/leadschain.yaml && \
    docker-compose -f leadschain.yaml up
elif [ "$type" = "node-dev" ]; then
    prepare
    clean
    echo "[DEVELOP] Deploying Leadschain node..."
    curl -L -O https://github.com/leadschain/leadschain/raw/develop/deploy/DOCKER/leadschain-develop.yaml && \
    docker-compose -f leadschain-develop.yaml up
elif [ "$type" = "clean" ]; then
    clean
elif [ "$type" = "update" ]; then
    clean
    echo "[RELEASE] Starting Leadschain node..."
    curl -L -O https://github.com/leadschain/leadschain/raw/master/deploy/DOCKER/leadschain.yaml && \
    docker-compose -f leadschain.yaml up
elif [ "$type" = "update-dev" ]; then
    clean
    echo "[DEVELOP] Starting Leadschain node..."
    curl -L -O https://github.com/leadschain/leadschain/raw/develop/deploy/DOCKER/leadschain.yaml && \
    docker-compose -f leadschain.yaml up
fi
