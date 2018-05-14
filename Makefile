SERVICE_NAME = "k8s-tracing-example"

# Shell to use for running scripts
SHELL := $(shell which bash)

# Get docker path or an empty string
DOCKER := $(shell command -v docker)

# Get docker-compose path or an empty string
DOCKER_COMPOSE := $(shell command -v docker-compose)

# Get the main unix group for the user running make (to be used by docker-compose later)
GID := $(shell id -g)

# Get the unix user id for the user running make (to be used by docker-compose later)
UID := $(shell id -u)

# Commit hash from git
COMMIT=$(shell git rev-parse --short HEAD)

# cmds
BUILD_BINARY_CMD := ./hack/scripts/build.sh
BUILD_IMAGE_CMD := IMAGE_VERSION=${COMMIT} ./hack/scripts/build-image.sh

# environment dirs
DEV_DIR := docker/dev

# The default action of this Makefile is to build the development docker image
.PHONY: default
default: build-image

# Test if the dependencies we need to run this Makefile are installed
deps-development:
ifndef DOCKER
	@echo "Docker is not available. Please install docker"
	@exit 1
endif
ifndef DOCKER_COMPOSE
	@echo "docker-compose is not available. Please install docker-compose"
	@exit 1
endif

# Run the development environment in non-daemonized mode (foreground)
dev: deps-development
	cd $(DEV_DIR) && \
	( docker-compose -p $(SERVICE_NAME) up; \
		docker-compose -p $(SERVICE_NAME) stop; \
		docker-compose -p $(SERVICE_NAME) rm -f; )

# Build production stuff.
build-binary:
	$(DOCKER_RUN_CMD) /bin/sh -c '$(BUILD_BINARY_CMD)'

.PHONY: build-image
build-image:
	$(BUILD_IMAGE_CMD)
