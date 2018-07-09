#!/bin/bash

# default values for environment variables
type=""
clean_all=true
export NODE_IP="$(dig +short myip.opendns.com @resolver1.opendns.com)"
export DATA_ROOT="$HOME/anychaindb"
export CONFIG_PATH="$DATA_ROOT/config"
export DB_PORT=27017
export P2P_PORT=26656
export GRPC_PORT=26657
export ABCI_PORT=26658
export REST_PORT=26659
export NODE_ARGS=""

# set environment variables
for i in "$@"
do
case $i in
    --type=*)
        type="${i#*=}"
        shift
    ;;
    --clean_all=*)
        clean_all="${i#*=}"
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

# clean method removes anychaindb docker containers and prunes volumes/networks, 
# that was used by anychaindb
clean () {
    echo "Removing AnychainDB images"
    docker stop $(docker ps | grep "anychaindb-" | awk '/ / { print $1 }')
    docker rm $(docker ps -a | grep "anychaindb-" | awk '/ / { print $1 }')
    docker volume rm $(docker volume ls -qf dangling=true)
    docker image prune -a -f
    docker system prune -f
}

# check for starting clean process
if [ "$clean_all" = true ]; then
    clean
fi

# check conditions
if [ "$type" = "node" ]; then
    prepare
    echo "[RELEASE] Deploying AnychainDB node..."
    curl -L -O https://github.com/eeonevision/anychaindb/raw/master/deploy/DOCKER/anychaindb.yaml && \
    docker-compose -f anychaindb.yaml up -d
elif [ "$type" = "node-dev" ]; then
    prepare
    echo "[DEVELOP] Deploying AnychainDB node..."
    curl -L -O https://github.com/eeonevision/anychaindb/raw/develop/deploy/DOCKER/anychaindb-develop.yaml && \
    docker-compose -f anychaindb-develop.yaml up -d
elif [ "$type" = "update" ]; then
    echo "[RELEASE] Starting AnychainDB node..."
    curl -L -O https://github.com/eeonevision/anychaindb/raw/master/deploy/DOCKER/anychaindb.yaml && \
    docker-compose -f anychaindb.yaml up -d
elif [ "$type" = "update-dev" ]; then
    echo "[DEVELOP] Starting AnychainDB node..."
    curl -L -O https://github.com/eeonevision/anychaindb/raw/develop/deploy/DOCKER/anychaindb-develop.yaml && \
    docker-compose -f anychaindb-develop.yaml up -d
fi
