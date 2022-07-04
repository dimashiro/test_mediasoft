SHELL := /bin/bash
VERSION := 1.0
APP_NAME = $(notdir $(CURDIR))
APP_BIN = app/build/service

# This will output the help for each task. thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Show this help
	@printf "\033[33m%s:\033[0m\n" 'Available commands'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[32m%-11s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build app binary file
	go build -o $(APP_BIN) ./app/service/main.go

build-migrate: ## Build migrate binary file
	go build -o .app/build/migrate ./app/migrate/main.go

lint: ## Run linter
	golangci-lint run

image: ## Build docker image with app
	docker build -f ./Dockerfile -t $(APP_NAME):local .
	@printf "\n   \e[30;42m %s \033[0m\n\n" 'Now you can use image like `docker run --rm $(APP_NAME):local ...`';