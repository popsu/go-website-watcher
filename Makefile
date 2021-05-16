.DEFAULT_GOAL := help
MIGRATIONS_PATH = sql/migrations

GOFILES = $(shell find . -name '*.go')

DOCKER_IMG_PRODUCER = ghcr.io/popsu/gww-producer
DOCKER_TAG_PRODUCER = latest
PRODUCER_PACKAGE = github.com/popsu/go-website-watcher/producer/cmd/producer

DOCKER_IMG_CONSUMER = ghcr.io/popsu/gww-consumer
DOCKER_TAG_CONSUMER = latest
CONSUMER_PACKAGE = github.com/popsu/go-website-watcher/consumer/cmd/consumer

# build flags
CGO_ENABLED = 0
LDFLAGS = -ldflags='-s -w'
BUILDFLAGS = -tags "osusergo netgo" -trimpath

.PHONY: help
help: ## This help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Run tests
# Use richgo if installed
# https://github.com/kyoh86/richgo
ifneq (, $(shell which richgo))
	richgo test -race -cover -v ./...
else
	go test -race -cover ./...
endif

.PHONY: tools
tools: ## Install required tools (go-migrate)
	go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@v4.14.1

.PHONY: tidy
tidy: ## Tidy go modules
	go mod tidy

.PHONY: clean
clean: ## Clean artifacts
	go clean ./...
	rm -vrf bin/*

.PHONY: build
build: build-producer build-consumer ## Build producer and consumer dev

.PHONY: build-producer
build-producer: bin/dev/gww-producer tidy ## Build producer dev

.PHONY: build-consumer
build-consumer: bin/dev/gww-consumer tidy ## Build consumer dev

.PHONY: db-migrate-up
db-migrate-up: ## Run db migrations up
	migrate -database ${POSTGRES_DBURL} -path ${MIGRATIONS_PATH} up

.PHONY: db-migrate-down
db-migrate-down: ## Run db migrations down
	migrate -database ${POSTGRES_DBURL} -path ${MIGRATIONS_PATH} down

.PHONY: generate-sql
generate-sql: ## Generate sql from template in sql/migrations/templates
	go run ./sql/migrations/templates > sql/migrations/001_create_initial_table.up.sql

bin/dev/gww-producer: $(GOFILES)
	go build -race $(LDFLAGS) $(BUILDFLAGS) -o $@ ${PRODUCER_PACKAGE}

bin/dev/gww-consumer: $(GOFILES)
	go build -race $(LDFLAGS) $(BUILDFLAGS) -o $@ ${CONSUMER_PACKAGE}

bin/release/gww-producer: $(GOFILES)
	CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) $(BUILDFLAGS) -o $@ ${PRODUCER_PACKAGE}
	upx --lzma $@

bin/release/gww-consumer: $(GOFILES)
	CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) $(BUILDFLAGS) -o $@ ${CONSUMER_PACKAGE}
	upx --lzma $@

.PHONY: docker-build-producer
docker-build-producer: ## Build producer docker image
	docker build -t ${DOCKER_IMG_PRODUCER}:${DOCKER_TAG_PRODUCER} -f producer/Dockerfile .

.PHONY: docker-build-consumer
docker-build-consumer: ## Build consumer docker image
	docker build -t ${DOCKER_IMG_CONSUMER}:${DOCKER_TAG_CONSUMER} -f consumer/Dockerfile .
