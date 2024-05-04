.PHONY: local/check
local/check: local/check-generation local/lint local/test ## Run all the checks locally

.PHONY: local/clean
local/clean: local/down ## Clean the project locally
	@rm -rf ./tmp/certs

.PHONY: local/up
local/up: ## Start the local environment
	@go run ./tools/generate-certs
	@docker-compose up -d

.PHONY: local/down
local/down: ## Stop the local environment
	@docker-compose down

.PHONY: local/lint
local/lint: ## Lint the code locally
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55 run ./...

.PHONY: local/test
local/test: local/up ## Perform tests locally
	@go test ./...

.PHONY: local/generate
local/generate: ## Generate files locally
	@go generate ./...

.PHONY: local/check-generation
local/check-generation: ## Check files are generated locally
	@sh ./scripts/check-generation.sh

.PHONY: local/publish
local/publish: dagger/publish ## Publish with tag on git, docker hub, etc. locally
	@git tag ${TAG} && git push origin ${TAG}
