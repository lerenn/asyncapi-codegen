.PHONY: brokers/up
brokers/up: ## Start the brokers
	@docker-compose up -d

.PHONY: brokers/logs
brokers/logs: ## Get the brokers logs
	@docker-compose logs

.PHONY: brokers/down
brokers/down: ## Stop the brokers
	@docker-compose down

.PHONY: lint
lint: ## Lint the code
	@LOG_LEVEL=error go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2 run

.PHONY: lint/fix
lint/fix: ## Fix what can be fixed regarding the linter
	@LOG_LEVEL=error go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2 run --fix

.PHONY: clean
clean: __examples/clean brokers/down ## Clean up the project

.PHONY: check
check: clean generate lint examples test ## Check that everything is ready for commit

.PHONY: __examples/clean
__examples/clean:
	@$(MAKE) -C examples clean

.PHONY: examples
examples: brokers/up ## Perform examples
	@$(MAKE) -C examples run

.PHONY: test
test: brokers/up ## Perform tests
	@go test ./... -p 1 -timeout=1m

.PHONY: generate
generate: ## Generate files
	@go generate ./...

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
