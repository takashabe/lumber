SUBPACKAGES := $(shell go list ./...)
APP_MAIN    := cmd/lumber/lumber.go

.DEFAULT_GOAL := help

##### Operation

build: $(APP_MAIN) ## Build application
	go build -a $(APP_MAIN)

run: $(APP_MAIN) ## Run application
	go run $(APP_MAIN)

docker-build: $(APP_MAIN) ## Build Docker image
	GOOS=linux GOARCH=amd64 go build -o release/lumber $(APP_MAIN)
	docker build -t "takashabe/lumber-api:prd" .

docker-push:  ## Push docker image to docker.io
	docker push "takashabe/lumber-api:prd"

##### Development

.PHONY: deps test vet lint

deps: ## Setup dependencies package
	dep ensure

test: ## Run go test
	go test -v -p 1 $(SUBPACKAGES)

vet: ## Check go vet
	go vet $(SUBPACKAGES)

lint: ## Check golint
	golint $(SUBPACKAGES)

##### Utilities

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
