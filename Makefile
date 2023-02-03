.PHONY: lint
lint: ## Lint the code
	@LOG_LEVEL=error go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1 run

.PHONY: test/integration
test/integration: ## Perform integration tests
	# TODO

.PHONY: test/unit
test/unit: ## Perform unit tests
	@go test ./... -coverprofile cover.out -v
	@go tool cover -func cover.out
	@rm cover.out

.PHONY: test
test: test/unit test/integration ## Perform all tests

.PHONY: generate
generate: ## Generate files
	@go generate ./...

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
