build:
	@sh -c "'$(CURDIR)/build.sh'"

push:
	@sh -c "'$(CURDIR)/push.sh'"

build_branch:
	@sh -c "'$(CURDIR)/build_branch.sh'"

push_branch:
	@sh -c "'$(CURDIR)/push_branch.sh'"	

login:
	@docker login -u ${DOCKER_USER} -p ${DOCKER_PASS}

.PHONY: login build build_branch push push_branch