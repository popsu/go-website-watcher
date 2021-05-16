.DEFAULT_GOAL := help
MIGRATIONS_PATH = "sql/migrations"

.PHONY: help
help: ## This help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: run-consumer
run-consumer: ## Run consumer
	go run ./consumer/cmd/consumer/

.PHONY: run-producer
run-producer: ## Run producer
	go run ./producer/cmd/producer/

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

.PHONY: db-migrate-up
db-migrate-up: ## Run db migrations up
	migrate -database ${GWW_DBURL} -path ${MIGRATIONS_PATH} up

.PHONY: db-migrate-down
db-migrate-down: ## Run db migrations down
	migrate -database ${GWW_DBURL} -path ${MIGRATIONS_PATH} down

.PHONY: generate-sql
generate-sql: ## Generate sql from template in sql/migrations/templates
	go run ./sql/migrations/templates > sql/migrations/001_create_initial_table.up.sql

# terraform output -raw kafka_access_key > ../kafka_access.key

# terraform output -raw kafka_access_cert > ../kafka_access.cert
