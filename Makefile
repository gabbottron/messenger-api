#### BEGIN CONFIGURATION ##################################
PROJECT_NAME := "messenger"

DOCKER_API_CONTAINER := "$(PROJECT_NAME)-api"
DOCKER_API_IMAGE     := "$(PROJECT_NAME)-api_api"

DOCKER := $(shell command -v docker 2> /dev/null)
GO := $(shell command -v go 2> /dev/null)

DEV_COMPOSE_FILE  := docker-compose.yml
TEST_COMPOSE_FILE := docker-compose.test.yml
PROD_COMPOSE_FILE := docker-compose.prod.yml
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
	@docker rmi $(DOCKER_API_IMAGE)

# Local development convenience targets
connect:
	@docker exec -it $(DOCKER_DB_CONTAINER) /bin/bash

run-local: check_go
	@go run src/$(DOCKER_API_CONTAINER)/main.go

.PHONY: build-img, build, run, up, down, logs, build-test, run-test, up-test, down-test, logs-test, build-prod, run-prod, up-prod, down-prod, logs-prod, purge, connect
.DEFAULT_GOAL := up