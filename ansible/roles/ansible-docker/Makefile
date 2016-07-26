.DEFAULT_GOAL := help
.PHONY: help deps deps-role test clean

MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_PATH := $(dir $(MAKEFILE_PATH))

ANSIBLE_DIR := /opt/ansible
ANSIBLE_VOLUME := $(CURRENT_PATH):$(ANSIBLE_DIR)
ROLES_DIR := $(ANSIBLE_DIR)/tests/roles
ROLE_NAME := $(notdir $(patsubst %/,%,$(dir $(MAKEFILE_PATH))))
ROLE_TEST_PATH := $(CURRENT_PATH)tests/roles

DOCKER_IMAGE := williamyeh/ansible:debian8

help: ## target descriptions and usage
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort | awk 'BEGIN {FS = ":.*?## "}; \
		{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## ensure all required dependencies are installed
	@hash docker > /dev/null 2>&1 || \
		(echo "Error: Please install Docker to continue"; exit 1)

deps-role: clean
	@echo "* Creating temporary directory for roles"
	@mkdir -p $(ROLE_TEST_PATH)

test: deps deps-role ## run role tests
	@docker run --rm -it \
		--privileged \
		--volume $(ANSIBLE_VOLUME) \
		--volume $(ANSIBLE_VOLUME)/tests/roles/$(ROLE_NAME):ro \
		-w "$(ANSIBLE_DIR)" \
		--name $(ROLE_NAME)-test \
		"$(DOCKER_IMAGE)" \
		bash -c "tests/test.sh"

clean: 
	@echo "* Removing temporary role directory"
	@-rm -rf "$(ROLE_TEST_PATH)/"*

