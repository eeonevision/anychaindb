#!/bin/bash

# default values for environment variables
mode=""
ext_node_address=""
clean_all=true
export TAG="latest"
export DATA_ROOT="$HOME/anychaindb"
export CONFIG_PATH="$DATA_ROOT/config"
export NODE_ARGS=""

# set environment variables
for i in "$@"
do
case $i in
    --mode=*)
        type="${i#*=}"
        shift
    ;;
    --tag=*)
        export TAG="${i#*=}"
        shift
    ;;
    --clean_all=*)
        clean_all="${i#*=}"
        shift
    ;;
    --ext_node_address=*)
        ext_node_address="${i#*=}"
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

# check if config file presented
if [ -f "$CONFIG_PATH/config.toml" ]; then
    # check to replace external ip address in config.toml
    if [ "$ext_node_address" != "" ]; then
        sed -i "s|external_address = \".*\"|external_address = \"${ext_node_address}\"|" "$CONFIG_PATH/config.toml"
    fi
fi

# check mode
if [ "$mode" = "deploy" ]; then
    prepare
    echo "[${TAG}] Deploying AnychainDB node..."
    docker-compose -f ./compose.yaml up -d
elif [ "$mode" = "update" ]; then
    echo "[${TAG}] Updating AnychainDB node..."
    docker-compose -f ./compose.yaml up -d
fi
