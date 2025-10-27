include tools/make/help.mk

.PHONY: check
check: check-generation lint test ## Run all the checks locally

.PHONY: check-generation
check-generation: ## Check files are generated locally
	@sh ./scripts/check-generation.sh

.PHONY: clean
clean: local-env/stop ## Clean the project locally
	@rm -rf ./tmp/certs

.PHONY: generate
generate: ## Generate files locally
	@go generate ./...

.PHONY: lint
lint: ## Lint the code locally
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.0 run ./...

# Dagger targets
.PHONY: dagger/check
dagger/check: ## Run all checks via Dagger
	@dagger call check --src-dir .

.PHONY: dagger/check-generation
dagger/check-generation: ## Check files are generated via Dagger
	@dagger call check-generation --src-dir .

.PHONY: dagger/develop
dagger/develop: ## Develop Dagger files
	@dagger develop

.PHONY: dagger/examples
dagger/examples: ## Run all examples via Dagger
	@dagger call examples --src-dir .

.PHONY: dagger/lint
dagger/lint: ## Lint the code via Dagger
	@dagger call lint --src-dir .

.PHONY: dagger/publish
dagger/publish: ## Publish with tag on git, docker hub, etc. via Dagger
	@dagger call publish --src-dir . --ssh-dir=~/.ssh

.PHONY: dagger/test
dagger/test: ## Perform tests via Dagger
	@dagger call test --src-dir .

.PHONY: dagger/all
dagger/all: dagger/check-generation dagger/lint dagger/test dagger/examples ## Run all CI checks via Dagger (same as GitHub workflow)

.PHONY: local-env/start
local-env/start: ## Start the local environment
	@go run ./tools/generate-certs
	@docker compose up -d

.PHONY: local-env/stop
local-env/stop: ## Stop the local environment
	@docker compose stop

.PHONY: local-env/teardown
local-env/teardown: ## Kill containers and delete volumes
	@docker compose down --volumes

.PHONY: publish
publish: dagger/publish ## Publish with tag on git, docker hub, etc. locally
	@git tag ${TAG} && git push origin ${TAG}

.PHONY: test
test: local-env/start ## Perform tests locally
	@go test ./...
