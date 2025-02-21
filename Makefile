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

.PHONY: local-env/start
local-env/start: ## Start the local environment
	@go run ./tools/generate-certs
	@docker-compose up -d

.PHONY: local-env/stop
local-env/stop: ## Stop the local environment
	@docker-compose down

.PHONY: publish
publish: dagger/publish ## Publish with tag on git, docker hub, etc. locally
	@git tag ${TAG} && git push origin ${TAG}

.PHONY: test
test: local-env/start ## Perform tests locally
	@go test ./...
