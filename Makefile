#### BEGIN CONFIGURATION ##################################
PROJECT_NAME := "messenger"

# Docker container and image names
DOCKER_API_CONTAINER := "$(PROJECT_NAME)_api"
DOCKER_API_IMAGE     := "$(PROJECT_NAME)-api_api"

# Docker and go binaries if they are present
DOCKER := $(shell command -v docker 2> /dev/null)
GO := $(shell command -v go 2> /dev/null)

# Compose files for different environments
DEV_COMPOSE_FILE  := docker-compose.yml
TEST_COMPOSE_FILE := docker-compose.test.yml
PROD_COMPOSE_FILE := docker-compose.prod.yml

# Default settings for the API when run locally
# NOTE: The settings for the API when running in a
#       container with make up/run/etc come from .env
PORT ?= 8088 # The port the API will run on
DB_HOSTNAME ?= 127.0.0.1 # The hostname of the DB
DB_PORT ?= 5439 # The port of the DB
#### END CONFIGURATION ###################################

# Utility...
UNAME_S := $(shell uname -s)
MKFILE_PATH := $(lastword $(MAKEFILE_LIST))
CURRENT_DIR := $(dir $(realpath $(MKFILE_PATH)))
CURRENT_DIR := $(CURRENT_DIR:/=)
TODAY = $(shell date +%Y-%m-%d.%H:%M:%S)
TAG_DATE = $(shell date +%Y%m%d)

#### TARGETS ---------------------------------------------
check_docker:
ifndef DOCKER
	$(error "You need to install Docker. https://store.docker.com/search?type=edition&offering=community")
endif

check_go:
ifndef GO
	$(error "You need to install Golang. https://golang.org/")
endif

# Development Targets
build-img: check_docker
	@docker build --build-arg PROJECT_NAME=$(PROJECT_NAME) -t $(DOCKER_API_IMAGE) .

build: check_docker
	@PROJECT_NAME=$(PROJECT_NAME) docker-compose up --build

run: check_docker
	@PROJECT_NAME=$(PROJECT_NAME) docker-compose up -d

up: run

down: check_docker
	@PROJECT_NAME=$(PROJECT_NAME) docker-compose down

logs: check_docker
	@PROJECT_NAME=$(PROJECT_NAME) docker-compose logs -f


# Test Targets
build-test: check_docker
	@PROJECT_NAME=$(PROJECT_NAME) docker-compose -f $(TEST_COMPOSE_FILE) up --build

run-test: check_docker
	@PROJECT_NAME=$(PROJECT_NAME) docker-compose -f $(TEST_COMPOSE_FILE) up -d

up-test: run-test

down-test: check_docker
	@PROJECT_NAME=$(PROJECT_NAME) docker-compose -f $(TEST_COMPOSE_FILE) down

logs-test: check_docker
	@PROJECT_NAME=$(PROJECT_NAME) docker-compose -f $(TEST_COMPOSE_FILE) logs -f


# Production Targets
build-prod: check_docker
	@PROJECT_NAME=$(PROJECT_NAME) docker-compose -f $(PROD_COMPOSE_FILE) up --build

run-prod: check_docker
	@PROJECT_NAME=$(PROJECT_NAME) docker-compose -f $(PROD_COMPOSE_FILE) up -d

up-prod: run-prod

down-prod: check_docker
	@PROJECT_NAME=$(PROJECT_NAME) docker-compose -f $(PROD_COMPOSE_FILE) down

logs-prod: check_docker
	@PROJECT_NAME=$(PROJECT_NAME) docker-compose -f $(PROD_COMPOSE_FILE) logs -f

# utility commands
purge: check_docker
	@PROJECT_NAME=$(PROJECT_NAME) docker-compose down
	@docker rmi --force $(DOCKER_API_IMAGE)

# Local development convenience targets
connect:
	@docker exec -it $(DOCKER_API_CONTAINER) /bin/sh

run-local: check_go
	@PORT=$(PORT) DB_HOSTNAME=$(DB_HOSTNAME) DB_PORT=$(DB_PORT) go run src/$(PROJECT_NAME)-api/main.go

build-local: check_go
	@CGO_ENABLED=0 go build -o ./bin/$(PROJECT_NAME)-api -a -installsuffix cgo ./src/$(PROJECT_NAME)-api/main.go

test-local: check_go
	@DB_HOSTNAME=$(DB_HOSTNAME) DB_PORT=$(DB_PORT) go test ./...

# For Debug: Prints the value of the variable to the right of print-
#            E.G. print-PORT will print the value of $(PORT)
print-%  : ; @echo $* = $($*)

.PHONY: build-img, build, run, up, down, logs, build-test, run-test, up-test, down-test, logs-test, build-prod, run-prod, up-prod, down-prod, logs-prod, purge, connect
.DEFAULT_GOAL := up