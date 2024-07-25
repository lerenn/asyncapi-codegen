include tools/make/help.mk

.PHONY: check
check: check-generation lint test ## Run all the checks locally

.PHONY: check-generation
check-generation: ## Check files are generated locally
	@sh ./scripts/check-generation.sh

.PHONY: clean
clean: down ## Clean the project locally
	@rm -rf ./tmp/certs

.PHONY: down
down: ## Stop the local environment
	@docker-compose down

.PHONY: generate
generate: ## Generate files locally
	@go generate ./...

.PHONY: lint
lint: ## Lint the code locally
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55 run ./...

.PHONY: publish
publish: dagger/publish ## Publish with tag on git, docker hub, etc. locally
	@git tag ${TAG} && git push origin ${TAG}

.PHONY: test
test: up ## Perform tests locally
	@go test ./...

.PHONY: up
up: ## Start the local environment
	@go run ./tools/generate-certs
	@docker-compose up -d