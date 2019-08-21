![anychaindb logo](https://i.imgur.com/co6woS3.png)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Feeonevision%2Fanychaindb.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Feeonevision%2Fanychaindb?ref=badge_shield)
> AnychainDB is distributed blockchain database suitables for any purposes.

[![Build Status](https://travis-ci.org/eeonevision/anychaindb.svg?branch=master)](https://travis-ci.org/eeonevision/anychaindb) [![Go Report Card](https://goreportcard.com/badge/github.com/eeonevision/anychaindb)](https://goreportcard.com/report/github.com/eeonevision/anychaindb) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

## Features
  * Blazing fast. The performance for data writing/reading is above of 2000 tx/s
  * MongoDB full-text search queries
  * Secure and private. AnychainDB allows to keep **any** kind of data in two ways: open *public* and **ECDH** encrypted *private*
  * Easy to deploy and use. Developed REST-API (see link in additional docs)
  * Clean, documented and tested code

## Installation
Officially AnychainDB provides the easiest way of installation using docker and docker-compose tools.
#### Prerequirements
Installed `docker ver. 17+` and `docker-compose` tools. Good manual are [this](https://docs.docker.com/install/linux/docker-ce/ubuntu/#install-docker-ce-1 "this") and [this](https://docs.docker.com/compose/install/#install-compose "this")

#### Choose your type of node
AnychainDB has two type of nodes:
- **Validator** node can validate transactions, configured once and need have better hardware than non-validator nodes.
- **Non-validator** node keeps the state and sends the transactions to validator nodes. The can be ease connected to the network because this type is not validates transaction.

#### Configure network
If you want to connect to existing network, then all you need is *genesis.json* and *config.toml* files from the ran network. Place in to config folder and follow to the next section. 
Configure a new network is little more sophisticated. Firstly you need to generate new config for validator node. Best way is use command from tendermint: ``tendermint init``. It creates all required configuration files for validator node (*config.toml*, *genesis.json*, *node_key.json* and *priv_validator.json*).

Examples of network configs placed [here](deploy/DOCKER/examples/schemas).
#### Deploy network
Deploy a network with the shell script:

```shell
sh deploy.sh --type=${JOB_TYPE} --node_ip=${NODE_IP} --config=${CONFIG_PATH}
```
Parameters:
* **type** - type of job for script. You can choose from *node*, *node-dev*, *update*, *update-dev*, *clean* types
* **node_ip** - ip address of node. Default script uses `dig` with `OpenDNS` as resolver. *Example: 127.0.0.1*
* **config** - full path to config folder. *Example: /home/ubuntu/CONFIG_FOLDER*
* clean_all - clean all cache and reset anychaindb state. *Default: true*
* db_port - port for communication with MongoDB container. *Default: 27017*
* p2p_port - port for communication between nodes. *Default: 26656*
* grpc_port - port for RPC client. *Default: 26657*
* abci_port - port for ABCI application. *Default: 26658*
* api_port - port for AnychainDB REST API. *Default: 26659*
* node_args - additional arguments for node. If you connect to existing network, you maybe need to set boot nodes addresses, like: `--node_args="--p2p.persistent_peers=id@host:port"`

By default script creates *AnychainDB* directory in home folder of user, where keeps all data from state. Do not remove it.
You may change it by setting another value in ${DATA_ROOT} script variable.

## Built With
  * [MongoDB] - the high-performable database with full-text search
  * [Tendermint] - the heart of blockchain platform
  * [Docker] - all components of platform wrapps in containers for fast deploy and ease to use
  * [MsgPack] - transport between messages in platform
  * [Golang] - fast and beautiful language

## Additional docs
  * [AnychainDB REST API] - REST API for AnychainDB client
  * [Tendermint Docs] - Tendermint documentation

## Contributing
Project is opensource and open for contributors and maintainers. You are welcome to do any *positive* changes in project :)

License
----
Apache 2.0

   [MongoDB]: <https://www.mongodb.com/>
   [Tendermint]: <https://github.com/tendermint/tendermint>
   [Docker]: <https://www.docker.com/>
   [MsgPack]: <https://msgpack.org/>
   [Golang]: <https://golang.org/>
   [AnychainDB REST API]: <https://anychaindb.docs.apiary.io/>
   [Tendermint Docs]: <http://tendermint.readthedocs.io/en/master/introduction.html>

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Feeonevision%2Fanychaindb.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Feeonevision%2Fanychaindb?ref=badge_large)