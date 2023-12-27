DAGGER_COMMAND := _EXPERIMENTAL_DAGGER_INTERACTIVE_TUI=true dagger run go run ./build/ci/dagger.go

ifndef EXAMPLE
	EXAMPLE=""
endif

ifndef TEST
	TEST=""
endif

.PHONY: ci
ci: ## Run the CI
	@${DAGGER_COMMAND} all

.PHONY: lint
lint: ## Lint the code
	@${DAGGER_COMMAND} linter

.PHONY: examples
examples: ## Perform examples
	@${DAGGER_COMMAND} examples -e ${EXAMPLE}

.PHONY: test
test: ## Perform tests
	@${DAGGER_COMMAND} test -t ${TEST}

.PHONY: generate
generate: ## Generate files
	@${DAGGER_COMMAND} generator

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
