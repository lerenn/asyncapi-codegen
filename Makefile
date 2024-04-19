DAGGER_COMMAND := dagger run -v go run ./build/ci/dagger.go

ifndef EXAMPLE
	EXAMPLE=""
endif

ifndef TAG
	TAG=""
endif

.PHONY: all
all: generate lint examples test ## Run all the checks

.PHONY: ci
ci: ## Run the CI
	@${DAGGER_COMMAND} all

.PHONY: clean
clean: dev/down ## Clean the project
	@rm -rf ./tmp/certs

.PHONY: dev/up
dev/up: ## Start the development environment
	@go run ./tools/generate-certs
	@docker-compose up -d

.PHONY: dev/down
dev/down: ## Stop the development environment
	@docker-compose down

.PHONY: lint
lint: ## Lint the code
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55 run ./...

.PHONY: examples
examples: ## Perform examples
	@${DAGGER_COMMAND} examples -e ${EXAMPLE}

.PHONY: test
test: dev/up ## Perform tests
	@go test ./...

.PHONY: generate
generate: ## Generate files
	@go generate ./...

.PHONY: publish
publish: ## Publish with tag on git, docker hub, etc.
	@git tag ${TAG} && git push origin ${TAG}
	@${DAGGER_COMMAND} publish --tag ${TAG}

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
