build:
	@sh -c "'$(CURDIR)/build.sh'"

push:
	@sh -c "'$(CURDIR)/push.sh'"

build_develop:
	# node:develop image uses latest tendermint release in master branch
	@docker build -t "leadschain/node:develop" -f dockerfiles/develop/leadschain-node.Dockerfile .
	@docker build -t "leadschain/abci:develop" -f dockerfiles/develop/leadschain-abci-develop.Dockerfile .
	@docker build -t "leadschain/api:develop" -f dockerfiles/develop/leadschain-rest-api-develop.Dockerfile .

push_develop:
	@docker push "leadschain/node:develop"
	@docker push "leadschain/abci:develop"
	@docker push "leadschain/api:develop"

login:
	@docker login -u ${DOCKER_USER} -p ${DOCKER_PASS}

.PHONY: login build build_develop push push_develop