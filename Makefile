.DEFAULT_GOAL := help

.PHONY: help
help: ## This help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: run-consumer
run-consumer: ## Run consumer
	go run ./consumer/cmd/consumer/

.PHONY: run-producer
run-producer: ## Run producer
	go run ./producer/cmd/producer/

# terraform output -raw kafka_access_key > ../kafka_access.key

# terraform output -raw kafka_access_cert > ../kafka_access.cert
