build:
	@sh -c "'$(CURDIR)/build.sh'"

push:
	@sh -c "'$(CURDIR)/push.sh'"

build_tag:
	@sh -c "'$(CURDIR)/build_tag.sh'"

push_tag:
	@sh -c "'$(CURDIR)/push_tag.sh'"	

login:
	@docker login -u ${DOCKER_USER} -p ${DOCKER_PASS}

.PHONY: login build build_tag push push_tag