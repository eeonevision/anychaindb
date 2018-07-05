build:
	@sh -c "'$(CURDIR)/build.sh'"

push:
	@sh -c "'$(CURDIR)/push.sh'"

build_develop:
	# node:develop image uses latest tendermint release in master branch
	@docker build -t "anychaindb/node:develop" -f dockerfiles/anychaindb-node.Dockerfile .
	@docker build -t "anychaindb/abci:develop" -f dockerfiles/develop/anychaindb-abci-develop.Dockerfile .
	@docker build -t "anychaindb/api:develop" -f dockerfiles/develop/anychaindb-rest-api-develop.Dockerfile .

push_develop:
	@docker push "anychaindb/node:develop"
	@docker push "anychaindb/abci:develop"
	@docker push "anychaindb/api:develop"

login:
	@docker login -u ${DOCKER_USER} -p ${DOCKER_PASS}

.PHONY: login build build_develop push push_develop